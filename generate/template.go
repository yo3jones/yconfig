package generate

import (
	"os"
	"path"
	"runtime"
	"text/template"
)

type TemplateContext struct {
	OS string
}

func generateTemplate(templateName, destinationName string) error {
	if err := prepareDestination(destinationName); err != nil {
		return err
	}

	t, err1 := template.ParseFiles(templateName)
	if err1 != nil {
		return err1
	}

	f, err2 := os.Create(destinationName)
	if err2 != nil {
		return err2
	}
	defer f.Close()

	context := &TemplateContext{
		OS: runtime.GOOS,
	}

	err3 := t.Execute(f, context)
	if err3 != nil {
		return err3
	}

	return nil
}

func prepareDestination(destinationName string) error {
	var err error

	err = makeDirAll(path.Dir(destinationName))
	if err != nil {
		return err
	}

	_, err = os.Stat(destinationName)
	if err == nil {
		err = os.Remove(destinationName)
		if err != nil {
			return err
		}
	}

	return nil
}
