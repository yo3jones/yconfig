package setup

import (
	"fmt"
	"time"

	"github.com/yo3jones/yconfig/set"
)

type Setuper interface {
	ScriptsConfig(scriptsConfig *any) Setuper
	PackageManagersConfig(packageManagersConfig *any) Setuper
	Config(config *any) Setuper
	Tags(tags []string) Setuper
	EntryNames(entryNames []string) Setuper
	Delay(delay int) Setuper
	OnProgress(onProgress func(setupState *SetupState)) Setuper
	Setup() (err error)
}

type setuper struct {
	scriptsConfig         *any
	packageManagersConfig *any
	config                *any
	tags                  *set.Set[string]
	entryNames            *set.Set[string]
	delay                 int
	onProgress            func(setupState *SetupState)
	scripts               []*SystemScript
	packageManagers       []*SystemPackageManager
	groups                []*EntryGroup
	systemScript          *SystemScript
	systemPackageManager  *SystemPackageManager
	entries               []*Entry
	state                 *SetupState
}

type Status int

const (
	StatusUknown Status = iota
	StatusWaiting
	StatusRunning
	StatusComplete
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusWaiting:
		return "waiting"
	case StatusRunning:
		return "running"
	case StatusComplete:
		return "complete"
	case StatusError:
		return "error"
	}

	return "unknown"
}

func (s Status) IsCompleted() bool {
	switch s {
	case StatusComplete:
		return true
	case StatusError:
		return true
	default:
		return false
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

type SetupState struct {
	Status       Status
	ErroredCount int
	EntryStates  []*EntryState
}

type EntryState struct {
	Entry    *Entry
	Status   Status
	Tries    int
	Retrying bool
	Out      []byte
}

func New() Setuper {
	return &setuper{}
}

func (s *setuper) Script() *SystemScript {
	return s.systemScript
}

func (s *setuper) PackageManager() *SystemPackageManager {
	return s.systemPackageManager
}

func (s *setuper) ScriptsConfig(scriptsConfig *any) Setuper {
	s.scriptsConfig = scriptsConfig
	return s
}

func (s *setuper) PackageManagersConfig(packageManagersConfig *any) Setuper {
	s.packageManagersConfig = packageManagersConfig
	return s
}

func (s *setuper) Config(config *any) Setuper {
	s.config = config
	return s
}

func (s *setuper) Tags(tags []string) Setuper {
	s.tags = set.New(tags...)
	return s
}

func (s *setuper) EntryNames(entryNames []string) Setuper {
	s.entryNames = set.New(entryNames...)
	return s
}

func (s *setuper) Delay(delay int) Setuper {
	s.delay = delay
	return s
}

func (s *setuper) OnProgress(onProgress func(setupState *SetupState)) Setuper {
	s.onProgress = onProgress
	return s
}

func (s *setuper) Setup() (err error) {
	if err = s.prepare(); err != nil {
		return err
	}

	s.notifyProgress()

	if err = s.execAll(); err != nil {
		return err
	}

	return nil
}

func (s *setuper) prepare() (err error) {
	if err = s.unmarshalConfigs(); err != nil {
		return err
	}

	// s.setup.Print()
	// fmt.Printf("\n Package Managers \n\n")
	// SlicePrint(s.packageManagers)
	// fmt.Printf("\n Scripts \n\n")
	// SlicePrint(s.scripts)

	if err = s.filter(); err != nil {
		return err
	}

	// s.systemScript.Print()
	// s.systemPackageManager.Print()
	// SlicePrint(s.values)

	s.prepareState()

	return nil
}

func (s *setuper) unmarshalConfigs() (err error) {
	if s.scripts, err = UnmarshalSystemScripts(s.scriptsConfig); err != nil {
		return err
	}

	s.packageManagers, err = UnmarshalSystemPackageManagers(
		s.packageManagersConfig,
	)
	if err != nil {
		return err
	}

	if s.groups, err = UnmarshalEntryGroups(s.config); err != nil {
		return err
	}

	return nil
}

func (s *setuper) filter() (err error) {
	filterer := NewFilterer().
		Tags(s.tags).
		EntryNames(s.entryNames)

	s.systemScript, err = filterer.FilterSystemScripts(s.scripts)
	if err != nil {
		return err
	}

	s.systemPackageManager, err = filterer.FilterSystemPackageManagers(
		s.packageManagers,
	)
	if err != nil {
		return err
	}

	if s.entries, err = filterer.FilterEntries(s.groups); err != nil {
		return err
	}

	return nil
}

func (s *setuper) prepareState() {
	setupState := &SetupState{
		Status:       StatusWaiting,
		ErroredCount: 0,
		EntryStates:  make([]*EntryState, len(s.entries)),
	}

	for i, e := range s.entries {
		setupState.EntryStates[i] = &EntryState{
			Entry:  e,
			Status: StatusWaiting,
			Out:    []byte{},
		}
	}

	s.state = setupState

	if s.onProgress == nil {
		s.onProgress = func(_ *SetupState) {}
	}
}

func (s *setuper) execAll() (err error) {
	for i := 0; !s.state.Status.IsCompleted(); i++ {
		if i >= len(s.entries) {
			i = 0
		}

		state := s.state.EntryStates[i]

		if err = s.exec(state); err != nil {
			return err
		}
	}

	return nil
}

func (s *setuper) exec(state *EntryState) (err error) {
	if state.Status.IsCompleted() {
		return nil
	}

	for {
		if err = s.execOnce(state); err != nil {
			return err
		}

		if !state.Status.IsCompleted() &&
			state.Entry.RetryBehavior == RetryBehaviorInPlace {
			continue
		}

		break
	}

	return nil
}

func (s *setuper) execOnce(state *EntryState) (err error) {
	s.doDelay()

	s.changeStatus(state, StatusRunning)

	cmd, args := state.Entry.commander.BuildCommand(s)
	writer := NewWriter(&state.Out, func() {
		s.notifyProgress()
	})

	err = Exec(cmd, args, writer)

	state.Tries++

	if err != nil && state.Entry.RetryCount+1 > state.Tries {
		state.Retrying = true
		s.changeStatus(state, StatusWaiting)
		return nil
	}

	if err != nil && state.Entry.ContinueOnError {
		s.changeStatus(state, StatusError)
		return nil
	}

	if err != nil {
		return err
	}

	s.doDelay()

	s.changeStatus(state, StatusComplete)

	return nil
}

func (s *setuper) doDelay() {
	if s.delay <= 0 {
		return
	}

	time.Sleep(time.Duration(s.delay) * time.Millisecond)
}

func (s *setuper) changeStatus(state *EntryState, status Status) {
	state.Status = status

	s.recalculateState()

	s.notifyProgress()
}

func (s *setuper) recalculateState() {
	setupStatus := StatusWaiting
	erroredCount := 0
	completedCount := 0

	for _, state := range s.state.EntryStates {
		switch state.Status {
		case StatusWaiting:
		case StatusRunning:
			setupStatus = StatusRunning
		case StatusComplete:
			completedCount++
		case StatusError:
			erroredCount++

			if state.Entry.ContinueOnError {
				setupStatus = StatusRunning
				completedCount++
			} else {
				setupStatus = StatusError
			}
		}
	}

	s.state.ErroredCount = erroredCount
	s.state.Status = setupStatus

	allComplete := completedCount >= len(s.entries)

	if allComplete && erroredCount > 0 {
		s.state.Status = StatusError
	} else if allComplete && erroredCount <= 0 {
		s.state.Status = StatusComplete
	}
}

func (s *setuper) notifyProgress() {
	s.onProgress(s.state)
}
