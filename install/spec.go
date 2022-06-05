package install

type Install struct {
	Groups []Group
}

type Group struct {
	Name     string
	Commands []Command
	Os       string
	Arch     string
}

type Command struct {
	Command string
	Os      string
	Arch    string
}
