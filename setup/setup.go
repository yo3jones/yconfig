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
	OnProgress(onProgress func(progress []*Progress)) Setuper
	Setup() (err error)
}

type setuper struct {
	scriptsConfig         *any
	packageManagersConfig *any
	config                *any
	tags                  *set.Set[string]
	entryNames            *set.Set[string]
	delay                 int
	onProgress            func(progress []*Progress)
	scripts               []*Script
	packageManagers       []*PackageManager
	setup                 *Setup
	systemScript          *Script
	systemPackageManager  *PackageManager
	values                []Value
	progress              []*Progress
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

func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

type Progress struct {
	Value  Value
	Status Status
	Out    []byte
}

func New() Setuper {
	return &setuper{}
}

func (s *setuper) Script() *Script {
	return s.systemScript
}

func (s *setuper) PackageManager() *PackageManager {
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

func (s *setuper) OnProgress(onProgress func(progress []*Progress)) Setuper {
	s.onProgress = onProgress
	return s
}

func (s *setuper) Setup() (err error) {
	if err = s.prepare(); err != nil {
		return err
	}

	s.onProgress(s.progress)

	if err = s.execAll(); err != nil {
		return err
	}

	return nil
}

func (s *setuper) prepare() (err error) {
	if err = s.parseConfigs(); err != nil {
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

	s.prepareProgress()

	return nil
}

func (s *setuper) parseConfigs() (err error) {
	if s.scripts, err = ParseScripts(s.scriptsConfig); err != nil {
		return err
	}

	s.packageManagers, err = ParsePackageManagers(s.packageManagersConfig)
	if err != nil {
		return err
	}

	if s.setup, err = Parse(s.config); err != nil {
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

	if s.values, err = filterer.FilterValues(s.setup); err != nil {
		return err
	}

	return nil
}

func (s *setuper) prepareProgress() {
	progressSlice := make([]*Progress, len(s.values))
	for i, v := range s.values {
		progressSlice[i] = &Progress{
			Value:  v,
			Status: StatusWaiting,
			Out:    []byte{},
		}
	}
	s.progress = progressSlice

	if s.onProgress == nil {
		s.onProgress = func(_ []*Progress) {}
	}
}

func (s *setuper) execAll() (err error) {
	for _, progress := range s.progress {
		if err = s.exec(progress); err != nil {
			return err
		}
	}

	return nil
}

func (s *setuper) exec(progress *Progress) (err error) {
	s.doDelay()

	progress.Status = StatusRunning
	s.onProgress(s.progress)

	cmd, args := progress.Value.BuildCommand(s)
	writer := NewWriter(&progress.Out, func() {
		s.onProgress(s.progress)
	})

	if err = Exec(cmd, args, writer); err != nil {
		progress.Status = StatusError
		s.onProgress(s.progress)

		if progress.Value.GetContinueOnError() {
			return nil
		} else {
			return err
		}
	}

	s.doDelay()

	progress.Status = StatusComplete
	s.onProgress(s.progress)

	return nil
}

func (s *setuper) doDelay() {
	if s.delay <= 0 {
		return
	}

	time.Sleep(time.Duration(s.delay) * time.Millisecond)
}
