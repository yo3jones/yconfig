package setup

import "github.com/yo3jones/yconfig/parse"

func typeGet(m *map[string]any, key string) (t Type, exists bool, err error) {
	var str *string

	if str, exists, err = parse.Get[string](m, key); err != nil {
		return t, exists, err
	} else if !exists {
		return t, false, nil
	}

	if t, err = TypeFromString(*str); err != nil {
		return t, true, err
	}

	return t, true, err
}

func retryBehaviorGet(
	m *map[string]any,
	key string,
) (behavior RetryBehavior, exists bool, err error) {
	var (
		str           *string
		retryBehavior RetryBehavior
	)

	if str, exists, err = parse.Get[string](m, key); err != nil {
		return behavior, false, err
	} else if !exists {
		return behavior, false, nil
	}

	if behavior, err = RetryBehaviorFromString(*str); err != nil {
		return retryBehavior, true, err
	}

	return behavior, true, nil
}

func retryBehaviorGetDefaultMap(
	m *map[string]any,
	key string,
	defaults *map[string]any,
) (behavior RetryBehavior, exists bool, err error) {
	if behavior, exists, err = retryBehaviorGet(m, key); err != nil {
		return behavior, false, err
	} else if exists {
		return behavior, true, nil
	}

	if defaults == nil {
		return behavior, false, nil
	}

	return retryBehaviorGet(defaults, key)
}
