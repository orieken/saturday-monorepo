package resources

import (
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/templates"
)

func TestProvider(t *testing.T) {
	// Setup dependencies
	logger := logging.NewLogger(os.Stderr)
	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	
	// Create provider
	provider := NewProvider(logger, loader)

	t.Run("List", func(t *testing.T) {
		resources := provider.List()
		if len(resources) == 0 {
			t.Fatal("Expected resources to be listed")
		}

		foundPage := false
		for _, r := range resources {
			if r.Name == "Template: page" {
				foundPage = true
				if r.URI != "saturday://templates/page" {
					t.Errorf("Expected URI saturday://templates/page, got %s", r.URI)
				}
			}
		}

		if !foundPage {
			t.Error("Expected page template to be listed")
		}
	})

	t.Run("Read_Valid", func(t *testing.T) {
		contents, err := provider.Read("saturday://templates/page")
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}

		if len(contents) != 1 {
			t.Fatalf("Expected 1 content item, got %d", len(contents))
		}

		// Type assert to TextResourceContents
		textContent, ok := contents[0].(mcp.TextResourceContents)
		if !ok {
			t.Fatalf("Expected TextResourceContents, got %T", contents[0])
		}

		if textContent.MIMEType != "text/plain" {
			t.Errorf("Expected text/plain, got %s", textContent.MIMEType)
		}

		if len(textContent.Text) == 0 {
			t.Error("Expected content text not to be empty")
		}
	})

	t.Run("Read_InvalidScheme", func(t *testing.T) {
		_, err := provider.Read("http://templates/page")
		if err == nil {
			t.Error("Expected error for invalid scheme")
		}
	})

	t.Run("Read_InvalidPath", func(t *testing.T) {
		_, err := provider.Read("saturday://templates")
		if err == nil {
			t.Error("Expected error for invalid path")
		}
	})

	t.Run("Read_UnknownTemplate", func(t *testing.T) {
		_, err := provider.Read("saturday://templates/unknown")
		if err == nil {
			t.Error("Expected error for unknown template")
		}
	})
}
