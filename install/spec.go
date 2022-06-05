package install

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Install struct {
	Groups []Group
}

type Group struct {
	Name     string
	Commands []Command
	Os       OsType
	Arch     ArchType
}

type Command struct {
	Command string
	Os      OsType
	Arch    ArchType
}

type HasEnv interface {
	GetOs() OsType
	GetArch() ArchType
}

func (group *Group) GetOs() OsType {
	return group.Os
}

func (group *Group) GetArch() ArchType {
	return group.Arch
}

func (command *Command) GetOs() OsType {
	return command.Os
}

func (command *Command) GetArch() ArchType {
	return command.Arch
}

type OsType int

const (
	OsAny OsType = iota
	OsAix
	OsAndroid
	OsDarwin
	OsDragonfly
	OsFreebsd
	OsHurd
	OsIllumos
	OsIos
	OsJs
	OsLinux
	OsNacl
	OsNetbsd
	OsOpenbsd
	OsPlan9
	OsSolaris
	OsWindows
	OsZos
)

func (osType OsType) String() string {
	switch osType {
	case OsAix:
		return "aix"
	case OsAndroid:
		return "android"
	case OsDarwin:
		return "darwin"
	case OsDragonfly:
		return "dragonfly"
	case OsFreebsd:
		return "freebsd"
	case OsHurd:
		return "hurd"
	case OsIllumos:
		return "illumos"
	case OsIos:
		return "ios"
	case OsJs:
		return "js"
	case OsLinux:
		return "linux"
	case OsNacl:
		return "nacl"
	case OsNetbsd:
		return "netbsd"
	case OsOpenbsd:
		return "openbsd"
	case OsPlan9:
		return "plan9"
	case OsSolaris:
		return "solaris"
	case OsWindows:
		return "windows"
	case OsZos:
		return "zos"
	}
	return "any"
}

func OsTypeFromString(str string) (OsType, error) {
	switch strings.ToLower(str) {
	case "any":
		return OsAny, nil
	case "aix":
		return OsAix, nil
	case "android":
		return OsAndroid, nil
	case "darwin":
		return OsDarwin, nil
	case "dragonfly":
		return OsDragonfly, nil
	case "freebsd":
		return OsFreebsd, nil
	case "hurd":
		return OsHurd, nil
	case "illumos":
		return OsIllumos, nil
	case "ios":
		return OsIos, nil
	case "js":
		return OsJs, nil
	case "linux":
		return OsLinux, nil
	case "nacl":
		return OsNacl, nil
	case "netbsd":
		return OsNetbsd, nil
	case "openbsd":
		return OsOpenbsd, nil
	case "plan9":
		return OsPlan9, nil
	case "solaris":
		return OsSolaris, nil
	case "windows":
		return OsWindows, nil
	case "zos":
		return OsZos, nil
	}
	return OsAny, fmt.Errorf("unknown os type %s", str)
}

func (osType OsType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, osType)), nil
}

func (osType OsType) ShouldExec(otherOsType OsType) bool {
	if otherOsType == OsAny {
		return true
	}
	return osType == otherOsType
}

type ArchType int

const (
	ArchAny ArchType = iota
	Arch386
	ArchAmd64
	ArchAmd64p32
	ArchArm
	ArchArm64
	ArchArm64be
	ArchArmbe
	ArchLoong64
	ArchMips
	ArchMips64
	ArchMips64le
	ArchMips64p32
	ArchMips64p32le
	ArchMipsle
	ArchPpc
	ArchPpc64
	ArchPpc64le
	ArchRiscv
	ArchRiscv64
	ArchS390
	ArchS390x
	ArchSparc
	ArchSparc64
	ArchWasm
)

func (archType ArchType) String() string {
	switch archType {
	case Arch386:
		return "386"
	case ArchAmd64:
		return "Amd64"
	case ArchAmd64p32:
		return "Amd64p32"
	case ArchArm:
		return "Arm"
	case ArchArm64:
		return "Arm64"
	case ArchArm64be:
		return "Arm64be"
	case ArchArmbe:
		return "Armbe"
	case ArchLoong64:
		return "Loong64"
	case ArchMips:
		return "Mips"
	case ArchMips64:
		return "Mips64"
	case ArchMips64le:
		return "Mips64le"
	case ArchMips64p32:
		return "Mips64p32"
	case ArchMips64p32le:
		return "Mips64p32le"
	case ArchMipsle:
		return "Mipsle"
	case ArchPpc:
		return "Ppc"
	case ArchPpc64:
		return "Ppc64"
	case ArchPpc64le:
		return "Ppc64le"
	case ArchRiscv:
		return "Riscv"
	case ArchRiscv64:
		return "Riscv64"
	case ArchS390:
		return "S390"
	case ArchS390x:
		return "S390x"
	case ArchSparc:
		return "Sparc"
	case ArchSparc64:
		return "Sparc64"
	case ArchWasm:
		return "Wasm"
	}
	return "any"
}

func ArchTypeFromString(str string) (ArchType, error) {
	switch strings.ToLower(str) {
	case "any":
		return ArchAny, nil
	case "386":
		return Arch386, nil
	case "amd64":
		return ArchAmd64, nil
	case "amd64p32":
		return ArchAmd64p32, nil
	case "arm":
		return ArchArm, nil
	case "arm64":
		return ArchArm64, nil
	case "arm64be":
		return ArchArm64be, nil
	case "armbe":
		return ArchArmbe, nil
	case "loong64":
		return ArchLoong64, nil
	case "mips":
		return ArchMips, nil
	case "mips64":
		return ArchMips64, nil
	case "mips64le":
		return ArchMips64le, nil
	case "mips64p32":
		return ArchMips64p32, nil
	case "mips64p32le":
		return ArchMips64p32le, nil
	case "mipsle":
		return ArchMipsle, nil
	case "ppc":
		return ArchPpc, nil
	case "ppc64":
		return ArchPpc64, nil
	case "ppc64le":
		return ArchPpc64le, nil
	case "riscv":
		return ArchRiscv, nil
	case "riscv64":
		return ArchRiscv64, nil
	case "s390":
		return ArchS390, nil
	case "s390x":
		return ArchS390x, nil
	case "sparc":
		return ArchSparc, nil
	case "sparc64":
		return ArchSparc64, nil
	case "wasm":
		return ArchWasm, nil
	}
	return ArchAny, fmt.Errorf("unknown arch type %s", str)
}

func (archType ArchType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, archType)), nil
}

func (archType ArchType) ShouldExec(otherArchType ArchType) bool {
	if otherArchType == ArchAny {
		return true
	}
	return archType == otherArchType
}

func (inst *Install) GetGroupByName(name string) (*Group, error) {
	for _, group := range inst.Groups {
		if group.Name == name {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("no group found with name %s", name)
}

func (inst *Install) Print() {
	jsonBytes, err := json.MarshalIndent(inst, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonBytes)
}

func (group *Group) Print() {
	jsonBytes, err := json.MarshalIndent(group, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonBytes)
}

func (command *Command) Print() {
	jsonBytes, err := json.MarshalIndent(command, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonBytes)
}
