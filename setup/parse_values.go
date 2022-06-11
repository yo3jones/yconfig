package setup

import (
	"fmt"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/parse"
)

type valueFields struct {
	t            Type
	os           ostypes.Os
	arch         archtypes.Arch
	tags         map[string]bool
	requiredTags map[string]bool
}

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

	inferedType = inferType(entriesConfig, entry)

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

func inferType(config *map[string]any, entry *Entry) (t Type) {
	var exists bool

	if entry.Type != TypeUnknown {
		return entry.Type
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
		vf        *valueFields
	)

	if configMap, err = parse.Cast[map[string]any](config); err != nil {
		return nil, err
	}

	if vf, err = parseValueFields(configMap, entry); err != nil {
		return nil, err
	}

	switch vf.t {
	case TypePackage:
		value, err = parsePackageValue(configMap, vf, entry)
	case TypeScript:
		value, err = parseScriptValue(configMap, vf, entry)
	case TypeCommand:
		value, err = parseCommandValue(configMap, vf, entry)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func parseValueFields(
	config *map[string]any,
	entry *Entry,
) (vf *valueFields, err error) {
	var (
		exists       bool
		t            Type
		os           ostypes.Os
		arch         archtypes.Arch
		tags         map[string]bool
		requiredTags map[string]bool
	)

	if t, exists, err = typeGet(config, "type"); err != nil {
		return nil, err
	} else if !exists && entry.Type != TypeUnknown {
		t = entry.Type
	} else if !exists {
		return nil,
			fmt.Errorf(
				"setup entry with name %s does not have required field type",
				entry.Name,
			)
	}

	if os, exists, err = parse.OsGet(config, "os"); err != nil {
		return nil, err
	} else if !exists {
		os = entry.Os
	}

	if arch, exists, err = parse.ArchGet(config, "arch"); err != nil {
		return nil, err
	} else if !exists {
		arch = entry.Arch
	}

	if tags, requiredTags, _, err = parse.TagsGet(config, "tags"); err != nil {
		return nil, err
	}

	for tag := range entry.Tags {
		tags[tag] = true
	}

	for tag := range entry.RequiredTags {
		requiredTags[tag] = true
	}

	vf = &valueFields{
		t:            t,
		os:           os,
		arch:         arch,
		tags:         tags,
		requiredTags: requiredTags,
	}

	return vf, nil
}

func parsePackageValueEntryMap(
	config *map[string]any,
	entry *Entry,
) (values []Value, err error) {
	value := &PackageValue{
		Os:           entry.Os,
		Arch:         entry.Arch,
		Tags:         entry.Tags,
		RequiredTags: entry.RequiredTags,
	}

	if err = parsePackageSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return []Value{value}, nil
}

func parsePackageValue(
	config *map[string]any,
	vf *valueFields,
	entry *Entry,
) (value *PackageValue, err error) {
	value = &PackageValue{
		Os:           vf.os,
		Arch:         vf.arch,
		Tags:         vf.tags,
		RequiredTags: vf.requiredTags,
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
		packages *[]string
	)

	if packages, exists, err = parse.StringSliceGet(config, "packages"); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(
			"setup entry %s package typed value requires property packages",
			entry.Name,
		)
	}

	value.Packages = *packages

	return nil
}

func parseScriptValueEntryMap(
	config *map[string]any,
	entry *Entry,
) (values []Value, err error) {
	value := &ScriptValue{
		Os:           entry.Os,
		Arch:         entry.Arch,
		Tags:         entry.Tags,
		RequiredTags: entry.RequiredTags,
	}

	if err = parseScriptValueSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return []Value{value}, nil
}

func parseScriptValue(
	config *map[string]any,
	vf *valueFields,
	entry *Entry,
) (value *ScriptValue, err error) {
	value = &ScriptValue{
		Os:           vf.os,
		Arch:         vf.arch,
		Tags:         vf.tags,
		RequiredTags: vf.requiredTags,
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
		Os:           entry.Os,
		Arch:         entry.Arch,
		Tags:         entry.Tags,
		RequiredTags: entry.RequiredTags,
	}

	if err = parseCommandValueSpecifics(config, value, entry); err != nil {
		return nil, err
	}

	return []Value{value}, nil
}

func parseCommandValue(
	config *map[string]any,
	vf *valueFields,
	entry *Entry,
) (value *CommandValue, err error) {
	value = &CommandValue{
		Os:           vf.os,
		Arch:         vf.arch,
		Tags:         vf.tags,
		RequiredTags: vf.requiredTags,
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
		args   *[]string
	)

	if cmd, exists, err = parse.Get[string](config, "cmd"); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(
			"setup entry %s script typed value requires property cmd",
			entry.Name,
		)
	}

	if args, exists, err = parse.StringSliceGet(config, "args"); err != nil {
		return err
	} else if !exists {
		args = &[]string{}
	}

	value.Cmd = *cmd
	value.Args = *args

	return nil
}
