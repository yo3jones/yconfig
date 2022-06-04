package generate

type Generator interface {
	TemplateRoot(templateRoot string) Generator
	DesinationRoot(destinationRoot string) Generator
	Include(include []string) Generator
	Exclude(exclude []string) Generator
	OnProgress(onProgress func(progress *Progress)) Generator
	Generate() error
}

type TemplateProgress struct {
	Path     string
	Progress string
}

type Progress struct {
	TemplatesProgress []*TemplateProgress
}

type generator struct {
	templateRoot    string
	destinationRoot string
	include         []string
	exclude         []string
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
			&TemplateProgress{template, "waiting"},
		)
	}
}

func (g *generator) notifyProgress(i int, newProgress string) {
	g.progress.TemplatesProgress[i].Progress = newProgress
	g.onProgress(g.progress)
}

func (g *generator) generateTemplate(i int) error {
	var err error
	templateName := g.templates[i]

	g.notifyProgress(i, "generating")

	destinationName := getDestPath(
		g.templateRoot,
		templateName,
		g.destinationRoot,
	)
	err = prepareDestination(destinationName)
	if err != nil {
		g.notifyProgress(i, "error")
		return err
	}

	err = generateTemplate(templateName, destinationName)
	if err != nil {
		g.notifyProgress(i, "error")
		return err
	}

	g.notifyProgress(i, "complete")

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
	return &generator{}
}
