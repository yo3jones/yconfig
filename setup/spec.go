package setup

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/set"
)

type System interface {
	PackageManager() *PackageManager
	Script() *Script
}

type Tagger interface {
	GetTags() *set.Set[string]
	GetRequiredTags() *set.Set[string]
}

type Oser interface {
	GetOs() ostypes.Os
}

type Archer interface {
	GetArch() archtypes.Arch
}

type PackageManager struct {
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         *set.Set[string]
	RequiredTags *set.Set[string]
	Script       string
}

func (pm *PackageManager) GetOs() ostypes.Os {
	return pm.Os
}

func (pm *PackageManager) GetArch() archtypes.Arch {
	return pm.Arch
}

func (pm *PackageManager) GetTags() *set.Set[string] {
	return pm.Tags
}

func (pm *PackageManager) GetRequiredTags() *set.Set[string] {
	return pm.RequiredTags
}

func (pm *PackageManager) BuildCommand(
	script *Script,
	packages []string,
) (cmd string, args []string) {
	pmScriptParts := make([]string, 0, len(packages)+1)
	pmScriptParts = append(pmScriptParts, pm.Script)
	pmScriptParts = append(pmScriptParts, packages...)

	pmScript := strings.Join(pmScriptParts, " ")

	return script.BuildCommand(pmScript)
}

type Script struct {
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         *set.Set[string]
	RequiredTags *set.Set[string]
	Cmd          string
	Args         []string
}

func (s *Script) GetOs() ostypes.Os {
	return s.Os
}

func (s *Script) GetArch() archtypes.Arch {
	return s.Arch
}

func (s *Script) GetTags() *set.Set[string] {
	return s.Tags
}

func (s *Script) GetRequiredTags() *set.Set[string] {
	return s.RequiredTags
}

func (s *Script) BuildCommand(script string) (cmd string, args []string) {
	cmd = s.Cmd

	args = make([]string, 0, len(s.Args)+1)
	args = append(args, s.Args...)
	args = append(args, script)

	return cmd, args
}

type Setup struct {
	Entries []*Entry
}

type Entry struct {
	Name            string
	Type            *Type
	Os              ostypes.Os
	Arch            archtypes.Arch
	Tags            *set.Set[string]
	RequiredTags    *set.Set[string]
	RetryCount      int
	RetryBehavior   RetryBehavior
	ContinueOnError bool
	Values          []Value
}

func (rs RetryBehavior) MarshalJSON() ([]byte, error) {
	str := rs.String()
	return json.Marshal(&str)
}

type Value interface {
	GetName() string
	GetType() Type
	Oser
	Archer
	Tagger
	GetContinueOnError() bool
	GetRetryCount() int
	GetRetryBehavior() RetryBehavior
	BuildCommand(system System) (cmd string, args []string)
	Printer
}

type PackageValue struct {
	Name            string
	Os              ostypes.Os
	Arch            archtypes.Arch
	Tags            *set.Set[string]
	RequiredTags    *set.Set[string]
	ContinueOnError bool
	RetryCount      int
	RetryBehavior   RetryBehavior
	Packages        []string
}

func (v *PackageValue) GetName() string {
	return v.Name
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

func (v *PackageValue) GetTags() *set.Set[string] {
	return v.Tags
}

func (v *PackageValue) GetRequiredTags() *set.Set[string] {
	return v.RequiredTags
}

func (v *PackageValue) GetContinueOnError() bool {
	return v.ContinueOnError
}

func (v *PackageValue) GetRetryCount() int {
	return v.RetryCount
}

func (v *PackageValue) GetRetryBehavior() RetryBehavior {
	return v.RetryBehavior
}

func (v *PackageValue) BuildCommand(system System) (cmd string, args []string) {
	return system.PackageManager().
		BuildCommand(system.Script(), v.Packages)
}

type ScriptValue struct {
	Name            string
	Os              ostypes.Os
	Arch            archtypes.Arch
	Tags            *set.Set[string]
	RequiredTags    *set.Set[string]
	ContinueOnError bool
	RetryCount      int
	RetryBehavior   RetryBehavior
	Script          string
}

func (v *ScriptValue) GetName() string {
	return v.Name
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

func (v *ScriptValue) GetTags() *set.Set[string] {
	return v.Tags
}

func (v *ScriptValue) GetRequiredTags() *set.Set[string] {
	return v.RequiredTags
}

func (v *ScriptValue) GetContinueOnError() bool {
	return v.ContinueOnError
}

func (v *ScriptValue) GetRetryCount() int {
	return v.RetryCount
}

func (v *ScriptValue) GetRetryBehavior() RetryBehavior {
	return v.RetryBehavior
}

func (v *ScriptValue) BuildCommand(system System) (cmd string, args []string) {
	return system.Script().BuildCommand(v.Script)
}

type CommandValue struct {
	Name            string
	Os              ostypes.Os
	Arch            archtypes.Arch
	Tags            *set.Set[string]
	RequiredTags    *set.Set[string]
	ContinueOnError bool
	RetryCount      int
	RetryBehavior   RetryBehavior
	Cmd             string
	Args            []string
}

func (v *CommandValue) GetName() string {
	return v.Name
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

func (v *CommandValue) GetTags() *set.Set[string] {
	return v.Tags
}

func (v *CommandValue) GetRequiredTags() *set.Set[string] {
	return v.RequiredTags
}

func (v *CommandValue) GetContinueOnError() bool {
	return v.ContinueOnError
}

func (v *CommandValue) GetRetryCount() int {
	return v.RetryCount
}

func (v *CommandValue) GetRetryBehavior() RetryBehavior {
	return v.RetryBehavior
}

func (v *CommandValue) BuildCommand(_ System) (cmd string, args []string) {
	return v.Cmd, v.Args
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

type Printer interface {
	Print()
}

func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t)), nil
}

type RetryBehavior int

const (
	RetryBehaviorInPlace RetryBehavior = iota
	RetryBehaviorAtEnd
)

func (rs RetryBehavior) String() string {
	switch rs {
	case RetryBehaviorInPlace:
		return "IN_PLACE"
	case RetryBehaviorAtEnd:
		return "AT_END"
	default:
		return "UNKNOWN"
	}
}

func RetryBehaviorFromString(str string) (RetryBehavior, error) {
	switch str {
	case "IN_PLACE":
		return RetryBehaviorInPlace, nil
	case "AT_END":
		return RetryBehaviorAtEnd, nil
	default:
		return RetryBehaviorInPlace,
			fmt.Errorf(
				"no retry strategy for string %s",
				str,
			)
	}
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

func (script *Script) Print() {
	Print(script)
}

func (v *PackageValue) Print() {
	Print(v)
}

func (v *ScriptValue) Print() {
	Print(v)
}

func (v *CommandValue) Print() {
	Print(v)
}

func SlicePrint[T Printer, E ~[]T](slice E) {
	for _, p := range slice {
		p.Print()
	}
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
