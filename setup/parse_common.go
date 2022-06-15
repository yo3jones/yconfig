package setup

import (
	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/parse"
	"github.com/yo3jones/yconfig/set"
)

type commonFields struct {
	t               *Type
	os              *ostypes.Os
	arch            *archtypes.Arch
	tags            *set.Set[string]
	requiredTags    *set.Set[string]
	continueOnError *bool
}

func parseCommonFields(config *map[string]any) (cf *commonFields, err error) {
	cf = &commonFields{}

	if cf.t, _, err = typePtrGet(config, "type"); err != nil {
		return nil, err
	}

	if cf.os, _, err = parse.OsPtrGet(config, "os"); err != nil {
		return nil, err
	}

	if cf.arch, _, err = parse.ArchPtrGet(config, "arch"); err != nil {
		return nil, err
	}

	cf.tags, cf.requiredTags, _, err = parse.TagsGet(config, "tags")
	if err != nil {
		return nil, err
	}

	cf.continueOnError, _, err = parse.Get[bool](config, "continueOnError")
	if err != nil {
		return nil, err
	}

	return cf, nil
}
