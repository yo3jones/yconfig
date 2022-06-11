package ostypes

import (
	"fmt"
	"strings"
)

type Os int

const (
	Any Os = iota
	Aix
	Android
	Darwin
	Dragonfly
	Freebsd
	Hurd
	Illumos
	Ios
	Js
	Linux
	Nacl
	Netbsd
	Openbsd
	Plan9
	Solaris
	Windows
	Zos
)

func (os Os) String() string {
	switch os {
	case Aix:
		return "aix"
	case Android:
		return "android"
	case Darwin:
		return "darwin"
	case Dragonfly:
		return "dragonfly"
	case Freebsd:
		return "freebsd"
	case Hurd:
		return "hurd"
	case Illumos:
		return "illumos"
	case Ios:
		return "ios"
	case Js:
		return "js"
	case Linux:
		return "linux"
	case Nacl:
		return "nacl"
	case Netbsd:
		return "netbsd"
	case Openbsd:
		return "openbsd"
	case Plan9:
		return "plan9"
	case Solaris:
		return "solaris"
	case Windows:
		return "windows"
	case Zos:
		return "zos"
	}
	return "any"
}

func OsFromString(str string) (Os, error) {
	switch strings.ToLower(str) {
	case "any":
		return Any, nil
	case "aix":
		return Aix, nil
	case "android":
		return Android, nil
	case "darwin":
		return Darwin, nil
	case "dragonfly":
		return Dragonfly, nil
	case "freebsd":
		return Freebsd, nil
	case "hurd":
		return Hurd, nil
	case "illumos":
		return Illumos, nil
	case "ios":
		return Ios, nil
	case "js":
		return Js, nil
	case "linux":
		return Linux, nil
	case "nacl":
		return Nacl, nil
	case "netbsd":
		return Netbsd, nil
	case "openbsd":
		return Openbsd, nil
	case "plan9":
		return Plan9, nil
	case "solaris":
		return Solaris, nil
	case "windows":
		return Windows, nil
	case "zos":
		return Zos, nil
	}
	return Any, fmt.Errorf("unknown os type for string%s", str)
}

func (os Os) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, os)), nil
}
