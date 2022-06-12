package setup

import (
	"fmt"
	"runtime"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
)

type Filterable interface {
	Oser
	Archer
	Tagger
}

type Filterer interface {
	Tags(tags map[string]bool) Filterer
	FilterSystemScripts(scripts []*Script) (systemScript *Script, err error)
	FilterSystemPackageManagers(
		packageManagers []*PackageManager,
	) (systemPackageManager *PackageManager, err error)
	FilterValues(setup *Setup) (values []Value, err error)
}

type filterer struct {
	runtimeOs   ostypes.Os
	runtimeArch archtypes.Arch
	runtimeTags map[string]bool
	initialized bool
}

func NewFilterer() Filterer {
	return &filterer{initialized: false}
}

func (f *filterer) Tags(tags map[string]bool) Filterer {
	f.runtimeTags = tags
	f.initialized = false
	return f
}

func (f *filterer) FilterSystemScripts(
	scripts []*Script,
) (systemScript *Script, err error) {
	if foundScript, found, err := filter(f, scripts); err != nil {
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
	if foundPackageManager, found, err := filter(f, packageManagers); err != nil {
		return nil, err
	} else if !found {
		return nil, fmt.Errorf("no system package manager found")
	} else {
		return foundPackageManager, nil
	}
}

func (f *filterer) FilterValues(setup *Setup) (values []Value, err error) {
	values = make([]Value, 0, len(setup.Entries))

	for _, entry := range setup.Entries {
		if foundValue, found, err := filter(f, entry.Values); err != nil {
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

	if f.runtimeArch, err = archtypes.ArchFromString(runtime.GOARCH); err != nil {
		return err
	}

	return nil
}

func filter[T Filterable, E ~[]T](
	f *filterer,
	items E,
) (foundItem T, found bool, err error) {
	if err = f.initialize(); err != nil {
		return foundItem, false, err
	}

	compatibleItems := make([]T, 0, len(items))
	tagMatchedItems := make([]T, 0, len(items))
	requiredTagMatchedItems := make([]T, 0, len(items))
	for i := range items {
		item := items[i]

		requiredTagsMatch := containsAllTags(
			f.runtimeTags,
			item.GetRequiredTags(),
		)

		if !requiredTagsMatch {
			continue
		}

		if item.GetOs() != ostypes.Any && item.GetOs() != f.runtimeOs {
			continue
		}

		if item.GetArch() != archtypes.Any && item.GetArch() != f.runtimeArch {
			continue
		}

		compatibleItems = append(compatibleItems, item)

		if len(item.GetRequiredTags()) > 0 {
			requiredTagMatchedItems = append(requiredTagMatchedItems, item)
		}

		if len(f.runtimeTags) > 0 &&
			!containsAtLeastOneTag(f.runtimeTags, item.GetTags()) {
			continue
		}

		tagMatchedItems = append(tagMatchedItems, item)
	}

	if len(compatibleItems) <= 0 {
		return foundItem, false, nil
	}

	if len(requiredTagMatchedItems) > 0 {
		return requiredTagMatchedItems[0], true, nil
	}

	if len(tagMatchedItems) > 0 {
		return tagMatchedItems[0], true, nil
	}

	return compatibleItems[0], true, nil
}

func containsAtLeastOneTag(runtimeTags, tags map[string]bool) bool {
	for runtimeTag := range runtimeTags {
		if _, exists := tags[runtimeTag]; exists {
			return true
		}
	}
	return false
}

func containsAllTags(runtimeTags, tags map[string]bool) bool {
	for tag := range tags {
		if _, exists := runtimeTags[tag]; !exists {
			return false
		}
	}
	return true
}
