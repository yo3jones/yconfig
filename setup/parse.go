package setup

import "github.com/yo3jones/yconfig/parse"

func Parse(config *any) (*Setup, error) {
	var (
		entriesConfig *[]any
		entries       []Entry
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

func typeGet(
	obj *map[string]any,
	key string,
) (t Type, exists bool, err error) {
	var str *string

	if str, exists, err = parse.Get[string](obj, key); err != nil {
		return t, exists, err
	} else if !exists {
		return t, false, nil
	}

	if t, err = TypeFromString(*str); err != nil {
		return t, true, err
	}

	return t, true, err
}
