package setup

import (
	"fmt"
	"strings"

	"github.com/yo3jones/yconfig/archtypes"
	"github.com/yo3jones/yconfig/ostypes"
	"github.com/yo3jones/yconfig/parse"
	"github.com/yo3jones/yconfig/set"
)

type System interface {
	PackageManager() *SystemPackageManager
	Script() *SystemScript
}

type SystemPackageManager struct {
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         *set.Set[string]
	RequiredTags *set.Set[string]
	Script       string
}

func UnmarshalSystemPackageManagers(
	a *any,
) (pms []*SystemPackageManager, err error) {
	var slicePtr *[]any
	if slicePtr, err = parse.Cast[[]any](a); err != nil {
		return nil, err
	}
	slice := *slicePtr

	pms = make([]*SystemPackageManager, len(slice))
	for i, pmAny := range slice {
		var pmMap *map[string]any
		if pmMap, err = parse.Cast[map[string]any](&pmAny); err != nil {
			return nil, err
		}
		pm := &SystemPackageManager{}
		if err = pm.UnmarshalMap(pmMap); err != nil {
			return nil, err
		}
		pms[i] = pm
	}

	return pms, nil
}

func (pm *SystemPackageManager) GetOs() ostypes.Os {
	return pm.Os
}

func (pm *SystemPackageManager) GetArch() archtypes.Arch {
	return pm.Arch
}

func (pm *SystemPackageManager) GetTags() *set.Set[string] {
	return pm.Tags
}

func (pm *SystemPackageManager) GetRequiredTags() *set.Set[string] {
	return pm.RequiredTags
}

func (pm *SystemPackageManager) BuildCommand(
	script *SystemScript,
	packages []string,
) (cmd string, args []string) {
	pmScriptParts := make([]string, 0, len(packages)+1)
	pmScriptParts = append(pmScriptParts, pm.Script)
	pmScriptParts = append(pmScriptParts, packages...)

	pmScript := strings.Join(pmScriptParts, " ")

	return script.BuildCommand(pmScript)
}

func (pm *SystemPackageManager) UnmarshalMap(m *map[string]any) (err error) {
	if pm.Os, _, err = parse.OsGet(m, "os"); err != nil {
		return err
	}
	if pm.Arch, _, err = parse.ArchGet(m, "arch"); err != nil {
		return err
	}
	if pm.Tags, pm.RequiredTags, _, err = parse.TagsGet(m, "tags"); err != nil {
		return err
	}
	var script *string
	if script, _, err = parse.Get[string](m, "script"); err != nil {
		return err
	}
	pm.Script = *script
	return nil
}

type SystemScript struct {
	Os           ostypes.Os
	Arch         archtypes.Arch
	Tags         *set.Set[string]
	RequiredTags *set.Set[string]
	Cmd          string
	Args         []string
}

func UnmarshalSystemScripts(a *any) (ss []*SystemScript, err error) {
	var slicePtr *[]any
	if slicePtr, err = parse.Cast[[]any](a); err != nil {
		return nil, err
	}
	slice := *slicePtr

	ss = make([]*SystemScript, len(slice))
	for i, sAny := range slice {
		var sMap *map[string]any
		if sMap, err = parse.Cast[map[string]any](&sAny); err != nil {
			return nil, err
		}
		s := &SystemScript{}
		if err = s.UnmarshalMap(sMap); err != nil {
			return nil, err
		}
		ss[i] = s
	}

	return ss, nil
}

func (s *SystemScript) GetOs() ostypes.Os {
	return s.Os
}

func (s *SystemScript) GetArch() archtypes.Arch {
	return s.Arch
}

func (s *SystemScript) GetTags() *set.Set[string] {
	return s.Tags
}

func (s *SystemScript) GetRequiredTags() *set.Set[string] {
	return s.RequiredTags
}

func (s *SystemScript) BuildCommand(script string) (cmd string, args []string) {
	cmd = s.Cmd

	args = make([]string, 0, len(s.Args)+1)
	args = append(args, s.Args...)
	args = append(args, fmt.Sprintf("\\\n%s", script))

	return cmd, args
}

func (s *SystemScript) UnmarshalMap(m *map[string]any) (err error) {
	if s.Os, _, err = parse.OsGet(m, "os"); err != nil {
		return err
	}
	if s.Arch, _, err = parse.ArchGet(m, "arch"); err != nil {
		return err
	}
	if s.Tags, s.RequiredTags, _, err = parse.TagsGet(m, "tags"); err != nil {
		return err
	}
	var cmd *string
	if cmd, _, err = parse.Get[string](m, "cmd"); err != nil {
		return err
	}
	s.Cmd = *cmd
	if s.Args, _, err = parse.StringSliceGet(m, "args"); err != nil {
		return err
	}
	return nil
}
