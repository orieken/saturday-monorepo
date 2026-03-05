package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"time"

	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code artifacts",
	Long:  `Generate Page Objects, Flows, Steps, Services, and other code artifacts.`,
}

var generatePageCmd = &cobra.Command{
	Use:   "page [name]",
	Short: "Generate a Page Object",
	Long:  `Generate a Page Object class with specified elements.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGeneratePage,
}

var generateFlowCmd = &cobra.Command{
	Use:   "flow [name]",
	Short: "Generate a Flow class",
	Long:  `Generate a Flow class with specified steps.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGenerateFlow,
}

var generateServiceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Generate a Service class",
	Long:  `Generate an API Service class with specified endpoints.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGenerateService,
}

var generateSiteCmd = &cobra.Command{
	Use:   "site [name]",
	Short: "Generate a Site class",
	Long:  `Generate a Site class with pages and flows.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGenerateSite,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generatePageCmd)
	generateCmd.AddCommand(generateFlowCmd)
	generateCmd.AddCommand(generateServiceCmd)
	generateCmd.AddCommand(generateSiteCmd)

	// Page flags
	generatePageCmd.Flags().String("path", "/", "Page path/URL")
	generatePageCmd.Flags().StringSlice("elements", []string{}, "Elements (format: name:selector:type)")
	generatePageCmd.Flags().Bool("write", false, "Write to file")

	// Flow flags
	generateFlowCmd.Flags().StringSlice("steps", []string{}, "Flow steps")
	generateFlowCmd.Flags().Bool("write", false, "Write to file")

	// Service flags
	generateServiceCmd.Flags().String("base-url", "", "Base URL for the service")
	generateServiceCmd.Flags().StringSlice("endpoints", []string{}, "Endpoints (format: name:method:path)")
	generateServiceCmd.Flags().Bool("write", false, "Write to file")

	// Site flags
	generateSiteCmd.Flags().String("url", "", "Base URL for the site")
	generateSiteCmd.Flags().StringSlice("pages", []string{}, "Pages")
	generateSiteCmd.Flags().StringSlice("flows", []string{}, "Flows")
	generateSiteCmd.Flags().Bool("write", false, "Write to file")
	generateSiteCmd.MarkFlagRequired("url")
}

func setupGenerator() (*generators.Generator, error) {
	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	if err := loader.LoadAll(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	validator := validators.NewValidator()

	siteGen := generators.NewSiteGenerator(processor, validator)
	pageGen := generators.NewPageGenerator(processor, validator)
	flowGen := generators.NewFlowGenerator(processor, validator)
	stepGen := generators.NewStepGenerator(processor, validator)
	elementGen := generators.NewElementGenerator(processor, validator)
	serviceGen := generators.NewServiceGenerator(processor, validator)
	migrationGen := generators.NewMigrationGenerator(processor, validator)
	docGen := generators.NewDocumentationGenerator(processor, validator)

	return generators.NewGenerator(siteGen, pageGen, flowGen, stepGen, elementGen, serviceGen, migrationGen, docGen), nil
}

func runGeneratePage(cmd *cobra.Command, args []string) error {
	name := args[0]
	path, _ := cmd.Flags().GetString("path")
	elementStrs, _ := cmd.Flags().GetStringSlice("elements")
	write, _ := cmd.Flags().GetBool("write")

	// Parse elements
	var elements []models.ElementDefinition
	for _, elemStr := range elementStrs {
		// Simple parsing: name:selector:type
		parts := parseElement(elemStr)
		if len(parts) == 3 {
			elements = append(elements, models.ElementDefinition{
				Name:     parts[0],
				Selector: parts[1],
				Type:     parts[2],
			})
		}
	}

	req := models.PageGenerationRequest{
		Name:     name,
		Path:     path,
		Elements: elements,
	}

	gen, err := setupGenerator()
	if err != nil {
		return err
	}

	resp, err := gen.PageGenerator.Generate(req)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	if write {
		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = "."
		}
		filePath := filepath.Join(outputDir, "lib", "pages", resp.FileName)
		if err := writeFile(filePath, resp.Code); err != nil {
			return err
		}
		fmt.Printf("✓ Generated: %s\n", filePath)
	} else {
		fmt.Println(resp.Code)
	}

	return nil
}

func runGenerateFlow(cmd *cobra.Command, args []string) error {
	name := args[0]
	steps, _ := cmd.Flags().GetStringSlice("steps")
	write, _ := cmd.Flags().GetBool("write")

	req := models.FlowGenerationRequest{
		Name:  name,
		Steps: steps,
	}

	gen, err := setupGenerator()
	if err != nil {
		return err
	}

	resp, err := gen.FlowGenerator.Generate(req)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	if write {
		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = "."
		}
		filePath := filepath.Join(outputDir, "lib", "flows", resp.FileName)
		if err := writeFile(filePath, resp.Code); err != nil {
			return err
		}
		fmt.Printf("✓ Generated: %s\n", filePath)
	} else {
		fmt.Println(resp.Code)
	}

	return nil
}

func runGenerateService(cmd *cobra.Command, args []string) error {
	name := args[0]
	baseURL, _ := cmd.Flags().GetString("base-url")
	endpointStrs, _ := cmd.Flags().GetStringSlice("endpoints")
	write, _ := cmd.Flags().GetBool("write")

	// Parse endpoints
	var endpoints []models.EndpointDefinition
	for _, epStr := range endpointStrs {
		parts := parseEndpoint(epStr)
		if len(parts) == 3 {
			endpoints = append(endpoints, models.EndpointDefinition{
				Name:   parts[0],
				Method: parts[1],
				Path:   parts[2],
			})
		}
	}

	req := models.ServiceGenerationRequest{
		Name:      name,
		BaseURL:   baseURL,
		Endpoints: endpoints,
	}

	gen, err := setupGenerator()
	if err != nil {
		return err
	}

	resp, err := gen.ServiceGenerator.Generate(req)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	if write {
		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = "."
		}
		filePath := filepath.Join(outputDir, "lib", "services", resp.FileName)
		if err := writeFile(filePath, resp.Code); err != nil {
			return err
		}
		fmt.Printf("✓ Generated: %s\n", filePath)
	} else {
		fmt.Println(resp.Code)
	}

	return nil
}

func runGenerateSite(cmd *cobra.Command, args []string) error {
	name := args[0]
	url, _ := cmd.Flags().GetString("url")
	pages, _ := cmd.Flags().GetStringSlice("pages")
	flows, _ := cmd.Flags().GetStringSlice("flows")
	write, _ := cmd.Flags().GetBool("write")

	req := models.SiteGenerationRequest{
		Name:    name,
		BaseURL: url,
		Pages:   pages,
		Flows:   flows,
	}

	gen, err := setupGenerator()
	if err != nil {
		return err
	}

	resp, err := gen.SiteGenerator.Generate(req)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	if write {
		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = "."
		}
		filePath := filepath.Join(outputDir, "lib", "sites", resp.FileName)
		if err := writeFile(filePath, resp.Code); err != nil {
			return err
		}
		fmt.Printf("✓ Generated: %s\n", filePath)
	} else {
		fmt.Println(resp.Code)
	}

	return nil
}

func parseElement(s string) []string {
	// Simple split by colon
	result := []string{}
	current := ""
	for _, ch := range s {
		if ch == ':' {
			result = append(result, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func parseEndpoint(s string) []string {
	return parseElement(s) // Same logic
}

func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
