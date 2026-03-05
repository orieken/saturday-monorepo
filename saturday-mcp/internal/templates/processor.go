package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

// Processor handles template processing and execution
type Processor struct {
	registry *Registry
	cache    *Cache
}

// NewProcessor creates a new template processor
func NewProcessor(registry *Registry, cache *Cache) *Processor {
	return &Processor{
		registry: registry,
		cache:    cache,
	}
}

// Process executes a template with the given data
func (p *Processor) Process(templateName string, data interface{}) (string, error) {
	// Check cache first
	if p.cache != nil {
		if cached, found := p.cache.Get(templateName, data); found {
			return cached, nil
		}
	}

	// Get template from registry
	tmpl, err := p.registry.Get(templateName)
	if err != nil {
		return "", err
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	result := buf.String()

	// Cache the result
	if p.cache != nil {
		p.cache.Set(templateName, data, result)
	}

	return result, nil
}

// ProcessWithFuncs executes a template with custom functions
func (p *Processor) ProcessWithFuncs(templateName string, data interface{}, funcMap template.FuncMap) (string, error) {
	tmpl, err := p.registry.Get(templateName)
	if err != nil {
		return "", err
	}

	// Clone template and add custom functions
	cloned, err := tmpl.Clone()
	if err != nil {
		return "", fmt.Errorf("failed to clone template %s: %w", templateName, err)
	}

	cloned = cloned.Funcs(funcMap)

	// Execute template
	var buf bytes.Buffer
	if err := cloned.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// Validate checks if a template can be executed without errors
func (p *Processor) Validate(templateName string, data interface{}) error {
	_, err := p.Process(templateName, data)
	return err
}
