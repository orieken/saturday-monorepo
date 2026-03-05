package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "saturday",
	Short: "Saturday Framework CLI - Code generation and analysis tools",
	Long: `Saturday Framework CLI provides command-line access to code generation,
analysis, and validation tools for Playwright test automation using the Saturday framework.

This CLI wraps the Saturday MCP server functionality, allowing you to:
- Generate Page Objects, Flows, Steps, and Services
- Analyze existing framework code
- Validate patterns and suggest improvements
- Migrate legacy code to Saturday patterns
- Generate project documentation`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output directory (default: current directory)")
}
