package templates

import (
	"fmt"
	"sync"
	"text/template"
)

// TemplateType represents the type of template
type TemplateType string

const (
	TemplateTypeSite    TemplateType = "site"
	TemplateTypePage    TemplateType = "page"
	TemplateTypeFlow    TemplateType = "flow"
	TemplateTypeElement TemplateType = "element"
	TemplateTypeService TemplateType = "service"
	TemplateTypeSteps   TemplateType = "steps"
)

// Registry manages template registration and retrieval
type Registry struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
}

// NewRegistry creates a new template registry
func NewRegistry() *Registry {
	return &Registry{
		templates: make(map[string]*template.Template),
	}
}

// Register adds a template to the registry
func (r *Registry) Register(name string, tmpl *template.Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[name]; exists {
		return fmt.Errorf("template %s already registered", name)
	}

	r.templates[name] = tmpl
	return nil
}

// Get retrieves a template by name
func (r *Registry) Get(name string) (*template.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tmpl, exists := r.templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s not found", name)
	}

	return tmpl, nil
}

// Has checks if a template exists
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.templates[name]
	return exists
}

// List returns all registered template names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.templates))
	for name := range r.templates {
		names = append(names, name)
	}
	return names
}

// Clear removes all templates from the registry
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.templates = make(map[string]*template.Template)
}
