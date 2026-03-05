package generators

import (
	"fmt"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// FlowGenerator generates Flow class code
type FlowGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// NewFlowGenerator creates a new flow generator
func NewFlowGenerator(processor *templates.Processor, validator *validators.Validator) *FlowGenerator {
	return &FlowGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates a Flow class from a request
func (g *FlowGenerator) Generate(req models.FlowGenerationRequest) (*models.GenerationResponse, error) {
	// Validate request
	if err := g.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Name":  req.Name,
		"Steps": req.Steps,
	}

	// Process template
	code, err := g.processor.Process("flow", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Generate filename
	fileName := fmt.Sprintf("%s-flow.ts", templates.KebabCase(req.Name))

	// Build response
	response := &models.GenerationResponse{
		Code:     code,
		FileName: fileName,
		Metadata: map[string]string{
			"type":        "flow",
			"name":        req.Name,
			"stepCount":   fmt.Sprintf("%d", len(req.Steps)),
			"description": req.Description,
		},
	}

	return response, nil
}

// ValidateRequest validates a flow generation request without generating code
func (g *FlowGenerator) ValidateRequest(req models.FlowGenerationRequest) error {
	return g.validator.Validate(req)
}
