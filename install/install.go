package install

import (
	"runtime"
)

type Installer interface {
	Install() error
}

type installer struct {
	config any
	inst   *Install
	os     OsType
	arch   ArchType
}

func (ir *installer) Install() error {
	var (
		inst *Install
		err  error
	)

	if inst, err = Parse(ir.config); err != nil {
		return err
	}

	ir.inst = inst

	Print(inst)

	return nil
}

func New(config any) (Installer, error) {
	var (
		osType   OsType
		archType ArchType
		err      error
	)

	if osType, err = OsTypeFromString(runtime.GOOS); err != nil {
		return nil, err
	}
	if archType, err = ArchTypeFromString(runtime.GOARCH); err != nil {
		return nil, err
	}

	instr := &installer{
		config: config,
		os:     osType,
		arch:   archType,
	}

	return instr, nil
}
