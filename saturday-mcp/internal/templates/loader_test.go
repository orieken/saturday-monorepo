package templates

import (
	"testing"
)

func TestLoader(t *testing.T) {
	registry := NewRegistry()
	loader := NewLoader(registry)

	t.Run("LoadAll", func(t *testing.T) {
		err := loader.LoadAll()
		if err != nil {
			t.Fatalf("Failed to load templates: %v", err)
		}

		// Check that templates were loaded
		names := registry.List()
		if len(names) == 0 {
			t.Error("Expected templates to be loaded")
		}

		// Verify specific templates exist
		expectedTemplates := []string{"page", "site"}
		for _, name := range expectedTemplates {
			if !registry.Has(name) {
				t.Errorf("Expected template %q to be loaded", name)
			}
		}
	})

	t.Run("LoadTemplate", func(t *testing.T) {
		newRegistry := NewRegistry()
		newLoader := NewLoader(newRegistry)

		err := newLoader.LoadTemplate("data/page.tmpl")
		if err != nil {
			t.Fatalf("Failed to load template: %v", err)
		}

		if !newRegistry.Has("page") {
			t.Error("Expected page template to be loaded")
		}
	})

	t.Run("LoadTemplateByType", func(t *testing.T) {
		newRegistry := NewRegistry()
		newLoader := NewLoader(newRegistry)

		err := newLoader.LoadTemplateByType(TemplateTypePage)
		if err != nil {
			t.Fatalf("Failed to load templates by type: %v", err)
		}

		if !newRegistry.Has("page") {
			t.Error("Expected page template to be loaded")
		}
	})
}
