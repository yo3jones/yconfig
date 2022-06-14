package setup

import (
	"fmt"
	"runtime"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/set"
)

type Filterable interface {
	Oser
	Archer
	Tagger
}

type Filterer interface {
	Tags(tags *set.Set[string]) Filterer
	EntryNames(entryNames *set.Set[string]) Filterer
	FilterSystemScripts(scripts []*Script) (systemScript *Script, err error)
	FilterSystemPackageManagers(
		packageManagers []*PackageManager,
	) (systemPackageManager *PackageManager, err error)
	FilterValues(setup *Setup) (values []Value, err error)
}

type filterer struct {
	runtimeOs         ostypes.Os
	runtimeArch       archtypes.Arch
	runtimeTags       *set.Set[string]
	runtimeEntryNames *set.Set[string]
	initialized       bool
}

const (
	restrictive    = true
	nonRestrictive = false
)

func NewFilterer() Filterer {
	return &filterer{initialized: false}
}

func (f *filterer) Tags(tags *set.Set[string]) Filterer {
	f.runtimeTags = tags
	f.initialized = false
	return f
}

func (f *filterer) EntryNames(entryNames *set.Set[string]) Filterer {
	f.runtimeEntryNames = entryNames
	return f
}

func (f *filterer) FilterSystemScripts(
	scripts []*Script,
) (systemScript *Script, err error) {
	foundScript, found, err := filter(f, scripts, nonRestrictive)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, fmt.Errorf("no system script found")
	} else {
		return foundScript, nil
	}
}

func (f *filterer) FilterSystemPackageManagers(
	packageManagers []*PackageManager,
) (systemPackageManager *PackageManager, err error) {
	foundPackageManager, found, err := filter(
		f,
		packageManagers,
		nonRestrictive,
	)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, fmt.Errorf("no system package manager found")
	} else {
		return foundPackageManager, nil
	}
}

func (f *filterer) FilterValues(setup *Setup) (values []Value, err error) {
	values = make([]Value, 0, len(setup.Entries))

	hasEntryNames := f.runtimeEntryNames.Len() > 0
	for _, entry := range setup.Entries {
		if hasEntryNames && !f.runtimeEntryNames.Contains(entry.Name) {
			continue
		}

		foundValue, found, err := filter(f, entry.Values, restrictive)
		if err != nil {
			return nil, err
		} else if found {
			values = append(values, foundValue)
		}
	}

	return values, nil
}

func (f *filterer) initialize() (err error) {
	if f.initialized {
		return nil
	}

	if f.runtimeOs, err = ostypes.OsFromString(runtime.GOOS); err != nil {
		return err
	}

	f.runtimeArch, err = archtypes.ArchFromString(runtime.GOARCH)
	if err != nil {
		return err
	}

	return nil
}

func filter[T Filterable, E ~[]T](
	f *filterer,
	items E,
	restrictive bool,
) (foundItem T, found bool, err error) {
	if err = f.initialize(); err != nil {
		return foundItem, false, err
	}

	compatibleItems := make([]T, 0, len(items))
	tagMatchedItems := make([]T, 0, len(items))
	requiredTagMatchedItems := make([]T, 0, len(items))
	for i := range items {
		item := items[i]

		keepCompatible, keepTag, keepRequiredTag := shouldKeep(
			f,
			item,
			restrictive,
		)

		if keepCompatible {
			compatibleItems = append(compatibleItems, item)
		}

		if keepTag {
			tagMatchedItems = append(tagMatchedItems, item)
		}

		if keepRequiredTag {
			requiredTagMatchedItems = append(requiredTagMatchedItems, item)
		}
	}

	if len(requiredTagMatchedItems) > 0 {
		return requiredTagMatchedItems[0], true, nil
	}

	if len(tagMatchedItems) > 0 {
		return tagMatchedItems[0], true, nil
	}

	if len(compatibleItems) > 0 {
		return compatibleItems[0], true, nil
	}

	return foundItem, false, nil
}

func shouldKeep[T Filterable](
	f *filterer,
	item T,
	restrictive bool,
) (keepCompatible, keepTag, keepRequiredTag bool) {
	compatible,
		hasTags,
		tagMatch,
		hasRequiredTags,
		requiredTagMatch,
		hasRuntimeTags := getMatches(f, item)

	if !compatible {
		return false, false, false
	}

	if hasRequiredTags && !requiredTagMatch {
		return false, false, false
	}

	keepCompatible = !restrictive
	keepTag = !hasRuntimeTags || (hasTags && tagMatch)
	keepRequiredTag = hasRequiredTags && requiredTagMatch && keepTag

	return keepCompatible, keepTag, keepRequiredTag
}

func getMatches[T Filterable](
	f *filterer,
	item T,
) (
	compatible,
	hasTags,
	tagMatch,
	hasRequiredTags,
	requiredTagMatch,
	hasRuntimeTags bool,
) {
	compatible = isCompatible(f, item)
	hasTags = item.GetTags().Len() > 0
	tagMatch = isTagMatch(f, item)
	hasRequiredTags = item.GetRequiredTags().Len() > 0
	requiredTagMatch = isRequiredTagMatch(f, item)
	hasRuntimeTags = f.runtimeTags.Len() > 0

	return compatible,
		hasTags,
		tagMatch,
		hasRequiredTags,
		requiredTagMatch,
		hasRuntimeTags
}

func isCompatible[T Filterable](f *filterer, item T) bool {
	if item.GetOs() != ostypes.Any && item.GetOs() != f.runtimeOs {
		return false
	}

	if item.GetArch() != archtypes.Any && item.GetArch() != f.runtimeArch {
		return false
	}

	return true
}

func isTagMatch[T Filterable](f *filterer, item T) bool {
	return f.runtimeTags.ContainedIn(item.GetTags())
}

func isRequiredTagMatch[T Filterable](f *filterer, item T) bool {
	return item.GetRequiredTags().ContainedIn(f.runtimeTags)
}
