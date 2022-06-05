package install

import (
	"encoding/json"
	"fmt"
)

func Parse(config any) *Install {
	configMap := *mustCastMap(config)
	groups := parseGroups(configMap["groups"])

	return &Install{
		Groups: groups,
	}
}

func parseGroups(config any) []Group {
	groupsSlice := *mustCast[[]any](config)
	groups := make([]Group, 0, len(groupsSlice))

	for _, groupConfig := range groupsSlice {
		groupMap := *mustCastMap(groupConfig)
		group := *parseGroup(groupMap)
		groups = append(groups, group)
	}

	return groups
}

func parseGroup(config map[string]any) *Group {
	name := mustGetStringDefault("name", config, "any")
	os := mustGetStringDefault("os", config, "any")
	arch := mustGetStringDefault("arch", config, "any")

	commandConfigs := *mustGet[[]any]("commands", config)
	commands := parseCommands(commandConfigs)

	return &Group{
		Name:     name,
		Os:       os,
		Arch:     arch,
		Commands: commands,
	}
}

func parseCommands(config []any) []Command {
	commands := make([]Command, 0, len(config))
	for _, commandConfig := range config {
		command := *parseCommand(commandConfig)
		commands = append(commands, command)
	}

	return commands
}

func parseCommand(config any) *Command {
	return mustParseEither(
		config,
		parseCommandMap,
		parseCommandString,
	)
}

func parseCommandMap(config map[string]any) *Command {
	command := *mustGet[string]("command", config)
	os := mustGetStringDefault("os", config, "any")
	arch := mustGetStringDefault("arch", config, "any")

	return &Command{
		Command: command,
		Os:      os,
		Arch:    arch,
	}
}

func parseCommandString(config string) *Command {
	return &Command{
		Command: config,
		Os:      "any",
		Arch:    "any",
	}
}

func mustParseEither[T any, E any, C any](
	config any,
	parseT func(config T) *C,
	parseE func(config E) *C,
) *C {
	result, err := parseEither(config, parseT, parseE)
	if err != nil {
		panic(err)
	}
	return result
}

func parseEither[T any, E any, C any](
	config any,
	parseT func(config T) *C,
	parseE func(config E) *C,
) (*C, error) {
	switch config := config.(type) {
	case T:
		return parseT(config), nil
	case E:
		return parseE(config), nil
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

func mustCast[T any](obj any) *T {
	result, err := cast[T](obj)
	if err != nil {
		panic(err)
	}
	return result
}

func mustCastMap(obj any) *map[string]any {
	return mustCast[map[string]any](obj)
}

func get[T any](key string, config map[string]any) (*T, error) {
	configValue, exists := config[key]
	if !exists {
		return nil, fmt.Errorf("expected map to contain %s but did not", key)
	}
	return cast[T](configValue)
}

func getDefault[T any](key string, config map[string]any, def *T) (*T, error) {
	configValue, exists := config[key]
	if !exists {
		return def, nil
	}
	return cast[T](configValue)
}

func mustGet[T any](key string, config map[string]any) *T {
	result, err := get[T](key, config)
	if err != nil {
		panic(err)
	}
	return result
}

func mustGetDefault[T any](key string, config map[string]any, def *T) *T {
	result, err := getDefault(key, config, def)
	if err != nil {
		panic(err)
	}
	return result
}

func mustGetStringDefault(
	key string,
	config map[string]any,
	def string,
) string {
	return *mustGetDefault(key, config, &def)
}

func Print(install *Install) {
	jsonBytes, err := json.MarshalIndent(install, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonBytes)
}
