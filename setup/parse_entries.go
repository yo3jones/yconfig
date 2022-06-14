package setup

import (
	"fmt"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/parse"
	"github.com/yo3jones/yconfig/set"
)

func parseEntries(config *[]any) ([]*Entry, error) {
	entries := make([]*Entry, 0, len(*config))
	for i := range *config {
		var (
			entry *Entry
			err   error
		)

		entryConfig := &(*config)[i]

		if entry, err = parseEntry(entryConfig); err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func parseEntry(config *any) (*Entry, error) {
	var (
		configMap *map[string]any
		values    []Value
		err       error
	)

	if configMap, err = parse.Cast[map[string]any](config); err != nil {
		return nil, err
	}

	var (
		exists       bool
		name         *string
		entryType    Type
		os           ostypes.Os
		arch         archtypes.Arch
		tags         *set.Set[string]
		requiredTags *set.Set[string]
	)

	if name, exists, err = parse.Get[string](configMap, "name"); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("setup entry missing required field %s", "name")
	}

	if entryType, _, err = typeGet(configMap, "type"); err != nil {
		return nil, err
	}

	if os, _, err = parse.OsGet(configMap, "os"); err != nil {
		return nil, err
	}

	if arch, _, err = parse.ArchGet(configMap, "arch"); err != nil {
		return nil, err
	}

	if tags, requiredTags, _, err = parse.TagsGet(configMap, "tags"); err != nil {
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
