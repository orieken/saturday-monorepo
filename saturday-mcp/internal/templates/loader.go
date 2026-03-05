package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"
)

//go:embed data/*.tmpl
var templateFS embed.FS

// Loader handles loading templates from embedded filesystem
type Loader struct {
	registry *Registry
	funcMap  template.FuncMap
}

// NewLoader creates a new template loader
func NewLoader(registry *Registry) *Loader {
	return &Loader{
		registry: registry,
		funcMap:  DefaultFuncMap(),
	}
}

// LoadAll loads all templates from the embedded filesystem
func (l *Loader) LoadAll() error {
	return fs.WalkDir(templateFS, "data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		return l.LoadTemplate(path)
	})
}

// LoadTemplate loads a single template file
func (l *Loader) LoadTemplate(path string) error {
	content, err := templateFS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", path, err)
	}

	name := filepath.Base(path)
	name = name[:len(name)-len(filepath.Ext(name))] // Remove .tmpl extension

	tmpl, err := template.New(name).Funcs(l.funcMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	return l.registry.Register(name, tmpl)
}

// LoadTemplateByType loads templates for a specific type
func (l *Loader) LoadTemplateByType(templateType TemplateType) error {
	pattern := fmt.Sprintf("data/%s*.tmpl", templateType)
	matches, err := fs.Glob(templateFS, pattern)
	if err != nil {
		return fmt.Errorf("failed to find templates for type %s: %w", templateType, err)
	}

	for _, match := range matches {
		if err := l.LoadTemplate(match); err != nil {
			return err
		}
	}

	return nil
}

// GetRawContent returns the raw content of a template by name
func (l *Loader) GetRawContent(name string) ([]byte, error) {
	path := fmt.Sprintf("data/%s.tmpl", name)
	content, err := templateFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", name, err)
	}
	return content, nil
}
