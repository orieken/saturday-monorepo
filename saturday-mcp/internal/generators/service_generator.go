package generators

import (
	"fmt"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// ServiceGenerator generates API Service class code
type ServiceGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// NewServiceGenerator creates a new service generator
func NewServiceGenerator(processor *templates.Processor, validator *validators.Validator) *ServiceGenerator {
	return &ServiceGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates a Service class from a request
func (g *ServiceGenerator) Generate(req models.ServiceGenerationRequest) (*models.GenerationResponse, error) {
	// Validate request
	if err := g.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Name":      req.Name,
		"BaseURL":   req.BaseURL,
		"Endpoints": req.Endpoints,
	}

	// Process template
	code, err := g.processor.Process("service", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Generate filename
	fileName := fmt.Sprintf("%s-service.ts", templates.KebabCase(req.Name))

	// Build response
	response := &models.GenerationResponse{
		Code:     code,
		FileName: fileName,
		Metadata: map[string]string{
			"type":          "service",
			"name":          req.Name,
			"path":          req.BaseURL,
			"endpointCount": fmt.Sprintf("%d", len(req.Endpoints)),
			"description":   req.Description,
		},
	}

	return response, nil
}
