package parse

import (
	"fmt"
	"strings"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/set"
)

func Cast[T any](obj *any) (*T, error) {
	switch obj := (*obj).(type) {
	case T:
		return &obj, nil
	}

	return nil, fmt.Errorf("error casting type %T", obj)
}

func StringSliceCast(obj *any) (strSlice []string, err error) {
	var slice *[]any

	if slice, err = Cast[[]any](obj); err != nil {
		return nil, err
	}

	strSlice = make([]string, len(*slice))
	for i, val := range *slice {
		var str *string

		if str, err = Cast[string](&val); err != nil {
			return nil, err
		}

		strSlice[i] = *str
	}

	return strSlice, nil
}

func Get[T any](
	obj *map[string]any,
	key string,
) (value *T, exists bool, err error) {
	rawVal, exists := (*obj)[key]
	if !exists {
		return nil, false, nil
	}

	if value, err = Cast[T](&rawVal); err != nil {
		return nil, true, err
	}

	return value, true, nil
}

func StringSliceGet(
	obj *map[string]any,
	key string,
) (strs []string, exists bool, err error) {
	var val *any

	if val, exists, err = Get[any](obj, key); err != nil {
		return nil, exists, err
	} else if !exists {
		return []string{}, false, nil
	}

	switch castedVal := (*val).(type) {
	case string:
		strs = make([]string, 1)
		strs[0] = castedVal
	case []any:
		if strs, err = StringSliceCast(val); err != nil {
			return nil, true, err
		}
	default:
		return nil,
			true,
			fmt.Errorf("expected either a string or []string but got %T", val)
	}

	return strs, true, nil
}

func OsGet(
	obj *map[string]any,
	key string,
) (os ostypes.Os, exists bool, err error) {
	var osPtr *ostypes.Os

	if osPtr, exists, err = OsPtrGet(obj, key); err != nil {
		return os, false, err
	} else if !exists {
		return os, false, nil
	} else {
		return *osPtr, true, nil
	}
}

func OsPtrGet(
	obj *map[string]any,
	key string,
) (osPtr *ostypes.Os, exists bool, err error) {
	var (
		str *string
		os  ostypes.Os
	)

	if str, exists, err = Get[string](obj, key); err != nil {
		return nil, exists, err
	} else if !exists {
		return nil, false, nil
	}

	if os, err = ostypes.OsFromString(*str); err != nil {
		return nil, true, err
	}

	return &os, true, err
}

func ArchGet(
	obj *map[string]any,
	key string,
) (arch archtypes.Arch, exists bool, err error) {
	var archPtr *archtypes.Arch

	if archPtr, exists, err = ArchPtrGet(obj, key); err != nil {
		return arch, false, err
	} else if !exists {
		return arch, false, nil
	} else {
		return *archPtr, true, nil
	}
}

func ArchPtrGet(
	obj *map[string]any,
	key string,
) (archPtr *archtypes.Arch, exists bool, err error) {
	var (
		str  *string
		arch archtypes.Arch
	)

	if str, exists, err = Get[string](obj, key); err != nil {
		return nil, false, err
	} else if !exists {
		return nil, false, nil
	}

	if arch, err = archtypes.ArchFromString(*str); err != nil {
		return nil, true, err
	}

	return &arch, true, err
}

func TagsGet(
	obj *map[string]any,
	key string,
) (tags, requiredTags *set.Set[string], exists bool, err error) {
	var rawTags []string

	tags = set.New[string]()
	requiredTags = set.New[string]()

	if rawTags, exists, err = StringSliceGet(obj, key); err != nil {
		return nil, nil, exists, err
	} else if !exists {
		return tags, requiredTags, false, nil
	}

	for _, tag := range rawTags {
		if strings.HasSuffix(tag, "!") {
			normalizedTag := tag[:len(tag)-1]
			requiredTags.Put(normalizedTag)
			tags.Put(normalizedTag)
		} else {
			tags.Put(tag)
		}
	}

	return tags, requiredTags, true, nil
}
