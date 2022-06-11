package setup

import (
	"fmt"
	"strings"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
)

type valueFields struct {
	t            Type
	os           ostypes.Os
	arch         archtypes.Arch
	tags         map[string]bool
	requiredTags map[string]bool
}

func Parse(config *any) (*Setup, error) {
	var (
		entriesConfig *[]any
		entries       []Entry
		err           error
	)

	if entriesConfig, err = cast[[]any](config); err != nil {
		return nil, err
	}

	if entries, err = parseEntries(entriesConfig); err != nil {
		return nil, err
	}

	setup := &Setup{
		Entries: entries,
	}

	return setup, nil
}

func parseEntries(config *[]any) ([]Entry, error) {
	entries := make([]Entry, 0, len(*config))
	for i := range *config {
		var (
			entry *Entry
			err   error
		)

		entryConfig := &(*config)[i]

		if entry, err = parseEntry(entryConfig); err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}

func parseEntry(config *any) (*Entry, error) {
	var (
		configMap *map[string]any
		values    []Value
		err       error
	)

	if configMap, err = cast[map[string]any](config); err != nil {
		return nil, err
	}

	var (
		exists       bool
		name         *string
		entryType    Type
		os           ostypes.Os
		arch         archtypes.Arch
		tags         map[string]bool
		requiredTags map[string]bool
	)

	if name, exists, err = get[string](configMap, "name"); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("setup entry missing required field %s", "name")
	}

	if entryType, _, err = typeGet(configMap, "type"); err != nil {
		return nil, err
	}

	if os, _, err = osGet(configMap, "os"); err != nil {
		return nil, err
	}

	if arch, _, err = archGet(configMap, "arch"); err != nil {
		return nil, err
	}

	if tags, requiredTags, _, err = tagsGet(configMap, "tags"); err != nil {
		return nil, err
	}

	entry := &Entry{
		Name:         *name,
		Type:         entryType,
		Os:           os,
		Arch:         arch,
		Tags:         tags,
		RequiredTags: requiredTags,
	}

	if values, err = parseValues(configMap, entry); err != nil {
		return nil, err
	}

	entry.Values = values

	return entry, nil
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

	if config, exists, err = get[any](entriesConfig, "values"); err != nil {
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

	if valuesConfig, err = cast[[]any](config); err != nil {
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

	if configMap, err = cast[map[string]any](config); err != nil {
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

	if os, exists, err = osGet(config, "os"); err != nil {
		return nil, err
	} else if !exists {
		os = entry.Os
	}

	if arch, exists, err = archGet(config, "arch"); err != nil {
		return nil, err
	} else if !exists {
		arch = entry.Arch
	}

	if tags, requiredTags, _, err = tagsGet(config, "tags"); err != nil {
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

	if packages, exists, err = stringSliceGet(config, "packages"); err != nil {
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

	if script, exists, err = get[string](config, "script"); err != nil {
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

	if cmd, exists, err = get[string](config, "cmd"); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(
			"setup entry %s script typed value requires property cmd",
			entry.Name,
		)
	}

	if args, exists, err = stringSliceGet(config, "args"); err != nil {
		return err
	} else if !exists {
		args = &[]string{}
	}

	value.Cmd = *cmd
	value.Args = *args

	return nil
}

func cast[T any](obj *any) (*T, error) {
	switch obj := (*obj).(type) {
	case T:
		return &obj, nil
	}

	return nil, fmt.Errorf("error casting type %T", obj)
}

func stringSliceCast(obj *any) (ptrStringSlice *[]string, err error) {
	var slice *[]any

	if slice, err = cast[[]any](obj); err != nil {
		return nil, err
	}

	stringSlice := make([]string, len(*slice))
	for i, val := range *slice {
		var str *string

		if str, err = cast[string](&val); err != nil {
			return nil, err
		}

		stringSlice[i] = *str
	}

	ptrStringSlice = &stringSlice

	return ptrStringSlice, nil
}

func get[T any](
	obj *map[string]any,
	key string,
) (value *T, exists bool, err error) {
	rawVal, exists := (*obj)[key]
	if !exists {
		return nil, false, nil
	}

	if value, err = cast[T](&rawVal); err != nil {
		return nil, true, err
	}

	return value, true, nil
}

func typeGet(
	obj *map[string]any,
	key string,
) (t Type, exists bool, err error) {
	var str *string

	if str, exists, err = get[string](obj, key); err != nil {
		return t, exists, err
	} else if !exists {
		return t, false, nil
	}

	if t, err = TypeFromString(*str); err != nil {
		return t, true, err
	}

	return t, true, err
}

func osGet(
	obj *map[string]any,
	key string,
) (os ostypes.Os, exists bool, err error) {
	var str *string

	if str, exists, err = get[string](obj, key); err != nil {
		return os, exists, err
	} else if !exists {
		return os, false, nil
	}

	if os, err = ostypes.OsFromString(*str); err != nil {
		return os, true, err
	}

	return os, true, err
}

func archGet(
	obj *map[string]any,
	key string,
) (arch archtypes.Arch, exists bool, err error) {
	var str *string

	if str, exists, err = get[string](obj, key); err != nil {
		return arch, exists, err
	} else if !exists {
		return arch, false, nil
	}

	if arch, err = archtypes.ArchFromString(*str); err != nil {
		return arch, true, err
	}

	return arch, true, err
}

func tagsGet(
	obj *map[string]any,
	key string,
) (tags, requiredTags map[string]bool, exists bool, err error) {
	var rawTags *[]string

	if rawTags, exists, err = stringSliceGet(obj, key); err != nil {
		return nil, nil, exists, err
	} else if !exists {
		return nil, nil, false, nil
	}

	tags = map[string]bool{}
	requiredTags = map[string]bool{}

	for _, tag := range *rawTags {
		if strings.HasSuffix(tag, "!") {
			normalizedTag := tag[:len(tag)-1]
			requiredTags[normalizedTag] = true
			tags[normalizedTag] = true
		} else {
			tags[tag] = true
		}
	}

	return tags, requiredTags, true, nil
}

func stringSliceGet(
	obj *map[string]any,
	key string,
) (strs *[]string, exists bool, err error) {
	var (
		val  *any
		vals []string
	)

	if val, exists, err = get[any](obj, key); err != nil {
		return nil, exists, err
	} else if !exists {
		return nil, false, nil
	}

	switch castedVal := (*val).(type) {
	case string:
		vals = make([]string, 1)
		vals[0] = castedVal
	case []any:
		var ptrVals *[]string
		if ptrVals, err = stringSliceCast(val); err != nil {
			return nil, true, err
		}
		vals = *ptrVals
	default:
		return nil,
			true,
			fmt.Errorf("expected either a string or []string but got %T", val)
	}

	strs = &vals

	return strs, true, nil
}
