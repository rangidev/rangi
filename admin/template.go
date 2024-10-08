package admin

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
	"github.com/rangidev/rangi/collection"
	"github.com/rangidev/rangi/config"
)

const (
	baseTemplateName = "base.html"

	keyInternalTemplateName   = "_internalTemplateName"
	keyInternalAllCollections = "_internalAllCollections"
)

var (
	//go:embed template
	templateFS embed.FS

	TemplateLogin      = &TemplateDefinition{name: "login.html", dependencies: []string{baseTemplateName}}
	TemplateDashboard  = &TemplateDefinition{name: "dashboard.html", dependencies: []string{baseTemplateName, "navbar.html"}}
	TemplateCollection = &TemplateDefinition{name: "collection.html", dependencies: []string{baseTemplateName, "navbar.html"}}
	TemplateEdit       = &TemplateDefinition{name: "edit.html", dependencies: []string{baseTemplateName, "navbar.html"}}
	TemplateSettings   = &TemplateDefinition{name: "settings.html", dependencies: []string{baseTemplateName, "navbar.html"}}
)

type TemplateDefinition struct {
	name         string
	dependencies []string
}

type Templates struct {
	config     *config.Config
	templates  map[string]*template.Template
	templateFS fs.FS
}

type TemplateData map[string]interface{}

type TemplateDataCollection struct {
	TemplateData
	Collection string
}

func NewTemplates(config *config.Config) (*Templates, error) {
	t := Templates{
		config:    config,
		templates: make(map[string]*template.Template),
	}
	if t.config.EnableTemplateDevelopment {
		// In development mode, we want to read from filesystem
		t.templateFS = os.DirFS(filepath.Join(config.ExecutableDir, "admin/template"))
	} else {
		// Read embedded templates
		templateSubFS, err := fs.Sub(templateFS, "template")
		if err != nil {
			return nil, fmt.Errorf("could not get filesystem sub directory: %v", err)
		}
		t.templateFS = templateSubFS
	}
	return &t, nil
}

// Render
// subTemplateName defines the sub template that should be executed (e. g. a "block" defined in the template string to render only a part of the original template). May be an empty string.
func (t *Templates) Render(w http.ResponseWriter, data TemplateData, templateDef *TemplateDefinition, collectionLoader *collection.CollectionLoader, subTemplateName string) error {
	// Always re-read templates in development mode
	// Otherwise cache the templates
	if _, ok := t.templates[templateDef.name]; !ok || t.config.EnableTemplateDevelopment {
		tmpl, err := template.New("").Funcs(sprig.FuncMap()).ParseFS(t.templateFS, append(templateDef.dependencies, templateDef.name)...)
		if err != nil {
			return fmt.Errorf("could not parse template %s: %v", templateDef.name, err)
		}
		t.templates[templateDef.name] = tmpl
	}
	if data == nil {
		data = TemplateData{}
	}
	allCollections, err := collectionLoader.GetAll()
	if err != nil {
		return fmt.Errorf("could not get all collections: %v", err)
	}
	data[keyInternalAllCollections] = allCollections
	data[keyInternalTemplateName] = templateDef.name
	if subTemplateName != "" {
		return t.templates[templateDef.name].ExecuteTemplate(w, subTemplateName, data)
	}
	return t.templates[templateDef.name].ExecuteTemplate(w, baseTemplateName, data)
}
