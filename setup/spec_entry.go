package setup

import (
	"fmt"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/parse"
	"github.com/yo3jones/yconfig/set"
)

type EntryGroup struct {
	Entries []*Entry
}

func UnmarshalEntryGroups(a *any) (groups []*EntryGroup, err error) {
	var s *[]any

	if s, err = parse.Cast[[]any](a); err != nil {
		return nil, err
	}

	groups = make([]*EntryGroup, len(*s))
	for i := range *s {
		groupAny := (*s)[i]

		group := &EntryGroup{}
		if err = group.unmarshalAny(&groupAny); err != nil {
			return nil, err
		}

		groups[i] = group
	}

	return groups, nil
}

func (eg *EntryGroup) unmarshalAny(a *any) (err error) {
	switch a := (*a).(type) {
	case string:
		return eg.unmarshalString(&a)
	case map[string]any:
		return eg.unmarshalMap(&a)
	}

	return fmt.Errorf(
		"setup entry group must be either a string or map[string]any but got %T",
		a,
	)
}

func (eg *EntryGroup) unmarshalString(str *string) error {
	entry := unmarshalEntryString(str)

	eg.Entries = []*Entry{entry}

	return nil
}

func (eg *EntryGroup) unmarshalMap(m *map[string]any) (err error) {
	var (
		exists    bool
		entryAnys *[]any
	)

	if entryAnys, exists, err = parse.Get[[]any](m, "entries"); err != nil {
		return err
	}

	if !exists {
		eg.Entries = make([]*Entry, 0, 1)
		return eg.unmarshalEntryMapDefaults(m, nil)
	}

	eg.Entries = make([]*Entry, 0, len(*entryAnys))
	for _, entryAny := range *entryAnys {
		var entryMap *map[string]any
		if entryMap, err = parse.Cast[map[string]any](&entryAny); err != nil {
			return err
		}
		if err = eg.unmarshalEntryMapDefaults(entryMap, m); err != nil {
			return err
		}
	}

	return nil
}

func (eg *EntryGroup) unmarshalEntryMapDefaults(
	m, defaults *map[string]any,
) (err error) {
	entry := &Entry{}
	if err = entry.UnmarshalMapDefaults(m, defaults); err != nil {
		return err
	}

	eg.Entries = append(eg.Entries, entry)

	return nil
}

type Entry struct {
	Name            string
	Type            Type
	Os              ostypes.Os
	Arch            archtypes.Arch
	Tags            *set.Set[string]
	RequiredTags    *set.Set[string]
	ContinueOnError bool
	RetryCount      int
	RetryBehavior   RetryBehavior
	commander       EntryCommander
}

func NewEntry(name string, commander EntryCommander) *Entry {
	return &Entry{
		Name:         name,
		Tags:         set.New[string](),
		RequiredTags: set.New[string](),
		commander:    commander,
	}
}

func (e *Entry) GetOs() ostypes.Os {
	return e.Os
}

func (e *Entry) GetArch() archtypes.Arch {
	return e.Arch
}

func (e *Entry) GetTags() *set.Set[string] {
	return e.Tags
}

func (e *Entry) GetRequiredTags() *set.Set[string] {
	return e.RequiredTags
}

func (e *Entry) UnmarshalMapDefaults(
	m, defaults *map[string]any,
) (err error) {
	var exists bool

	var name *string
	name, exists, err = parse.GetDefaultMap[string](m, "name", defaults)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("setup entry name field is required")
	} else {
		e.Name = *name
	}

	if e.Type, err = inferEntryTypeDefault(m, defaults); err != nil {
		return err
	}

	if e.Os, _, err = parse.OsGetDefaultMap(m, "os", defaults); err != nil {
		return err
	}

	e.Arch, _, err = parse.ArchGetDefaultMap(m, "arch", defaults)
	if err != nil {
		return err
	}

	e.Tags, e.RequiredTags, _, err = parse.TagsGetDefaultMap(
		m,
		"tags",
		defaults,
	)
	if err != nil {
		return err
	}

	var continueOnError *bool
	continueOnError, exists, err = parse.GetDefaultMap[bool](
		m,
		"continueOnError",
		defaults,
	)
	if err != nil {
		return err
	} else if exists {
		e.ContinueOnError = *continueOnError
	}

	if err = e.unmarshalRetryMapDefaults(m, defaults); err != nil {
		return err
	}

	var commander EntryCommanderUnmarshaler
	if commander, err = newEntryCommander(e.Type); err != nil {
		return err
	}
	if err = commander.unmarshalMapDefaults(m, defaults); err != nil {
		return err
	}
	e.commander = commander

	return nil
}

func (e *Entry) unmarshalRetryMapDefaults(
	m, defaults *map[string]any,
) (err error) {
	var a, aDefaults *any

	if a, _, err = parse.Get[any](m, "retry"); err != nil {
		return err
	}

	if defaults != nil {
		if aDefaults, _, err = parse.Get[any](defaults, "retry"); err != nil {
			return err
		}
	}

	return e.unmarshalRetryAnyDefaults(a, aDefaults)
}

func (e *Entry) unmarshalRetryAnyDefaults(a, defaults *any) (err error) {
	var (
		exists      bool
		m           *map[string]any
		defaultsMap *map[string]any
	)

	if m, err = castRetry(a); err != nil {
		return err
	}

	if defaultsMap, err = castRetry(a); err != nil {
		return err
	}

	var retryCount *int
	retryCount, exists, err = parse.GetDefaultMap[int](m, "count", defaultsMap)
	if err != nil {
		return err
	} else if exists {
		e.RetryCount = *retryCount
	}

	e.RetryBehavior, _, err = retryBehaviorGetDefaultMap(
		m,
		"behavior",
		defaultsMap,
	)
	if err != nil {
		return err
	}

	return nil
}

func castRetry(a *any) (m *map[string]any, err error) {
	if a == nil {
		return &map[string]any{}, nil
	}

	switch a := (*a).(type) {
	case int:
		return &map[string]any{"count": a}, nil
	case map[string]any:
		return &a, nil
	default:
		return nil,
			fmt.Errorf(
				"setup retry value must be of type string or map[string]any but got %T",
				a,
			)
	}
}
