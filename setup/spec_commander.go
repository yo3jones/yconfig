package setup

import (
	"fmt"
	"strings"

	"github.com/yo3jones/yconfig/parse"
)

type EntryCommander interface {
	BuildCommand(system System) (cmd string, args []string)
}

type EntryCommanderUnmarshaler interface {
	EntryCommander

	unmarshalMapDefaults(m, defaults *map[string]any) (err error)
}

func unmarshalEntryString(str *string) *Entry {
	commander := &PackageEntry{
		packages: []string{*str},
	}
	return NewEntry(*str, commander)
}

func inferEntryTypeDefault(m, defaults *map[string]any) (t Type, err error) {
	var exists bool

	if t, exists, err = inferEntryType(m); err != nil {
		return t, err
	} else if exists {
		return t, nil
	}

	if t, exists, err = inferEntryType(defaults); err != nil {
		return t, err
	} else if exists {
		return t, nil
	}

	return t, fmt.Errorf("unable to infer type")
}

func inferEntryType(m *map[string]any) (t Type, exists bool, err error) {
	if t, exists, err = typeGet(m, "type"); err != nil {
		return t, false, err
	} else if exists {
		return t, true, nil
	}

	if _, exists, err = parse.Get[any](m, "packages"); err != nil {
		return t, false, err
	} else if exists {
		return TypePackage, true, nil
	}

	if _, exists, err = parse.Get[any](m, "script"); err != nil {
		return t, false, err
	} else if exists {
		return TypeScript, true, nil
	}

	if _, exists, err = parse.Get[any](m, "cmd"); err != nil {
		return t, false, err
	} else if exists {
		return TypeCommand, true, nil
	}

	if _, exists, err = parse.Get[any](m, "repo"); err != nil {
		return t, false, err
	} else if exists {
		return TypeGit, true, nil
	}

	return t, false, nil
}

func newEntryCommander(t Type) (EntryCommanderUnmarshaler, error) {
	switch t {
	case TypePackage:
		return &PackageEntry{}, nil
	case TypeScript:
		return &ScriptEntry{}, nil
	case TypeCommand:
		return &CommandEntry{}, nil
	case TypeGit:
		return &GitEntry{}, nil
	}
	return nil, fmt.Errorf("unable to instantiate entry for type %s", t)
}

type PackageEntry struct {
	packages []string
}

func (e *PackageEntry) unmarshalMapDefaults(
	m, defaults *map[string]any,
) (err error) {
	e.packages, _, err = parse.StringSliceGetDefaultMap(m, "packages", defaults)
	if err != nil {
		return err
	}
	return nil
}

func (e *PackageEntry) BuildCommand(
	system System,
) (cmd string, args []string) {
	return system.PackageManager().BuildCommand(system.Script(), e.packages)
}

type GitEntry struct {
	repo     string
	dest     string
	depth    int
	behavior GitCleanupBehavior
}

type GitCleanupBehavior int

const (
	GitCleanupBehaviorRemove GitCleanupBehavior = iota
	GitCleanupBehaviorNothing
)

func GitCleanupBehaviorFromString(
	str string,
) (behavior GitCleanupBehavior, err error) {
	switch str {
	case "REMOVE":
		return GitCleanupBehaviorRemove, nil
	case "NOTHING":
		return GitCleanupBehaviorNothing, nil
	}
	return behavior, fmt.Errorf("unknown GitCleanupBehavior for string %s", str)
}

func (behavior GitCleanupBehavior) String() string {
	switch behavior {
	case GitCleanupBehaviorRemove:
		return "REMOVE"
	case GitCleanupBehaviorNothing:
		return "NOTHING"
	}

	return "UNKNOWN"
}

func GitCleanupBehaviorGet(
	m *map[string]any,
	key string,
) (behavior GitCleanupBehavior, exists bool, err error) {
	var str string
	if str, exists, err = parse.StringGet(m, key); err != nil {
		return behavior, false, err
	} else if !exists {
		return behavior, false, nil
	}

	if behavior, err = GitCleanupBehaviorFromString(str); err != nil {
		return behavior, true, err
	}

	return behavior, true, nil
}

func GitCleanupBehaviorGetDefaultMap(
	m *map[string]any,
	key string,
	defaults *map[string]any,
) (behavior GitCleanupBehavior, exists bool, err error) {
	if behavior, exists, err = GitCleanupBehaviorGet(m, key); err != nil {
		return behavior, false, err
	} else if exists {
		return behavior, true, nil
	}

	if defaults == nil {
		return behavior, false, nil
	}

	return GitCleanupBehaviorGet(defaults, key)
}

func (e *GitEntry) unmarshalMapDefaults(
	m, defaults *map[string]any,
) (err error) {
	var exists bool

	e.repo, exists, err = parse.StringGetDefaultMap(m, "repo", defaults)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("git setup entry requires field repo")
	}

	e.dest, exists, err = parse.StringGetDefaultMap(m, "dest", defaults)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("git setup entry requires field dest")
	}

	e.depth, exists, err = parse.IntGetDefaultMap(m, "depth", defaults)
	if err != nil {
		return err
	} else if !exists {
		e.depth = 1
	}

	e.behavior, _, err = GitCleanupBehaviorGetDefaultMap(
		m,
		"behavior",
		defaults,
	)
	if err != nil {
		return err
	}

	return nil
}

func (e *GitEntry) BuildCommand(
	system System,
) (cmd string, args []string) {
	repo := e.repo
	if !strings.HasPrefix(strings.ToLower(repo), "https://") {
		repo = fmt.Sprintf("https://github.com/%s", repo)
	}
	if !strings.HasSuffix(strings.ToLower(repo), ".git") {
		repo = fmt.Sprintf("%s.git", repo)
	}

	depth := fmt.Sprintf(" --depth=%d ", e.depth)
	if e.depth < 1 {
		depth = " "
	}

	script := fmt.Sprintf(
		"git clone%s\\\n    %s \\\n    %s",
		depth,
		repo,
		e.dest,
	)

	if e.behavior == GitCleanupBehaviorRemove {
		script = fmt.Sprintf("rm -rf %s \\\n; %s", e.dest, script)
	}

	return system.Script().BuildCommand(script)
}

type ScriptEntry struct {
	script string
}

func (e *ScriptEntry) unmarshalMapDefaults(
	m, defaults *map[string]any,
) (err error) {
	var script *string
	script, _, err = parse.GetDefaultMap[string](m, "script", defaults)
	if err != nil {
		return err
	}
	e.script = *script
	return nil
}

func (e *ScriptEntry) BuildCommand(
	system System,
) (cmd string, args []string) {
	return system.Script().BuildCommand(e.script)
}

type CommandEntry struct {
	cmd  string
	args []string
}

func (e *CommandEntry) unmarshalMapDefaults(
	m, defaults *map[string]any,
) (err error) {
	var cmd *string
	cmd, _, err = parse.GetDefaultMap[string](m, "cmd", defaults)
	if err != nil {
		return err
	}
	e.cmd = *cmd
	e.args, _, err = parse.StringSliceGetDefaultMap(m, "args", defaults)
	if err != nil {
		return err
	}
	return nil
}

func (e *CommandEntry) BuildCommand(
	system System,
) (cmd string, args []string) {
	return e.cmd, e.args
}
