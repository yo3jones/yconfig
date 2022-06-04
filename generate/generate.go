package generate

import (
	"path"
	"time"
)

type Generator interface {
	TemplateRoot(templateRoot string) Generator
	DesinationRoot(destinationRoot string) Generator
	Include(include []string) Generator
	Exclude(exclude []string) Generator
	Link(link bool) Generator
	Delay(delay int) Generator
	OnProgress(onProgress func(progress *Progress)) Generator
	Generate() error
}

type TemplateProgress struct {
	Path   string
	Status ProgressStatus
}

type Progress struct {
	TemplatesProgress []*TemplateProgress
}

type ProgressStatus int

const (
	Unknown ProgressStatus = iota
	Waiting
	Generating
	Linking
	Complete
	Error
)

func (ps ProgressStatus) String() string {
	switch ps {
	case Waiting:
		return "Waiting"
	case Generating:
		return "Generating"
	case Linking:
		return "Linking"
	case Complete:
		return "Complete"
	case Error:
		return "Error"
	}
	return "Uknown"
}

type generator struct {
	templateRoot    string
	destinationRoot string
	include         []string
	exclude         []string
	link            bool
	delay           int
	onProgress      func(progress *Progress)
	templates       []string
	progress        *Progress
}

func (g *generator) TemplateRoot(templateRoot string) Generator {
	g.templateRoot = templateRoot
	return g
}

func (g *generator) DesinationRoot(destinationRoot string) Generator {
	g.destinationRoot = destinationRoot
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

func (g *generator) Link(link bool) Generator {
	g.link = link
	return g
}

func (g *generator) Delay(delay int) Generator {
	g.delay = delay
	return g
}

func (g *generator) OnProgress(onProgress func(progress *Progress)) Generator {
	g.onProgress = onProgress
	return g
}

func (g *generator) prepare() {
	if g.onProgress == nil {
		g.onProgress = func(_ *Progress) {}
	}
}

func (g *generator) initTempalates() {
	g.templates = glob(g.templateRoot, g.include, g.exclude)
}

func (g *generator) initProgress() {
	g.progress = &Progress{[]*TemplateProgress{}}
	for _, template := range g.templates {
		g.progress.TemplatesProgress = append(
			g.progress.TemplatesProgress,
			&TemplateProgress{template, Waiting},
		)
	}
}

func (g *generator) notifyProgress(i int, newStatus ProgressStatus) {
	g.progress.TemplatesProgress[i].Status = newStatus
	g.onProgress(g.progress)
}

func (g *generator) sleep() {
	if g.delay <= 0 {
		return
	}
	time.Sleep(time.Duration(g.delay) * time.Millisecond)
}

func (g *generator) generateTemplate(i int) error {
	templateName := g.templates[i]

	g.sleep()

	g.notifyProgress(i, Generating)

	relativeName := getRelativePath(g.templateRoot, templateName)
	destinationName := path.Join(g.destinationRoot, relativeName)

	if err := generateTemplate(templateName, destinationName); err != nil {
		g.notifyProgress(i, Error)
		return err
	}

	if g.link {
		g.sleep()
		g.notifyProgress(i, Linking)
		if err := makeLink(relativeName, destinationName); err != nil {
			return err
		}
	}

	g.sleep()

	g.notifyProgress(i, Complete)

	return nil
}

func (g *generator) generateTemplates() error {
	var err error

	for i := range g.templates {
		err = g.generateTemplate(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) Generate() error {
	var err error

	g.prepare()

	g.initTempalates()
	g.initProgress()

	g.onProgress(g.progress)

	err = g.generateTemplates()
	if err != nil {
		return err
	}

	return nil
}

func New() Generator {
	return &generator{link: true}
}
