package setup

import "fmt"

type Setuper interface {
	ScriptsConfig(scriptsConfig *any) Setuper
	PackageManagersConfig(packageManagersConfig *any) Setuper
	Config(config *any) Setuper
	Setup() (err error)
}

type setuper struct {
	scriptsConfig         *any
	packageManagersConfig *any
	config                *any
	scripts               []*Script
	packageManagers       []*PackageManager
	setup                 *Setup
}

func New() Setuper {
	return &setuper{}
}

func (s *setuper) ScriptsConfig(scriptsConfig *any) Setuper {
	s.scriptsConfig = scriptsConfig
	return s
}

func (s *setuper) PackageManagersConfig(packageManagersConfig *any) Setuper {
	s.packageManagersConfig = packageManagersConfig
	return s
}

func (s *setuper) Config(config *any) Setuper {
	s.config = config
	return s
}

func (s *setuper) Setup() (err error) {
	if err = s.parseConfigs(); err != nil {
		return err
	}

	s.setup.Print()
	fmt.Printf("\n Package Managers \n\n")
	SlicePrint(s.packageManagers)
	fmt.Printf("\n Scripts \n\n")
	SlicePrint(s.scripts)

	return nil
}

func (s *setuper) parseConfigs() (err error) {
	if s.scripts, err = ParseScripts(s.scriptsConfig); err != nil {
		return err
	}

	s.packageManagers, err = ParsePackageManagers(s.packageManagersConfig)
	if err != nil {
		return err
	}

	if s.setup, err = Parse(s.config); err != nil {
		return err
	}

	return nil
}
