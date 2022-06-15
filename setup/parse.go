package setup

import "github.com/yo3jones/yconfig/parse"

func Parse(config *any) (*Setup, error) {
	var (
		entriesConfig *[]any
		entries       []*Entry
		err           error
	)

	if entriesConfig, err = parse.Cast[[]any](config); err != nil {
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

func typePtrGet(
	obj *map[string]any,
	key string,
) (typePtr *Type, exists bool, err error) {
	var (
		str *string
		t   Type
	)

	if str, exists, err = parse.Get[string](obj, key); err != nil {
		return nil, exists, err
	} else if !exists {
		return nil, false, nil
	}

	if t, err = TypeFromString(*str); err != nil {
		return nil, true, err
	}

	return &t, true, err
}

func retryBehaviorPtrGet(
	obj *map[string]any,
	key string,
) (retryBehaviorPtr *RetryBehavior, exists bool, err error) {
	var (
		str           *string
		retryBehavior RetryBehavior
	)

	if str, exists, err = parse.Get[string](obj, key); err != nil {
		return nil, false, err
	} else if !exists {
		return nil, false, nil
	}

	if retryBehavior, err = RetryBehaviorFromString(*str); err != nil {
		return nil, true, err
	}

	return &retryBehavior, true, nil
}
