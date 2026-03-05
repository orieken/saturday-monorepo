package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate code against Saturday patterns",
	Long:  `Validate code against Saturday framework patterns and naming conventions.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Flags
	validateCmd.Flags().Bool("json", false, "Output as JSON")
	validateCmd.Flags().Bool("strict", false, "Fail on warnings")
}

func runValidate(cmd *cobra.Command, args []string) error {
	path := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")
	strict, _ := cmd.Flags().GetBool("strict")

	logger := logging.NewLogger(os.Stderr)
	validator := analyzers.NewPatternValidator(logger)

	result, err := validator.Validate(path)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Pretty print
	fmt.Printf("🔍 Pattern Validation Results\n\n")
	
	if result.Valid {
		fmt.Printf("✓ All patterns valid!\n")
		return nil
	}

	fmt.Printf("Found %d issue(s):\n\n", len(result.Issues))
	
	errorCount := 0
	warningCount := 0
	
	for i, issue := range result.Issues {
		icon := "⚠️"
		if issue.Severity == "error" {
			icon = "❌"
			errorCount++
		} else {
			warningCount++
		}
		
		fmt.Printf("%d. %s [%s] %s\n", i+1, icon, issue.Severity, issue.Message)
		if issue.File != "" {
			fmt.Printf("   File: %s", issue.File)
			if issue.Line > 0 {
				fmt.Printf(":%d", issue.Line)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)

	if strict && warningCount > 0 {
		return fmt.Errorf("validation failed in strict mode")
	}

	if errorCount > 0 {
		return fmt.Errorf("validation failed with errors")
	}

	return nil
}
