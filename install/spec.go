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

func OsTypeFromString(str string) OsType {
	switch strings.ToLower(str) {
	case "aix":
		return OsAix
	case "android":
		return OsAndroid
	case "darwin":
		return OsDarwin
	case "dragonfly":
		return OsDragonfly
	case "freebsd":
		return OsFreebsd
	case "hurd":
		return OsHurd
	case "illumos":
		return OsIllumos
	case "ios":
		return OsIos
	case "js":
		return OsJs
	case "linux":
		return OsLinux
	case "nacl":
		return OsNacl
	case "netbsd":
		return OsNetbsd
	case "openbsd":
		return OsOpenbsd
	case "plan9":
		return OsPlan9
	case "solaris":
		return OsSolaris
	case "windows":
		return OsWindows
	case "zos":
		return OsZos
	}
	return OsAny
}

func (osType OsType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, osType)), nil
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

func ArchTypeFromString(str string) ArchType {
	switch strings.ToLower(str) {
	case "386":
		return Arch386
	case "amd64":
		return ArchAmd64
	case "amd64p32":
		return ArchAmd64p32
	case "arm":
		return ArchArm
	case "arm64":
		return ArchArm64
	case "arm64be":
		return ArchArm64be
	case "armbe":
		return ArchArmbe
	case "loong64":
		return ArchLoong64
	case "mips":
		return ArchMips
	case "mips64":
		return ArchMips64
	case "mips64le":
		return ArchMips64le
	case "mips64p32":
		return ArchMips64p32
	case "mips64p32le":
		return ArchMips64p32le
	case "mipsle":
		return ArchMipsle
	case "ppc":
		return ArchPpc
	case "ppc64":
		return ArchPpc64
	case "ppc64le":
		return ArchPpc64le
	case "riscv":
		return ArchRiscv
	case "riscv64":
		return ArchRiscv64
	case "s390":
		return ArchS390
	case "s390x":
		return ArchS390x
	case "sparc":
		return ArchSparc
	case "sparc64":
		return ArchSparc64
	case "wasm":
		return ArchWasm
	}
	return ArchAny
}

func (archType ArchType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, archType)), nil
}

func Print(install *Install) {
	jsonBytes, err := json.MarshalIndent(install, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonBytes)
}
