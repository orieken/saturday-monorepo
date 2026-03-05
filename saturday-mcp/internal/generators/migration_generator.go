package generators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// MigrationGenerator helps migrate legacy code to framework patterns
type MigrationGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// NewMigrationGenerator creates a new migration generator
func NewMigrationGenerator(processor *templates.Processor, validator *validators.Validator) *MigrationGenerator {
	return &MigrationGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates framework code from legacy source
func (g *MigrationGenerator) Generate(req models.MigrationRequest) (*models.GenerationResponse, error) {
	// Validate request
	if err := g.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if req.Type == "page" {
		return g.generatePageFromSource(req.SourceCode)
	}

	return nil, fmt.Errorf("unsupported migration type: %s", req.Type)
}

func (g *MigrationGenerator) generatePageFromSource(source string) (*models.GenerationResponse, error) {
	// Simple regex extraction for MVP
	// Matches: page.action('selector') or page.action("selector")
	// Groups: 1=action, 2=selector (single quote), 3=selector (double quote)
	actionRegex := regexp.MustCompile(`page\.(click|fill|type|check|uncheck)\((?:'([^']*)'|"([^"]*)")`)
	
	matches := actionRegex.FindAllStringSubmatch(source, -1)
	
	elements := make([]models.ElementDefinition, 0)
	seenSelectors := make(map[string]bool)
	
	for i, match := range matches {
		if len(match) < 3 {
			continue
		}
		
		// action := match[1]
		selector := match[2] // single quote match
		if selector == "" {
			selector = match[3] // double quote match
		}

		// Escape double quotes for the generated code
		selector = strings.ReplaceAll(selector, "\"", "\\\"")

		if seenSelectors[selector] {
			continue
		}
		seenSelectors[selector] = true

		// Simple heuristics for element naming and type
		name := fmt.Sprintf("element%d", i+1)
		elementType := "unknown"

		if strings.Contains(selector, "button") || strings.Contains(selector, "btn") {
			name = fmt.Sprintf("button%d", i+1)
			elementType = "button"
		} else if strings.Contains(selector, "input") {
			name = fmt.Sprintf("input%d", i+1)
			elementType = "input"
		} else if strings.Contains(selector, "link") || strings.Contains(selector, "a[") {
			name = fmt.Sprintf("link%d", i+1)
			elementType = "link"
		}

		elements = append(elements, models.ElementDefinition{
			Name:     name,
			Selector: selector,
			Type:     elementType,
		})
	}

	// Use the "Page" template to generate a draft page
	data := map[string]interface{}{
		"Name":     "DraftMigrated",
		"Path":     "/migrated/path",
		"Elements": elements,
	}

	code, err := g.processor.Process("page", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	return &models.GenerationResponse{
		Code:     code,
		FileName: "draft-migrated-page.ts",
		Metadata: map[string]string{
			"type":         "page",
			"elementCount": fmt.Sprintf("%d", len(elements)),
			"migrated":     "true",
		},
	}, nil
}
