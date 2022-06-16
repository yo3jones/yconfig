package setup

import (
	"fmt"

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
