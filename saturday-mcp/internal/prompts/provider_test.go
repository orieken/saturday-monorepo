package prompts

import (
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/logging"
)

func TestProvider(t *testing.T) {
	logger := logging.NewLogger(os.Stderr)
	provider := NewProvider(logger)

	t.Run("List", func(t *testing.T) {
		prompts := provider.List()
		if len(prompts) < 2 {
			t.Errorf("Expected at least 2 prompts, got %d", len(prompts))
		}

		foundPlan := false
		for _, p := range prompts {
			if p.Name == "plan_feature" {
				foundPlan = true
				break
			}
		}

		if !foundPlan {
			t.Error("Expected plan_feature prompt to be listed")
		}
	})

	t.Run("Get_PlanFeature", func(t *testing.T) {
		args := map[string]string{
			"feature": "Login System",
		}
		messages, err := provider.Get("plan_feature", args)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}

		// Type assert content
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}

		if !strings.Contains(textContent.Text, "Login System") {
			t.Error("Expected content to contain feature name")
		}
	})

	t.Run("Get_Unknown", func(t *testing.T) {
		_, err := provider.Get("unknown_prompt", nil)
		if err == nil {
			t.Error("Expected error for unknown prompt")
		}
	})

	t.Run("Get_DebugError", func(t *testing.T) {
		args := map[string]string{
			"error":   "Element not found",
			"context": "await page.click('#submt');",
		}
		messages, err := provider.Get("debug_error", args)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}
		if !strings.Contains(textContent.Text, "Element not found") {
			t.Error("Expected content to contain error message")
		}
	})

	t.Run("Get_GenerateGherkin", func(t *testing.T) {
		args := map[string]string{
			"requirements": "User should check out",
		}
		messages, err := provider.Get("generate_gherkin", args)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}
		if !strings.Contains(textContent.Text, "User should check out") {
			t.Error("Expected content to contain requirements")
		}
	})

	t.Run("Get_VisualPageObject", func(t *testing.T) {
		args := map[string]string{
			"componentName": "CustomPage",
			"path":          "/custom",
		}
		messages, err := provider.Get("visual_page_object", args)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}
		if !strings.Contains(textContent.Text, "CustomPage") {
			t.Error("Expected content to contain component name")
		}
		if !strings.Contains(textContent.Text, "/custom") {
			t.Error("Expected content to contain path")
		}
	})

	t.Run("Get_SaturdaySME", func(t *testing.T) {
		messages, err := provider.Get("saturday_sme", nil)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}
		if !strings.Contains(textContent.Text, "Page Object Model") {
			t.Error("Expected content to mention Page Object Model")
		}
		if !strings.Contains(textContent.Text, "Fluent Flows") {
			t.Error("Expected content to mention Fluent Flows")
		}
	})

	t.Run("Get_MigrateTest", func(t *testing.T) {
		args := map[string]string{
			"legacy_code": "cy.get('.login').click()",
		}
		messages, err := provider.Get("migrate_test", args)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}
		if !strings.Contains(textContent.Text, "cy.get('.login').click()") {
			t.Error("Expected content to contain legacy code")
		}
	})

	t.Run("Get_OtelMetricsExpert", func(t *testing.T) {
		args := map[string]string{
			"flow_code": "fluentFlow.doTask()",
		}
		messages, err := provider.Get("otel_metrics_expert", args)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}
		if !strings.Contains(textContent.Text, "fluentFlow.doTask()") {
			t.Error("Expected content to contain flow code")
		}
	})

	t.Run("Get_SelfHealTest", func(t *testing.T) {
		args := map[string]string{
			"failure_log": "Element is not attached to the DOM",
		}
		messages, err := provider.Get("self_heal_test", args)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		textContent, ok := messages[0].Content.(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", messages[0].Content)
		}
		if !strings.Contains(textContent.Text, "Element is not attached to the DOM") {
			t.Error("Expected content to contain failure log")
		}
	})
}

