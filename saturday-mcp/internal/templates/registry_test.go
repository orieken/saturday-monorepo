package templates

import (
	"testing"
	"text/template"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	t.Run("Register and Get", func(t *testing.T) {
		tmpl := template.Must(template.New("test").Parse("Hello {{.}}"))
		err := registry.Register("test", tmpl)
		if err != nil {
			t.Fatalf("Failed to register template: %v", err)
		}

		retrieved, err := registry.Get("test")
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Retrieved template is nil")
		}
	})

	t.Run("Register duplicate", func(t *testing.T) {
		tmpl := template.Must(template.New("dup").Parse("Test"))
		err := registry.Register("dup", tmpl)
		if err != nil {
			t.Fatalf("Failed to register first template: %v", err)
		}

		err = registry.Register("dup", tmpl)
		if err == nil {
			t.Error("Expected error when registering duplicate template")
		}
	})

	t.Run("Get non-existent", func(t *testing.T) {
		_, err := registry.Get("nonexistent")
		if err == nil {
			t.Error("Expected error when getting non-existent template")
		}
	})

	t.Run("Has", func(t *testing.T) {
		if !registry.Has("test") {
			t.Error("Expected Has to return true for existing template")
		}

		if registry.Has("nonexistent") {
			t.Error("Expected Has to return false for non-existent template")
		}
	})

	t.Run("List", func(t *testing.T) {
		names := registry.List()
		if len(names) == 0 {
			t.Error("Expected List to return template names")
		}

		found := false
		for _, name := range names {
			if name == "test" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected List to include 'test' template")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		registry.Clear()
		names := registry.List()
		if len(names) != 0 {
			t.Errorf("Expected empty registry after Clear, got %d templates", len(names))
		}
	})
}
