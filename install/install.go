package install

import (
	"os"
	"runtime"
)

type Installer interface {
	Groups(groups []string) Installer
	Install() error
}

type installer struct {
	config *any
	inst   *Install
	os     OsType
	arch   ArchType
	groups []string
}

func (instr *installer) Groups(groups []string) Installer {
	instr.groups = groups
	return instr
}

func (instr *installer) Install() error {
	var err error

	if err = instr.parse(); err != nil {
		return err
	}

	if err = filter(instr.inst); err != nil {
		return err
	}

	instr.inst.Print()

	if err = instr.execGroups(); err != nil {
		return err
	}

	return nil
}

func New(config *any) (Installer, error) {
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

func (instr *installer) parse() error {
	var (
		inst *Install
		err  error
	)

	if inst, err = Parse(instr.config); err != nil {
		return err
	}

	instr.inst = inst

	return nil
}

func (instr *installer) execGroups() error {
	for _, groupName := range instr.groups {
		var (
			group *Group
			err   error
		)

		if group, err = instr.inst.GetGroupByName(groupName); err != nil {
			return err
		}

		if err = instr.execGroup(group); err != nil {
			return err
		}
	}

	return nil
}

func (instr *installer) execGroup(group *Group) error {
	if group.Status == Skipped {
		return nil
	}

	for _, command := range group.Commands {
		var err error

		if err = instr.execCommand(&command); err != nil {
			return err
		}
	}
	return nil
}

func (instr *installer) execCommand(command *Command) error {
	if command.Status == Skipped {
		return nil
	}

	var err error

	if err = ExecBashCommand(command.Command, os.Stdout); err != nil {
		return err
	}

	return nil
}
