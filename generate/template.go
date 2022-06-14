package generate

import (
	"os"
	"runtime"
	"strings"
	"text/template"
)

type TemplateContext struct {
	OS   string
	Tags map[string]bool
}

func (c *TemplateContext) ForTag(tag string) bool {
	_, exists := c.Tags[strings.ToLower(tag)]
	return exists
}

func (c *TemplateContext) NotForTag(tag string) bool {
	return !c.ForTag(tag)
}

func generateTemplate(
	templateName, destinationName string,
	tags map[string]bool,
) error {
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
		OS:   runtime.GOOS,
		Tags: tags,
	}

	err3 := t.Execute(f, context)
	if err3 != nil {
		return err3
	}

	return nil
}

func prepareDestination(destinationName string) error {
	var err error

	err = makeDirAll(destinationName)
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
