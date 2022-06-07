package install

import (
	"runtime"
	"strings"
)

type context struct {
	all         bool
	runtimeOs   OsType
	runtimeArch ArchType
	groupNames  map[string]bool
}

func newContext(groupNames []string, all bool) (*context, error) {
	var (
		runtimeOs   OsType
		runtimeArch ArchType
		err         error
	)

	if runtimeOs, err = OsTypeFromString(runtime.GOOS); err != nil {
		return nil, err
	}
	if runtimeArch, err = ArchTypeFromString(runtime.GOARCH); err != nil {
		return nil, err
	}

	groupNamesSet := map[string]bool{}
	for _, groupName := range groupNames {
		groupNamesSet[strings.ToLower(groupName)] = true
	}

	context := &context{
		all:         all,
		runtimeOs:   runtimeOs,
		runtimeArch: runtimeArch,
		groupNames:  groupNamesSet,
	}

	return context, nil
}

func (context *context) shouldExecute(config HasEnv) bool {
	if !shouldExecuteOs(context.runtimeOs, config.GetOs()) {
		return false
	}
	if !shouldExecuteArch(context.runtimeArch, config.GetArch()) {
		return false
	}
	return true
}

func (context *context) getConfigInitialStatus(config HasEnv) Status {
	if !context.shouldExecute(config) {
		return Skipped
	}
	return Waiting
}

func shouldExecuteOs(runtimeOs, configOs OsType) bool {
	if configOs == OsAny {
		return true
	}
	if configOs == runtimeOs {
		return true
	}
	return false
}

func shouldExecuteArch(runtimeArch, configArch ArchType) bool {
	if configArch == ArchAny {
		return true
	}
	if configArch == runtimeArch {
		return true
	}
	return false
}

func filter(inst *Install, groupNames []string, all bool) error {
	var (
		cxt *context
		err error
	)

	if cxt, err = newContext(groupNames, all); err != nil {
		return err
	}

	filterGroups(cxt, inst.Groups)

	return nil
}

func filterGroups(cxt *context, groups []Group) {
	for i := range groups {
		filterGroup(cxt, &groups[i])
	}
}

func filterGroup(cxt *context, group *Group) {
	group.Status = cxt.getConfigInitialStatus(group)

	if !cxt.all {
		_, contains := cxt.groupNames[strings.ToLower(group.Name)]
		if !contains {
			group.Status = Skipped
		}
	}

	filterCommands(cxt, group, group.Commands)
}

func filterCommands(cxt *context, group *Group, commands []Command) {
	for i := range commands {
		filterCommand(cxt, group, &commands[i])
	}
}

func filterCommand(cxt *context, group *Group, command *Command) {
	if group.Status == Skipped {
		command.Status = Skipped
	} else {
		command.Status = cxt.getConfigInitialStatus(command)
	}
}
