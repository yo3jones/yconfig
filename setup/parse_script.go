package setup

import (
	"fmt"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/parse"
)

func ParseScripts(config *any) (scripts []*Script, err error) {
	var configSlice *[]any

	if configSlice, err = parse.Cast[[]any](config); err != nil {
		return nil, err
	}

	return parseScripts(*configSlice)
}

func parseScripts(config []any) (scripts []*Script, err error) {
	scripts = make([]*Script, len(config))

	for i := range config {
		var (
			configMap *map[string]any
			script    *Script
		)

		if configMap, err = parse.Cast[map[string]any](&config[i]); err != nil {
			return nil, err
		}

		if script, err = parseScript(*configMap); err != nil {
			return nil, err
		}

		scripts[i] = script
	}

	return scripts, nil
}

func parseScript(config map[string]any) (script *Script, err error) {
	var (
		exists bool
		os     ostypes.Os
		arch   archtypes.Arch
		tags   map[string]bool
		cmd    *string
		args   []string
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

	if cmd, exists, err = parse.Get[string](&config, "cmd"); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("cmd field is required for script elements")
	}

	if args, _, err = parse.StringSliceGet(&config, "args"); err != nil {
		return nil, err
	}

	script = &Script{
		Os:   os,
		Arch: arch,
		Tags: tags,
		Cmd:  *cmd,
		Args: args,
	}

	return script, nil
}
