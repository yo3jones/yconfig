package install

import (
	"runtime"
)

type Installer interface {
	Groups(groups []string) Installer
	OnProgress(onProgress func(inst *Install)) Installer
	Install() error
}

type installer struct {
	config     *any
	inst       *Install
	os         OsType
	arch       ArchType
	groups     []string
	onProgress func(inst *Install)
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

func (instr *installer) Groups(groups []string) Installer {
	instr.groups = groups
	return instr
}

func (instr *installer) OnProgress(onProgress func(inst *Install)) Installer {
	instr.onProgress = onProgress
	return instr
}

func (instr *installer) Install() error {
	var err error

	instr.prepare()

	if err = instr.parse(); err != nil {
		return err
	}

	if err = filter(instr.inst, instr.groups); err != nil {
		return err
	}

	instr.triggerProgess()

	if err = instr.execGroups(); err != nil {
		return err
	}

	return nil
}

func (instr *installer) prepare() {
	if instr.onProgress == nil {
		instr.onProgress = func(_ *Install) {}
	}
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

func (instr *installer) triggerProgess() {
	instr.onProgress(instr.inst)
}

func (instr *installer) execGroups() error {
	groups := instr.inst.Groups

	for i := range groups {
		var (
			group = &groups[i]
			err   error
		)

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

	for i := range group.Commands {
		var (
			command = &group.Commands[i]
			err     error
		)

		if err = instr.execCommand(command); err != nil {
			return err
		}
	}
	return nil
}

func (instr *installer) execCommand(command *Command) error {
	if command.Status == Skipped {
		return nil
	}

	command.Status = Running
	instr.triggerProgess()

	var err error

	writer := newCommandWriter(instr, command)

	if err = ExecBashCommand(command.Command, writer); err != nil {
		command.Status = Error
		instr.triggerProgess()
		return err
	}

	command.Status = Complete
	instr.triggerProgess()

	return nil
}
