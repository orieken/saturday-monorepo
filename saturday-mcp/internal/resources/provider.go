package resources

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/templates"
)

const (
	// ResourceScheme is the URI scheme for Saturday resources
	ResourceScheme = "saturday"
	// TemplatesPath is the path prefix for template resources
	TemplatesPath = "templates"
)

// Provider manages MCP resources
type Provider struct {
	logger *logging.Logger
	loader *templates.Loader
}

// NewProvider creates a new resource provider
func NewProvider(logger *logging.Logger, loader *templates.Loader) *Provider {
	return &Provider{
		logger: logger,
		loader: loader,
	}
}

// List returns all available resources
func (p *Provider) List() []mcp.Resource {
	var resources []mcp.Resource

	// List templates
	tmplNames := []string{"site", "page", "flow", "steps"}
	for _, name := range tmplNames {
		uri := fmt.Sprintf("%s://%s/%s", ResourceScheme, TemplatesPath, name)
		resources = append(resources, mcp.Resource{
			URI:         uri,
			Name:        fmt.Sprintf("Template: %s", name),
			Description: fmt.Sprintf("Template for generating %s", name),
			MIMEType:    "text/plain",
		})
	}

	return resources
}

// Read returns the content of a specific resource
func (p *Provider) Read(uri string) ([]mcp.ResourceContents, error) {
	// Parse URI
	// Format: saturday://templates/{name}
	if !strings.HasPrefix(uri, ResourceScheme+"://") {
		return nil, fmt.Errorf("invalid resource scheme: %s", uri)
	}

	path := strings.TrimPrefix(uri, ResourceScheme+"://")
	parts := strings.Split(path, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource path format: %s", path)
	}

	category := parts[0]
	name := parts[1]

	if category == TemplatesPath {
		content, err := p.loader.GetRawContent(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load template %s: %w", name, err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "text/plain",
				Text:     string(content),
			},
		}, nil
	}

	return nil, fmt.Errorf("resource category not found: %s", category)
}
