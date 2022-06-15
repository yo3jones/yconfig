package setup

import (
	"fmt"

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
	retryCount      *int
	retryBehavior   *RetryBehavior
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

	if err = parseRetry(config, cf); err != nil {
		return nil, err
	}

	return cf, nil
}

func parseRetry(config *map[string]any, cf *commonFields) (err error) {
	var (
		exists      bool
		retryConfig *any
	)

	if retryConfig, exists, err = parse.Get[any](config, "retry"); err != nil {
		return err
	} else if !exists {
		return nil
	}

	switch retryConfig := (*retryConfig).(type) {
	case int:
		cf.retryCount = &retryConfig
		return nil
	case map[string]any:
		return parseRetryMap(&retryConfig, cf)
	default:
		return fmt.Errorf(
			"setup retry must be either a number or map[string]any but got %T",
			retryConfig,
		)
	}
}

func parseRetryMap(config *map[string]any, cf *commonFields) (err error) {
	if cf.retryCount, _, err = parse.Get[int](config, "count"); err != nil {
		return err
	}

	cf.retryBehavior, _, err = retryBehaviorPtrGet(config, "behavior")
	if err != nil {
		return err
	}

	return nil
}
