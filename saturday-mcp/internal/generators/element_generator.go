package generators

import (
	"fmt"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// ElementGenerator generates Element/Component class code
type ElementGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// NewElementGenerator creates a new element generator
func NewElementGenerator(processor *templates.Processor, validator *validators.Validator) *ElementGenerator {
	return &ElementGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates an Element class from a request
func (g *ElementGenerator) Generate(req models.ElementGenerationRequest) (*models.GenerationResponse, error) {
	// Validate request
	if err := g.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Name":        req.Name,
		"RootSelector": req.RootSelector,
		"Methods":     req.Methods,
	}

	// Process template
	code, err := g.processor.Process("element", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Generate filename
	fileName := fmt.Sprintf("%s.ts", templates.KebabCase(req.Name))

	// Build response
	response := &models.GenerationResponse{
		Code:     code,
		FileName: fileName,
		Metadata: map[string]string{
			"type":         "element",
			"name":         req.Name,
			"rootSelector": req.RootSelector,
			"methodCount":  fmt.Sprintf("%d", len(req.Methods)),
			"description":  req.Description,
		},
	}

	return response, nil
}
