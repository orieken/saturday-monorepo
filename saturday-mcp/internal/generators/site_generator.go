package generators

import (
	"fmt"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// SiteGenerator generates Site class code
type SiteGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// NewSiteGenerator creates a new site generator
func NewSiteGenerator(processor *templates.Processor, validator *validators.Validator) *SiteGenerator {
	return &SiteGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates a Site class from a request
func (g *SiteGenerator) Generate(req models.SiteGenerationRequest) (*models.GenerationResponse, error) {
	// Validate request
	if err := g.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Name":    req.Name,
		"BaseURL": req.BaseURL,
		"Pages":   req.Pages,
		"Imports": req.Pages, // Same as pages for imports
		"Flows":   req.Flows,
	}

	// Process template
	code, err := g.processor.Process("site", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Generate filename
	fileName := fmt.Sprintf("%s-site.ts", templates.KebabCase(req.Name))

	// Build response
	response := &models.GenerationResponse{
		Code:     code,
		FileName: fileName,
		Metadata: map[string]string{
			"type":        "site",
			"name":        req.Name,
			"baseUrl":     req.BaseURL,
			"pageCount":   fmt.Sprintf("%d", len(req.Pages)),
			"description": req.Description,
		},
	}

	return response, nil
}

// ValidateRequest validates a site generation request without generating code
func (g *SiteGenerator) ValidateRequest(req models.SiteGenerationRequest) error {
	return g.validator.Validate(req)
}
