package setup

import (
	"fmt"

	"github.com/yo3jones/yconfig/parse"
)

func parseValues(
	entriesConfig *map[string]any,
	entry *Entry,
) (values []Value, err error) {
	var (
		exists      bool
		config      *any
		inferedType Type
	)

	if config, exists, err = parse.Get[any](entriesConfig, "values"); err != nil {
		return nil, err
	} else if exists {
		return parseFromValues(config, entry)
	}

	inferedType = inferType(entriesConfig, entry.Type)

	switch inferedType {
	case TypePackage:
		return parsePackageValueEntryMap(entriesConfig, entry)
	case TypeScript:
		return parseScriptValueEntryMap(entriesConfig, entry)
	case TypeCommand:
		return parseCommandValueEntryMap(entriesConfig, entry)
	}

	return nil,
		fmt.Errorf(
			"cannot infer type from setup entry %s",
			entry.Name,
		)
}

func parseFromValues(config *any, entry *Entry) (values []Value, err error) {
	var valuesConfig *[]any

	if valuesConfig, err = parse.Cast[[]any](config); err != nil {
		return nil, err
	}

	values = make([]Value, len(*valuesConfig))

	for i := range *valuesConfig {
		var value Value

		valueConfig := &(*valuesConfig)[i]

		if value, err = parseValue(valueConfig, entry); err != nil {
			return nil, err
		}

		values[i] = value
	}

	return values, nil
}

func parseValue(config *any, entry *Entry) (value Value, err error) {
	var (
		configMap *map[string]any
		cf        *commonFields
	)

	if configMap, err = parse.Cast[map[string]any](config); err != nil {
		return nil, err
	}

	if cf, err = parseCommonFields(configMap); err != nil {
		return nil, err
	}

	if cf.t == nil {
		cf.t = entry.Type
	}
	inferedType := inferType(configMap, cf.t)
	cf.t = &inferedType

	if cf.os == nil {
		cf.os = &entry.Os
	}

	if cf.arch == nil {
		cf.arch = &entry.Arch
	}

	cf.tags.PutAll(entry.Tags.Iter()...)

	cf.requiredTags.PutAll(entry.RequiredTags.Iter()...)

	if cf.continueOnError == nil {
		cf.continueOnError = &entry.ContinueOnError
	}

	if cf.retryCount == nil {
		cf.retryCount = &entry.RetryCount
	}

	if cf.retryBehavior == nil {
		cf.retryBehavior = &entry.RetryBehavior
	}

	switch *cf.t {
	case TypePackage:
		value, err = parsePackageValue(configMap, cf, entry)
	case TypeScript:
		value, err = parseScriptValue(configMap, cf, entry)
	case TypeCommand:
		value, err = parseCommandValue(configMap, cf, entry)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func inferType(config *map[string]any, inType *Type) (t Type) {
	var exists bool

	if inType != nil {
		return *inType
	}

	if _, exists = (*config)["packages"]; exists {
		return TypePackage
	}

	if _, exists = (*config)["script"]; exists {
		return TypeScript
	}

	if _, exists = (*config)["cmd"]; exists {
		return TypeCommand
	}

	return t
}

func parsePackageValueEntryMap(
	config *map[string]any,
	entry *Entry,
) (values []Value, err error) {
	value := &PackageValue{
		Name:            entry.Name,
		Os:              entry.Os,
		Arch:            entry.Arch,
		Tags:            entry.Tags,
		RequiredTags:    entry.RequiredTags,
		ContinueOnError: entry.ContinueOnError,
		RetryCount:      entry.RetryCount,
		RetryBehavior:   entry.RetryBehavior,
	}

	if err = parsePackageSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return []Value{value}, nil
}

func parsePackageValue(
	config *map[string]any,
	cf *commonFields,
	entry *Entry,
) (value *PackageValue, err error) {
	value = &PackageValue{
		Name:            entry.Name,
		Os:              *cf.os,
		Arch:            *cf.arch,
		Tags:            cf.tags,
		RequiredTags:    cf.requiredTags,
		ContinueOnError: *cf.continueOnError,
		RetryCount:      *cf.retryCount,
		RetryBehavior:   *cf.retryBehavior,
	}

	if err = parsePackageSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return value, nil
}

func parsePackageSpecifics(
	config *map[string]any,
	value *PackageValue,
	entry *Entry,
) (err error) {
	var (
		exists   bool
		packages []string
	)

	if packages, exists, err = parse.StringSliceGet(config, "packages"); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(
			"setup entry %s package typed value requires property packages",
			entry.Name,
		)
	}

	value.Packages = packages

	return nil
}

func parseScriptValueEntryMap(
	config *map[string]any,
	entry *Entry,
) (values []Value, err error) {
	value := &ScriptValue{
		Name:            entry.Name,
		Os:              entry.Os,
		Arch:            entry.Arch,
		Tags:            entry.Tags,
		RequiredTags:    entry.RequiredTags,
		ContinueOnError: entry.ContinueOnError,
		RetryCount:      entry.RetryCount,
		RetryBehavior:   entry.RetryBehavior,
	}

	if err = parseScriptValueSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return []Value{value}, nil
}

func parseScriptValue(
	config *map[string]any,
	cf *commonFields,
	entry *Entry,
) (value *ScriptValue, err error) {
	value = &ScriptValue{
		Name:            entry.Name,
		Os:              *cf.os,
		Arch:            *cf.arch,
		Tags:            cf.tags,
		RequiredTags:    cf.requiredTags,
		ContinueOnError: *cf.continueOnError,
		RetryCount:      *cf.retryCount,
		RetryBehavior:   *cf.retryBehavior,
	}

	if err = parseScriptValueSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return value, nil
}

func parseScriptValueSpecifics(
	config *map[string]any,
	value *ScriptValue,
	entry *Entry,
) (err error) {
	var (
		exists bool
		script *string
	)

	if script, exists, err = parse.Get[string](config, "script"); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(
			"setup entry %s script typed value requires property script",
			entry.Name,
		)
	}

	value.Script = *script

	return nil
}

func parseCommandValueEntryMap(
	config *map[string]any,
	entry *Entry,
) (values []Value, err error) {
	value := &CommandValue{
		Name:            entry.Name,
		Os:              entry.Os,
		Arch:            entry.Arch,
		Tags:            entry.Tags,
		RequiredTags:    entry.RequiredTags,
		ContinueOnError: entry.ContinueOnError,
		RetryCount:      entry.RetryCount,
		RetryBehavior:   entry.RetryBehavior,
	}

	if err = parseCommandValueSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return []Value{value}, nil
}

func parseCommandValue(
	config *map[string]any,
	cf *commonFields,
	entry *Entry,
) (value *CommandValue, err error) {
	value = &CommandValue{
		Name:            entry.Name,
		Os:              *cf.os,
		Arch:            *cf.arch,
		Tags:            cf.tags,
		RequiredTags:    cf.requiredTags,
		ContinueOnError: *cf.continueOnError,
		RetryCount:      *cf.retryCount,
		RetryBehavior:   *cf.retryBehavior,
	}

	if err = parseCommandValueSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return value, nil
}

func parseCommandValueSpecifics(
	config *map[string]any,
	value *CommandValue,
	entry *Entry,
) (err error) {
	var (
		exists bool
		cmd    *string
		args   []string
	)

	if cmd, exists, err = parse.Get[string](config, "cmd"); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(
			"setup entry %s script typed value requires property cmd",
			entry.Name,
		)
	}

	if args, _, err = parse.StringSliceGet(config, "args"); err != nil {
		return err
	}

	value.Cmd = *cmd
	value.Args = args

	return nil
}
