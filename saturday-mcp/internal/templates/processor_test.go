package templates

import (
	"strings"
	"testing"
	"text/template"
)

func TestProcessor(t *testing.T) {
	registry := NewRegistry()
	cache := NewCache(0) // No TTL for tests
	processor := NewProcessor(registry, cache)

	// Register a test template
	tmpl := template.Must(template.New("greeting").Parse("Hello {{.Name}}!"))
	registry.Register("greeting", tmpl)

	t.Run("Process template", func(t *testing.T) {
		data := map[string]string{"Name": "World"}
		result, err := processor.Process("greeting", data)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		expected := "Hello World!"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Process non-existent template", func(t *testing.T) {
		_, err := processor.Process("nonexistent", nil)
		if err == nil {
			t.Error("Expected error when processing non-existent template")
		}
	})

	t.Run("Cache hit", func(t *testing.T) {
		data := map[string]string{"Name": "Cache"}

		// First call - cache miss
		result1, err := processor.Process("greeting", data)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		// Second call - should hit cache
		result2, err := processor.Process("greeting", data)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		if result1 != result2 {
			t.Error("Expected same result from cache")
		}

		if cache.Size() == 0 {
			t.Error("Expected cache to have entries")
		}
	})

	t.Run("ProcessWithFuncs", func(t *testing.T) {
		funcTmpl := template.Must(template.New("func").Parse("{{.}}"))
		registry.Register("func", funcTmpl)

		customFuncs := template.FuncMap{
			"custom": func(s string) string {
				return strings.ToUpper(s)
			},
		}

		result, err := processor.ProcessWithFuncs("func", "test", customFuncs)
		if err != nil {
			t.Fatalf("Failed to process template with funcs: %v", err)
		}

		if result != "test" {
			t.Errorf("Expected test, got %s", result)
		}
	})

	t.Run("Validate", func(t *testing.T) {
		data := map[string]string{"Name": "Valid"}
		err := processor.Validate("greeting", data)
		if err != nil {
			t.Errorf("Expected validation to pass: %v", err)
		}

		err = processor.Validate("nonexistent", data)
		if err == nil {
			t.Error("Expected validation to fail for non-existent template")
		}
	})
}
