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
	switch config := (*config).(type) {
	case string:
		return parseEntryString(config)
	case map[string]any:
		return parseEntryMap(&config)
	default:
		return nil,
			fmt.Errorf(
				"setup entries must be of type string or map[string]any got %T",
				config,
			)
	}
}

func parseEntryString(entryString string) (entry *Entry, err error) {
	t := TypePackage

	return &Entry{
		Name:         entryString,
		Type:         &t,
		Tags:         set.New[string](),
		RequiredTags: set.New[string](),
		Values: []Value{&PackageValue{
			Name:         entryString,
			Tags:         set.New[string](),
			RequiredTags: set.New[string](),
			Packages:     []string{entryString},
		}},
	}, nil
}

func parseEntryMap(config *map[string]any) (entry *Entry, err error) {
	var (
		exists bool
		name   *string
		cf     *commonFields
		values []Value
	)

	entry = &Entry{}

	if name, exists, err = parse.Get[string](config, "name"); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("setup entry missing required field %s", "name")
	} else {
		entry.Name = *name
	}

	if cf, err = parseCommonFields(config); err != nil {
		return nil, err
	}

	entry.Type = cf.t

	if cf.os == nil {
		entry.Os = ostypes.Any
	} else {
		entry.Os = *cf.os
	}

	if cf.arch == nil {
		entry.Arch = archtypes.Any
	} else {
		entry.Arch = *cf.arch
	}

	entry.Tags = cf.tags

	entry.RequiredTags = cf.requiredTags

	if cf.continueOnError != nil {
		entry.ContinueOnError = *cf.continueOnError
	}

	if values, err = parseValues(config, entry); err != nil {
		return nil, err
	}

	entry.Values = values

	return entry, nil
}
