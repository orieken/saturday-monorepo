package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze code and projects",
	Long:  `Analyze existing code for framework usage, patterns, and performance issues.`,
}

var analyzeFrameworkCmd = &cobra.Command{
	Use:   "framework [path]",
	Short: "Analyze framework structure",
	Long:  `Scan a project directory and analyze its Saturday framework structure.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAnalyzeFramework,
}

var analyzePerformanceCmd = &cobra.Command{
	Use:   "performance [path]",
	Short: "Analyze performance issues",
	Long:  `Scan code for performance bottlenecks and anti-patterns.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAnalyzePerformance,
}

var suggestCmd = &cobra.Command{
	Use:   "suggest [path]",
	Short: "Suggest code improvements",
	Long:  `Analyze code and suggest improvements based on best practices.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSuggest,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(suggestCmd)
	
	analyzeCmd.AddCommand(analyzeFrameworkCmd)
	analyzeCmd.AddCommand(analyzePerformanceCmd)

	// Flags
	analyzeFrameworkCmd.Flags().Bool("json", false, "Output as JSON")
	analyzePerformanceCmd.Flags().Bool("json", false, "Output as JSON")
	suggestCmd.Flags().Bool("json", false, "Output as JSON")
}

func runAnalyzeFramework(cmd *cobra.Command, args []string) error {
	path := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	logger := logging.NewLogger(os.Stderr)
	analyzer := analyzers.NewFrameworkAnalyzer(logger)

	result, err := analyzer.Analyze(path)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if jsonOutput {
		return printJSON(result)
	}

	// Pretty print
	fmt.Printf("📊 Framework Analysis Results\n\n")
	fmt.Printf("Pages:  %d\n", result.Stats.PageCount)
	fmt.Printf("Flows:  %d\n", result.Stats.FlowCount)
	fmt.Printf("Steps:  %d\n", result.Stats.StepDefCount)
	fmt.Printf("\n")

	if len(result.Inventory.Pages) > 0 {
		fmt.Printf("📄 Pages:\n")
		for _, page := range result.Inventory.Pages {
			fmt.Printf("  - %s (%s)\n", page.Name, page.FilePath)
		}
		fmt.Printf("\n")
	}

	if len(result.Inventory.Flows) > 0 {
		fmt.Printf("🔄 Flows:\n")
		for _, flow := range result.Inventory.Flows {
			fmt.Printf("  - %s (%s)\n", flow.Name, flow.FilePath)
		}
		fmt.Printf("\n")
	}

	return nil
}

func runAnalyzePerformance(cmd *cobra.Command, args []string) error {
	path := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	logger := logging.NewLogger(os.Stderr)
	analyzer := analyzers.NewPerformanceAnalyzer(logger)

	result, err := analyzer.Analyze(path)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if jsonOutput {
		return printJSON(result)
	}

	// Pretty print
	fmt.Printf("⚡ Performance Analysis Results\n\n")
	
	if len(result.Suggestions) == 0 {
		fmt.Printf("✓ No performance issues found!\n")
		return nil
	}

	fmt.Printf("Found %d issue(s):\n\n", len(result.Suggestions))
	
	for i, issue := range result.Suggestions {
		fmt.Printf("%d. [%s] %s\n", i+1, issue.Severity, issue.Message)
		if issue.File != "" {
			fmt.Printf("   File: %s", issue.File)
			if issue.Line > 0 {
				fmt.Printf(":%d", issue.Line)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}

	return nil
}

func runSuggest(cmd *cobra.Command, args []string) error {
	path := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	logger := logging.NewLogger(os.Stderr)
	analyzer := analyzers.NewImprovementAnalyzer(logger)

	result, err := analyzer.Analyze(path)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if jsonOutput {
		return printJSON(result)
	}

	// Pretty print
	fmt.Printf("💡 Code Improvement Suggestions\n\n")
	
	if len(result.Suggestions) == 0 {
		fmt.Printf("✓ No improvements suggested - code looks good!\n")
		return nil
	}

	fmt.Printf("Found %d suggestion(s):\n\n", len(result.Suggestions))
	
	for i, issue := range result.Suggestions {
		fmt.Printf("%d. [%s] %s\n", i+1, issue.Severity, issue.Message)
		if issue.File != "" {
			fmt.Printf("   File: %s", issue.File)
			if issue.Line > 0 {
				fmt.Printf(":%d", issue.Line)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}

	return nil
}

func printJSONAnalysis(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
