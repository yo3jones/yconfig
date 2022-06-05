package install

type Installer interface {
	Install() error
}

type installer struct {
	config any
	inst   *Install
}

func (ir *installer) Install() error {
	var (
		inst *Install
		err  error
	)

	if inst, err = Parse(ir.config); err != nil {
		return err
	}

	ir.inst = inst

	Print(inst)

	return nil
}

func New(config any) Installer {
	return &installer{
		config: config,
	}
}
