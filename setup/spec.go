package setup

import (
	"encoding/json"
	"fmt"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
)

type PackageManager struct {
	Os     ostypes.Os
	Arch   archtypes.Arch
	Tags   map[string]bool
	Script string
}

type Setup struct {
	Entries []Entry
}

type Entry struct {
	Name         string
	Type         Type
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         map[string]bool
	RequiredTags map[string]bool
	Values       []Value
}

type Value interface {
	GetType() Type
	GetOs() ostypes.Os
	GetArch() archtypes.Arch
	GetTags() map[string]bool
	GetRequiredTags() map[string]bool
}

type PackageValue struct {
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         map[string]bool
	RequiredTags map[string]bool

	Packages []string
}

func (v *PackageValue) GetType() Type {
	return TypePackage
}

func (v *PackageValue) GetOs() ostypes.Os {
	return v.Os
}

func (v *PackageValue) GetArch() archtypes.Arch {
	return v.Arch
}

func (v *PackageValue) GetTags() map[string]bool {
	return v.Tags
}

func (v *PackageValue) GetRequiredTags() map[string]bool {
	return v.RequiredTags
}

type ScriptValue struct {
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         map[string]bool
	RequiredTags map[string]bool

	Script string
}

func (v *ScriptValue) GetType() Type {
	return TypeScript
}

func (v *ScriptValue) GetOs() ostypes.Os {
	return v.Os
}

func (v *ScriptValue) GetArch() archtypes.Arch {
	return v.Arch
}

func (v *ScriptValue) GetTags() map[string]bool {
	return v.Tags
}

func (v *ScriptValue) GetRequiredTags() map[string]bool {
	return v.RequiredTags
}

type CommandValue struct {
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         map[string]bool
	RequiredTags map[string]bool

	Cmd  string
	Args []string
}

func (v *CommandValue) GetType() Type {
	return TypeScript
}

func (v *CommandValue) GetOs() ostypes.Os {
	return v.Os
}

func (v *CommandValue) GetArch() archtypes.Arch {
	return v.Arch
}

func (v *CommandValue) GetTags() map[string]bool {
	return v.Tags
}

func (v *CommandValue) GetRequiredTags() map[string]bool {
	return v.RequiredTags
}

type Type int

const (
	TypeUnknown Type = iota
	TypeScript
	TypePackage
	TypeCommand
)

func (t Type) String() string {
	switch t {
	case TypeScript:
		return "script"
	case TypePackage:
		return "package"
	case TypeCommand:
		return "command"
	}
	return "unknown"
}

func TypeFromString(str string) (Type, error) {
	switch str {
	case "script":
		return TypeScript, nil
	case "package":
		return TypePackage, nil
	case "command":
		return TypeCommand, nil
	}
	return TypeUnknown, fmt.Errorf("no setup type for string %s", str)
}

func (t Type) Marshal() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t)), nil
}

func (setup *Setup) Print() {
	Print(setup)
}

func (entry *Entry) Print() {
	Print(entry)
}

func (pm *PackageManager) Print() {
	Print(pm)
}

func Print(obj any) {
	var (
		jsonOut []byte
		err     error
	)
	if jsonOut, err = json.MarshalIndent(obj, "", "  "); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", jsonOut)
}
