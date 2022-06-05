package install

import (
	"fmt"
)

func Parse(config any) (*Install, error) {
	var (
		configMap *map[string]any
		groups    []Group
		err       error
	)

	if configMap, err = cast[map[string]any](config); err != nil {
		return nil, err
	}
	if groups, err = parseGroups((*configMap)["groups"]); err != nil {
		return nil, err
	}

	install := &Install{
		Groups: groups,
	}

	return install, nil
}

func parseGroups(config any) ([]Group, error) {
	var (
		groupsSlice *[]any
		err         error
	)

	if groupsSlice, err = cast[[]any](config); err != nil {
		return nil, err
	}
	groups := make([]Group, 0, len(*groupsSlice))

	for _, groupConfig := range *groupsSlice {
		var (
			groupMap *map[string]any
			group    *Group
		)

		if groupMap, err = cast[map[string]any](groupConfig); err != nil {
			return nil, err
		}
		if group, err = parseGroup(*groupMap); err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}

	return groups, nil
}

func parseGroup(config map[string]any) (*Group, error) {
	var (
		name           string
		os             string
		arch           string
		commandConfigs *[]any
		commands       []Command
		err            error
	)

	if name, err = getStringDefault("name", config, "any"); err != nil {
		return nil, err
	}
	if os, err = getStringDefault("os", config, "any"); err != nil {
		return nil, err
	}
	if arch, err = getStringDefault("arch", config, "any"); err != nil {
		return nil, err
	}

	if commandConfigs, err = get[[]any]("commands", config); err != nil {
		return nil, err
	}
	if commands, err = parseCommands(*commandConfigs); err != nil {
		return nil, err
	}

	group := &Group{
		Name:     name,
		Os:       OsTypeFromString(os),
		Arch:     ArchTypeFromString(arch),
		Commands: commands,
	}

	return group, nil
}

func parseCommands(config []any) ([]Command, error) {
	commands := make([]Command, 0, len(config))
	for _, commandConfig := range config {
		var (
			command *Command
			err     error
		)

		if command, err = parseCommand(commandConfig); err != nil {
			return nil, err
		}

		commands = append(commands, *command)
	}

	return commands, nil
}

func parseCommand(config any) (*Command, error) {
	return parseEither(
		config,
		parseCommandMap,
		parseCommandString,
	)
}

func parseCommandMap(config map[string]any) (*Command, error) {
	var (
		commandString string
		os            string
		arch          string
		err           error
	)

	if commandString, err = getString("command", config); err != nil {
		return nil, err
	}
	if os, err = getStringDefault("os", config, "any"); err != nil {
		return nil, err
	}
	if arch, err = getStringDefault("arch", config, "any"); err != nil {
		return nil, err
	}

	command := &Command{
		Command: commandString,
		Os:      OsTypeFromString(os),
		Arch:    ArchTypeFromString(arch),
	}

	return command, nil
}

func parseCommandString(config string) (*Command, error) {
	command := &Command{
		Command: config,
		Os:      OsAny,
		Arch:    ArchAny,
	}

	return command, nil
}

func parseEither[T any, E any, C any](
	config any,
	parseT func(config T) (*C, error),
	parseE func(config E) (*C, error),
) (*C, error) {
	switch config := config.(type) {
	case T:
		return parseT(config)
	case E:
		return parseE(config)
	}

	return nil, fmt.Errorf(
		"expected %T or %T but got %T",
		*new(T),
		*new(E),
		config,
	)
}

func cast[T any](obj any) (*T, error) {
	switch obj := obj.(type) {
	case T:
		return &obj, nil
	}
	return nil, fmt.Errorf("expected a map[string] but got %T", obj)
}

func get[T any](key string, config map[string]any) (*T, error) {
	configValue, exists := config[key]
	if !exists {
		return nil, fmt.Errorf("expected map to contain %s but did not", key)
	}
	return cast[T](configValue)
}

func getString(key string, config map[string]any) (string, error) {
	result, err := get[string](key, config)
	if err != nil {
		return "", err
	}
	return *result, nil
}

func getDefault[T any](key string, config map[string]any, def *T) (*T, error) {
	configValue, exists := config[key]
	if !exists {
		return def, nil
	}
	return cast[T](configValue)
}

func getStringDefault(
	key string,
	config map[string]any,
	def string,
) (string, error) {
	result, err := getDefault(key, config, &def)
	if err != nil {
		return "", err
	}
	return *result, nil
}
