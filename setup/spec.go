package setup

import (
	"encoding/json"
	"fmt"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/set"
)

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

func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t)), nil
}

func (rs RetryBehavior) MarshalJSON() ([]byte, error) {
	str := rs.String()
	return json.Marshal(&str)
}

func (entry *Entry) Print() {
	Print(entry)
}

func (pm *SystemPackageManager) Print() {
	Print(pm)
}

func (script *SystemScript) Print() {
	Print(script)
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
