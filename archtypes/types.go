package archtypes

import (
	"errors"
	"fmt"
	"strings"
)

type Arch int

const (
	Any Arch = iota
	Arch386
	Amd64
	Amd64p32
	Arm
	Arm64
	Arm64be
	Armbe
	Loong64
	Mips
	Mips64
	Mips64le
	Mips64p32
	Mips64p32le
	Mipsle
	Ppc
	Ppc64
	Ppc64le
	Riscv
	Riscv64
	S390
	S390x
	Sparc
	Sparc64
	Wasm
)

const anyStr = "any"

var ArchNotFoundError = errors.New("unknown arch type for string")

// nolint:funlen,cyclop
func (arch Arch) String() string {
	switch arch {
	case Arch386:
		return "386"
	case Amd64:
		return "amd64"
	case Amd64p32:
		return "amd64p32"
	case Arm:
		return "arm"
	case Arm64:
		return "arm64"
	case Arm64be:
		return "arm64be"
	case Armbe:
		return "armbe"
	case Loong64:
		return "loong64"
	case Mips:
		return "mips"
	case Mips64:
		return "mips64"
	case Mips64le:
		return "mips64le"
	case Mips64p32:
		return "mips64p32"
	case Mips64p32le:
		return "mips64p32le"
	case Mipsle:
		return "mipsle"
	case Ppc:
		return "ppc"
	case Ppc64:
		return "ppc64"
	case Ppc64le:
		return "ppc64le"
	case Riscv:
		return "riscv"
	case Riscv64:
		return "riscv64"
	case S390:
		return "s390"
	case S390x:
		return "s390x"
	case Sparc:
		return "sparc"
	case Sparc64:
		return "sparc64"
	case Wasm:
		return "wasm"
	case Any:
		return anyStr
	}

	return anyStr
}

// nolint:funlen,cyclop
func ArchFromString(str string) (Arch, error) {
	switch strings.ToLower(str) {
	case anyStr:
		return Any, nil
	case "386":
		return Arch386, nil
	case "amd64":
		return Amd64, nil
	case "amd64p32":
		return Amd64p32, nil
	case "arm":
		return Arm, nil
	case "arm64":
		return Arm64, nil
	case "arm64be":
		return Arm64be, nil
	case "armbe":
		return Armbe, nil
	case "loong64":
		return Loong64, nil
	case "mips":
		return Mips, nil
	case "mips64":
		return Mips64, nil
	case "mips64le":
		return Mips64le, nil
	case "mips64p32":
		return Mips64p32, nil
	case "mips64p32le":
		return Mips64p32le, nil
	case "mipsle":
		return Mipsle, nil
	case "ppc":
		return Ppc, nil
	case "ppc64":
		return Ppc64, nil
	case "ppc64le":
		return Ppc64le, nil
	case "riscv":
		return Riscv, nil
	case "riscv64":
		return Riscv64, nil
	case "s390":
		return S390, nil
	case "s390x":
		return S390x, nil
	case "sparc":
		return Sparc, nil
	case "sparc64":
		return Sparc64, nil
	case "wasm":
		return Wasm, nil
	}

	return Any, fmt.Errorf("%w : %s", ArchNotFoundError, str)
}

func (arch Arch) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, arch)), nil
}
