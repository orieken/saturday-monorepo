package generators

import (
	"fmt"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// StepGenerator generates Cucumber step definition code
type StepGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// NewStepGenerator creates a new step definition generator
func NewStepGenerator(processor *templates.Processor, validator *validators.Validator) *StepGenerator {
	return &StepGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates step definitions from a request
func (g *StepGenerator) Generate(req models.StepDefinitionRequest) (*models.GenerationResponse, error) {
	// Validate request
	if err := g.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Default language to TypeScript
	language := req.Language
	if language == "" {
		language = "typescript"
	}

	// Prepare template data
	data := map[string]interface{}{
		"Steps":    req.Steps,
		"Language": language,
	}

	// Process template
	code, err := g.processor.Process("steps", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Generate filename based on language
	var fileName string
	if language == "typescript" {
		fileName = "steps.ts"
	} else {
		fileName = "steps.js"
	}

	// Build response
	response := &models.GenerationResponse{
		Code:     code,
		FileName: fileName,
		Metadata: map[string]string{
			"type":      "steps",
			"language":  language,
			"stepCount": fmt.Sprintf("%d", len(req.Steps)),
		},
	}

	return response, nil
}

// ValidateRequest validates a step definition request without generating code
func (g *StepGenerator) ValidateRequest(req models.StepDefinitionRequest) error {
	return g.validator.Validate(req)
}
