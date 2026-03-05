package prompts

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// PromptHandler handles prompt execution
type PromptHandler func(arguments map[string]string) ([]mcp.PromptMessage, error)

// PromptDefinition defines a single prompt
type PromptDefinition struct {
	Name        string
	Description string
	Arguments   []mcp.PromptArgument
	Handler     PromptHandler
}

// Provider manages MCP prompts
type Provider struct {
	logger  *logging.Logger
	prompts map[string]PromptDefinition
}

// NewProvider creates a new prompt provider
func NewProvider(logger *logging.Logger) *Provider {
	p := &Provider{
		logger:  logger,
		prompts: make(map[string]PromptDefinition),
	}
	p.registerDefaults()
	return p
}

// registerDefaults registers the default set of prompts
func (p *Provider) registerDefaults() {
	p.prompts["plan_feature"] = PromptDefinition{
		Name:        "plan_feature",
		Description: "Plan a new feature implementation",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "feature",
				Description: "Name or description of the feature",
				Required:    true,
			},
		},
		Handler: p.handlePlanFeature,
	}

	p.prompts["explain_framework"] = PromptDefinition{
		Name:        "explain_framework",
		Description: "Explain how the Saturday framework works",
		Arguments:   []mcp.PromptArgument{},
		Handler:     p.handleExplainFramework,
	}

	p.prompts["debug_error"] = PromptDefinition{
		Name:        "debug_error",
		Description: "Debug a test failure or error message",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "error",
				Description: "The error message or failure log",
				Required:    true,
			},
			{
				Name:        "context",
				Description: "Optional code context or step definition",
				Required:    false,
			},
		},
		Handler: p.handleDebugError,
	}

	p.prompts["generate_gherkin"] = PromptDefinition{
		Name:        "generate_gherkin",
		Description: "Generate Gherkin scenarios from requirements",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "requirements",
				Description: "User story or feature requirements",
				Required:    true,
			},
		},
		Handler: p.handleGenerateGherkin,
	}

	p.prompts["visual_page_object"] = PromptDefinition{
		Name:        "visual_page_object",
		Description: "Generate a Page Object from a UI screenshot (requires attaching an image)",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "componentName",
				Description: "Name of the page/component (e.g. LoginPage)",
				Required:    false,
			},
			{
				Name:        "path",
				Description: "URL path (e.g. /login)",
				Required:    false,
			},
		},
		Handler: p.handleVisualPageObject,
	}

	p.prompts["implement_feature"] = PromptDefinition{
		Name:        "implement_feature",
		Description: "Orchestrate the 'Autonomous QA Engineer' workflow to implement and verify a feature",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "feature",
				Description: "The name of the feature to implement (e.g., 'User Login')",
				Required:    true,
			},
		},
		Handler: p.handleImplementFeature,
	}

	p.prompts["saturday_sme"] = PromptDefinition{
		Name:        "saturday_sme",
		Description: "Act as a Subject Matter Expert for the Saturday testing framework",
		Arguments:   []mcp.PromptArgument{},
		Handler:     p.handleSaturdaySME,
	}

	p.prompts["migrate_test"] = PromptDefinition{
		Name:        "migrate_test",
		Description: "Migrate legacy tests to the Saturday framework and utilize Playwright MCP",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "legacy_code",
				Description: "The legacy test code to migrate",
				Required:    true,
			},
		},
		Handler: p.handleMigrateTest,
	}

	p.prompts["otel_metrics_expert"] = PromptDefinition{
		Name:        "otel_metrics_expert",
		Description: "Analyze flows to add OpenTelemetry tracing and custom metrics",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "flow_code",
				Description: "The Saturday flow code to analyze",
				Required:    true,
			},
		},
		Handler: p.handleOtelMetricsExpert,
	}

	p.prompts["self_heal_test"] = PromptDefinition{
		Name:        "self_heal_test",
		Description: "Self-heal a failing test using Saturday patterns and Playwright MCP",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "failure_log",
				Description: "The log of the failing test",
				Required:    true,
			},
		},
		Handler: p.handleSelfHealTest,
	}
}


