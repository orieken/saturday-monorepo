package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate legacy code to Saturday patterns",
	Long:  `Analyze and migrate legacy Playwright code to Saturday Framework patterns.`,
}

var migratePageCmd = &cobra.Command{
	Use:   "page [file]",
	Short: "Migrate legacy code to Page Object",
	Long:  `Parse legacy Playwright code and generate a draft Page Object.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMigratePage,
}

var docsCmd = &cobra.Command{
	Use:   "docs [project-path] [output-file]",
	Short: "Generate project documentation",
	Long:  `Scan a project and generate markdown documentation for all Page Objects.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runGenerateDocs,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(docsCmd)
	
	migrateCmd.AddCommand(migratePageCmd)

	// Flags
	migratePageCmd.Flags().Bool("write", false, "Write to file")
}

func runMigratePage(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	write, _ := cmd.Flags().GetBool("write")

	// Read source file
	sourceBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	req := models.MigrationRequest{
		SourceCode: string(sourceBytes),
		Type:       "page",
	}

	gen, err := setupGenerator()
	if err != nil {
		return err
	}

	resp, err := gen.MigrationGenerator.Generate(req)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if write {
		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = "."
		}
		outPath := filepath.Join(outputDir, "lib", "pages", resp.FileName)
		if err := writeFile(outPath, resp.Code); err != nil {
			return err
		}
		fmt.Printf("✓ Migrated to: %s\n", outPath)
		fmt.Printf("⚠️  Please review and adjust the generated code as needed.\n")
	} else {
		fmt.Println(resp.Code)
	}

	return nil
}

func runGenerateDocs(cmd *cobra.Command, args []string) error {
	projectPath := args[0]
	outputPath := args[1]

	req := models.DocumentationRequest{
		ProjectPath: projectPath,
		OutputPath:  outputPath,
	}

	gen, err := setupGenerator()
	if err != nil {
		return err
	}

	resp, err := gen.DocumentationGenerator.Generate(req)
	if err != nil {
		return fmt.Errorf("documentation generation failed: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(resp.Code), 0644); err != nil {
		return fmt.Errorf("failed to write documentation: %w", err)
	}

	fmt.Printf("✓ Documentation generated: %s\n", outputPath)
	if pageCount, ok := resp.Metadata["pageCount"]; ok {
		fmt.Printf("  Pages documented: %s\n", pageCount)
	}

	return nil
}
