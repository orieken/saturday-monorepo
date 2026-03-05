package generators

import (
	"fmt"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// PageGenerator generates Page class code
type PageGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// NewPageGenerator creates a new page generator
func NewPageGenerator(processor *templates.Processor, validator *validators.Validator) *PageGenerator {
	return &PageGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates a Page class from a request
func (g *PageGenerator) Generate(req models.PageGenerationRequest) (*models.GenerationResponse, error) {
	// Validate request
	if err := g.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Name":     req.Name,
		"Path":     req.Path,
		"Elements": req.Elements,
	}

	// Process template
	code, err := g.processor.Process("page", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Generate filename
	fileName := fmt.Sprintf("%s-page.ts", templates.KebabCase(req.Name))

	// Build response
	response := &models.GenerationResponse{
		Code:     code,
		FileName: fileName,
		Metadata: map[string]string{
			"type":         "page",
			"name":         req.Name,
			"path":         req.Path,
			"elementCount": fmt.Sprintf("%d", len(req.Elements)),
			"description":  req.Description,
		},
	}

	return response, nil
}

// ValidateRequest validates a page generation request without generating code
func (g *PageGenerator) ValidateRequest(req models.PageGenerationRequest) error {
	return g.validator.Validate(req)
}
