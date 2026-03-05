package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// DocumentationGenerator generates project documentation
type DocumentationGenerator struct {
	processor *templates.Processor
	validator *validators.Validator
}

// PageDoc represents documentation for a single page
type PageDoc struct {
	Name     string
	Elements []string
	Methods  []string
}

// NewDocumentationGenerator creates a new documentation generator
func NewDocumentationGenerator(processor *templates.Processor, validator *validators.Validator) *DocumentationGenerator {
	return &DocumentationGenerator{
		processor: processor,
		validator: validator,
	}
}

// Generate generates documentation for the project
func (g *DocumentationGenerator) Generate(req models.DocumentationRequest) (*models.GenerationResponse, error) {
	// 1. Scan project for pages
	pages, err := g.scanPages(req.ProjectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan pages: %w", err)
	}

	// 2. Apply template
	data := map[string]interface{}{
		"Pages": pages,
	}

	code, err := g.processor.Process("documentation", data)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	return &models.GenerationResponse{
		Code:     code,
		FileName: "PROJECT_DOCS.md",
		Metadata: map[string]string{
			"pageCount": fmt.Sprintf("%d", len(pages)),
		},
	}, nil
}

func (g *DocumentationGenerator) scanPages(rootPath string) ([]PageDoc, error) {
	var pages []PageDoc

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".ts" {
			return nil
		}

		// Check if it's likely a page object (contains "Page" in name or path)
		if !strings.Contains(path, "Page") && !strings.Contains(path, "page") {
			return nil
		}

		// Read file
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip unreadable files
		}
		content := string(contentBytes)

		// Extract info
		pageDoc := g.parsePage(content)
		if pageDoc.Name != "" {
			pages = append(pages, pageDoc)
		}

		return nil
	})

	return pages, err
}

func (g *DocumentationGenerator) parsePage(content string) PageDoc {
	var doc PageDoc

	// Extract Class Name
	classRegex := regexp.MustCompile(`export class (\w+)`)
	if match := classRegex.FindStringSubmatch(content); len(match) > 1 {
		doc.Name = match[1]
	} else {
		return doc // Not a class file
	}

	// Extract Elements (this.registerElement('name', ...))
	elementRegex := regexp.MustCompile(`this\.registerElement\(['"]([^'"]+)['"]`)
	elementMatches := elementRegex.FindAllStringSubmatch(content, -1)
	for _, match := range elementMatches {
		if len(match) > 1 {
			doc.Elements = append(doc.Elements, match[1])
		}
	}

	// Extract Methods (public async methodName() or async methodName())
	// Simple assumption: anything looking like a method definition inside the class
	methodRegex := regexp.MustCompile(`(?:public\s+)?async\s+(\w+)\(`)
	methodMatches := methodRegex.FindAllStringSubmatch(content, -1)
	for _, match := range methodMatches {
		if len(match) > 1 {
			methodName := match[1]
			// Filter out common lifecycle methods if they are async (usually not, but good safe guard)
			if methodName != "constructor" {
				doc.Methods = append(doc.Methods, methodName)
			}
		}
	}

	return doc
}
