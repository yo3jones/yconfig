package generate

import (
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

type TemplateContext struct {
	OS string
}

func generateTemplate(templateName, destinationName string) error {
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

func getRelativePath(root, name string) string {
	relativePath := strings.TrimPrefix(name, root)
	relativePath = strings.TrimPrefix(relativePath, "/")

	return relativePath
}

func getDestPath(templateRoot, template, destinationRoot string) string {
	relativePath := getRelativePath(templateRoot, template)
	return path.Join(destinationRoot, relativePath)
}

func prepareDestination(destinationName string) error {
	var err error

	err = os.MkdirAll(path.Dir(destinationName), 0755)
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
