package generate

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

type Generator interface {
	TemplateRoot(templateRoot string) Generator
	Include(include []string) Generator
	Exclude(exclude []string) Generator
	Progress(progress func(progress *Progress)) Generator
	Generate() error
}

type TemplateProgress struct {
	Path     string
	Progress string
}

type Progress struct {
	TemplatesProgress []TemplateProgress
}

type generator struct {
	templateRoot string
	include      []string
	exclude      []string
	progress     func(progress *Progress)
}

func (g *generator) TemplateRoot(templateRoot string) Generator {
	g.templateRoot = templateRoot
	return g
}

func (g *generator) Include(include []string) Generator {
	g.include = include
	return g
}

func (g *generator) Exclude(exclude []string) Generator {
	g.exclude = exclude
	return g
}

func (g *generator) Progress(progress func(progress *Progress)) Generator {
	g.progress = progress
	return g
}

func (g *generator) prepare() {
	if g.progress == nil {
		g.progress = func(_ *Progress) {}
	}
}

func (g *generator) Generate() error {
	g.prepare()

	templates := glob(g.templateRoot, g.include, g.exclude)

	progress := &Progress{[]TemplateProgress{}}
	for _, template := range templates {
		progress.TemplatesProgress = append(
			progress.TemplatesProgress,
			TemplateProgress{template, "waiting"},
		)
	}

	g.progress(progress)

	return nil
}

func New() Generator {
	return &generator{}
}

func generate() {
	type Foo struct {
		Bar string
	}

	foo := Foo{"bar"}

	t, err1 := template.New("config").
		Parse("something here {{- .Bar }} something after")
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	writer := bytes.NewBuffer([]byte{})

	t.Execute(writer, foo)

	fmt.Printf("%s\n", writer.Bytes())
}

func Generate(glob string) {
	_, err1 := filepath.Glob(glob)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	// for _, m := range matches {
	// 	fmt.Println(m)
	// }
}