func (p *Provider) handleVisualPageObject(args map[string]string) ([]mcp.PromptMessage, error) {
	name := args["componentName"]
	if name == "" {
		name = "MyPage"
	}
	path := args["path"]
	if path == "" {
		path = "/"
	}

	content := fmt.Sprintf(`You are an expert Saturday Framework developer with visual intelligence.
The user has attached a UI screenshot. **Please analyze this image and generate a complete Saturday Page Object.**

**Configuration:**
- Class Name: `+"`%s`"+`
- Path: `+"`%s`"+`

**Instructions:**
1.  **Analyze the UI**: Identify all interactive elements (buttons, inputs, links, dropdowns).
2.  **Naming Conventions**: Use camelCase with type suffixes (e.g., `+"`emailInput`"+`, `+"`submitButton`"+`).
3.  **Selectors**: Suggest robust CSS selectors (avoid generic classes).
4.  **Structure**: Generate the TypeScript code adhering to the Saturday pattern:
    - Extend `+"`Page`"+` from `+"`@/core`"+`.
    - Register elements in the constructor.

**Example Output:**
`+"```typescript"+`
import { Page } from 'src/lib/pages/Page';

export class %s extends Page {
    constructor() {
        super('%s');
    }

    // Register elements here based on the image
}
`+"```"+`
`, name, path, name, path)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

// List returns all available prompts
func (p *Provider) List() []mcp.Prompt {
	var results []mcp.Prompt
	for _, def := range p.prompts {
		results = append(results, mcp.Prompt{
			Name:        def.Name,
			Description: def.Description,
			Arguments:   def.Arguments,
		})
	}
	return results
}

// Get executes a specific prompt
func (p *Provider) Get(name string, arguments map[string]string) ([]mcp.PromptMessage, error) {
	def, ok := p.prompts[name]
	if !ok {
		return nil, fmt.Errorf("prompt not found: %s", name)
	}
	return def.Handler(arguments)
}

// Handler implementations

func (p *Provider) handlePlanFeature(args map[string]string) ([]mcp.PromptMessage, error) {
	feature := args["feature"]
	content := fmt.Sprintf(`You are an expert Saturday Framework developer.
The user wants to implement the following feature: "%s".

Please create a step-by-step implementation plan that includes:
1.  **Page Objects**: What new pages are needed?
2.  **Flows**: What user journeys (Flows) should be created?
3.  **Step Definitions**: What Gherkin steps are required?
4.  **Tests**: Examples of usage.

Use the Saturday Framework patterns (Page Object Model, Fluent Flows).
`, feature)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleExplainFramework(args map[string]string) ([]mcp.PromptMessage, error) {
	content := `Explain the Saturday Framework architecture, focusing on:
1.  **Site**: The root aggregator.
2.  **Pages**: Encapsulation of UI interactables.
3.  **Flows**: Reusable sequences of interactions.
4.  **Step Definitions**: Clean connection to Cucumber.

Provide examples of how they fit together.`

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleDebugError(args map[string]string) ([]mcp.PromptMessage, error) {
	errMsg := args["error"]
	context := args["context"]

	content := fmt.Sprintf(`You are an expert Saturday Framework debugger.
The user has encountered the following error during a test run:

Error:
%s

Context:
%s

Please analyze this error and provide:
1.  **Root Cause Analysis**: What likely went wrong?
2.  **Debugging Steps**: Specific actions to verify the issue (e.g., check selector, check network tab).
3.  **Fix Suggestions**: Code changes or configuration updates to resolve the issue.

Consider common Saturday pitfalls like:
- Missing await on async calls.
- Incorrect page object registration.
- Selector specificity issues.
`, errMsg, context)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleGenerateGherkin(args map[string]string) ([]mcp.PromptMessage, error) {
	requirements := args["requirements"]

	content := fmt.Sprintf(`You are an expert BDD practitioner using the Saturday Framework.
The user wants to generate Gherkin scenarios for the following requirements:

"%s"

Please generate a complete .feature file content that includes:
1.  **Feature**: A clear, value-oriented feature title and description.
2.  **Background**: (Optional) Common setup steps.
3.  **Scenarios**: 3-5 scenarios covering Happy Path, Edge Cases, and potentially Negative Path.
4.  **Tags**: Appropriate tags like @sanity, @regression, @smoke.

Follow these best practices:
- Use declarative steps (e.g., "Given I am logged in" vs "Given I enter username...").
- Keep scenarios independent.
- Use Scenario Outline for data-driven tests.
`, requirements)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleImplementFeature(args map[string]string) ([]mcp.PromptMessage, error) {
	feature := args["feature"]

	content := fmt.Sprintf(`You are acting as an **Autonomous QA Engineer** using the Saturday Framework.
Your goal is to implement and verify the following feature: **"%s"**

Please follow this strict workflow (Chain-of-Thought):

**Phase 1: Visualization & Modeling**
1.  Ask the user for a UI Screenshot if not provided.
2.  Use the `+"`visual_page_object`"+` prompt instructions (mentally) to analyze the image.
3.  Generate the Page Object code.
4.  Write the file to `+"`lib/pages`"+`.

**Phase 2: Test Generation**
1.  Create a Playwright spec file in `+"`tests/`"+`.
2.  Import your new Page Object.
3.  Write a test that verifies the feature flows.

**Phase 3: Verification (Self-Healing)**
1.  Run the test using `+"`run_tests`"+`.
2.  **IF PASS**: Celebrate!
3.  **IF FAIL**:
    a. Use `+"`parse_test_failure`"+`.
    b. Analyze the error.
    c. **FIX** the code (Page Object or Test) using `+"`replace_file_content`"+`.
    d. **RETRY** the test.

**Phase 4: Safety Check**
1.  Run `+"`analyze_impact`"+` on your new files.
2.  Ensure no regressions in existing code.

Begin by confirming you understand the feature and asking for the Visual Input (Screenshot).
`, feature)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleSaturdaySME(args map[string]string) ([]mcp.PromptMessage, error) {
	content := `You are the Saturday Framework Subject Matter Expert (SME).
You are an expert in writing tests, implementing page objects, and debugging issues using the Saturday test automation framework.

**Core Principles:**
1.  **Page Object Model**: UI elements are encapsulated in Page classes extending from `+"`@/core`"+`. Avoid raw selectors in tests.
2.  **Fluent Flows**: Reusable sequences of interactions should be organized into Flows.
3.  **Clean Step Definitions**: Gherkin steps should be declarative and map cleanly to Flows/Pages.
4.  **Observability built-in**: Use the custom OTel logger and metrics properties for visibility.
5.  **Boy Scout Rule**: Always leave code better than you found it. Refactor aggressively.

When generating code or answering questions, always adhere to these rules and the Saturday framework conventions.`

	return []mcp.PromptMessage{
		{
			Role: "system",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleMigrateTest(args map[string]string) ([]mcp.PromptMessage, error) {
	legacyCode := args["legacy_code"]
	content := fmt.Sprintf(`You are the Saturday Framework Migration Specialist.
Your task is to migrate the following legacy test code to the Saturday framework:

Legacy code:
%s

**Instructions:**
1.  **Understand Intent:** What is the test trying to achieve?
2.  **Playwright MCP Integration:** Check if the Playwright MCP tools (e.g. `+"`playwright_navigate`"+`, `+"`playwright_evaluate`"+`) are available in your environment. If they are, use them to navigate to the application being tested, inspect the live DOM, and locate the most robust selectors (e.g. user-facing roles). If they are not available, rely entirely on the provided legacy code.
3.  **Generate Page Object:** Create a Saturday Page Object (`+"`lib/pages/XPage.ts`"+`) extending from `+"`@/core/Page`"+`.
4.  **Generate Flow:** Create a Fluent Flow covering the sequence of actions.
5.  **Generate Gherkin:** Generate the Gherkin steps mapping to the flow.
6.  **Refactor & Verify:** Make sure the output adheres to the Boy Scout rule and observability principles.

You MUST use your file editing capabilities to apply the changes directly to the codebase. Do not just output the code in a markdown block.
`, legacyCode)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleOtelMetricsExpert(args map[string]string) ([]mcp.PromptMessage, error) {
	flowCode := args["flow_code"]
	content := fmt.Sprintf(`You are the Saturday OpenTelemetry Expert.
Your task is to analyze the following Saturday Flow code and recommend observability improvements:

Flow Code:
%s

**Instructions:**
1.  **Identify Bottlenecks & Value:** Analyze the flow to understand business value steps or potential performance bottlenecks.
2.  **Suggest Instrumenting:** Suggest where to add custom span metrics or counters to track the flow's success, duration, or failure reasons.
3.  **Code Updates:** Provide the updated Flow code. You MUST use the Saturday framework's built-in hooks and loggers (e.g., `+"`import { logger } from '@orieken/saturday-core'`"+`). Do NOT use raw `+"`@opentelemetry/api`"+` boilerplate.
4.  **Grafana Context:** Explain how these new metrics map to dashboards and how they improve observability.

You MUST use your file editing capabilities to apply the instrumentation directly to the codebase. Do not just output the code in a markdown block.
`, flowCode)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func (p *Provider) handleSelfHealTest(args map[string]string) ([]mcp.PromptMessage, error) {
	failureLog := args["failure_log"]
	content := fmt.Sprintf(`You are the Saturday Self-Healing Orchestrator.
Your task is to fix a failing Saturday test based on the following failure log:

Failure Log:
%s

**Instructions:**
1.  **Analyze the Failure:** Review the log to understand the context of the failure (e.g., stale element, changed ID, timing issue).
2.  **Playwright MCP Integration:** Check if the Playwright MCP tools (like `+"`playwright_navigate`"+`, `+"`playwright_evaluate`"+`, `+"`playwright_click`"+`) are available. If they are, use them to:
    a. Launch the application or navigate to the failing page.
    b. Inspect the live DOM to find the currently valid locators (preferring user-facing roles).
    c. Reproduce the failing step to confirm the issue.
    If Playwright MCP tools are not available, rely entirely on the provided failure log and context.
3.  **Apply Fix:** You MUST use your file editing capabilities to update the local code (e.g., the Saturday Page Object or test file) directly. Avoid hacky fixes like hardcoded sleep (`+"`page.waitForTimeout`"+`). Do not just output the code in a markdown block.
4.  **Verify:** Re-run the test (or use Playwright MCP) to confirm the fix works.
`, failureLog)

	return []mcp.PromptMessage{
		{
			Role: "user",
			Content: mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}
