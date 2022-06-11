package setup

import (
	"fmt"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/parse"
)

func ParsePackageManagers(config *any) (pms []PackageManager, err error) {
	var configSlice *[]any

	if configSlice, err = parse.Cast[[]any](config); err != nil {
		return nil, err
	}

	return parsePackageManagers(*configSlice)
}

func parsePackageManagers(config []any) (pms []PackageManager, err error) {
	pms = make([]PackageManager, len(config))

	for i := range config {
		var (
			configMap *map[string]any
			pm        *PackageManager
		)
		configElement := config[i]

		configMap, err = parse.Cast[map[string]any](&configElement)
		if err != nil {
			return nil, err
		}

		if pm, err = parsePackageManager(*configMap); err != nil {
			return nil, err
		}

		pms[i] = *pm
	}

	return pms, nil
}

func parsePackageManager(
	config map[string]any,
) (pm *PackageManager, err error) {
	var (
		exists bool
		os     ostypes.Os
		arch   archtypes.Arch
		tags   map[string]bool
		script *string
	)

	if os, _, err = parse.OsGet(&config, "os"); err != nil {
		return nil, err
	}

	if arch, _, err = parse.ArchGet(&config, "arch"); err != nil {
		return nil, err
	}

	if tags, _, _, err = parse.TagsGet(&config, "tags"); err != nil {
		return nil, err
	}

	if script, exists, err = parse.Get[string](&config, "script"); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("packageManager requires a script field")
	}

	pm = &PackageManager{
		Os:     os,
		Arch:   arch,
		Tags:   tags,
		Script: *script,
	}

	return pm, nil
}
