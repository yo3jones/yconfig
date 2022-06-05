package install

import (
	"encoding/json"
	"fmt"
)

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

func Print(install *Install) {
	jsonBytes, err := json.MarshalIndent(install, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonBytes)
}
