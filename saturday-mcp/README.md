# Saturday MCP Server

A Model Context Protocol (MCP) server for the Saturday testing framework that provides code generation, analysis, and validation tools.

## Overview

The Saturday MCP Server enables AI assistants (like Claude) to generate and analyze Saturday framework code through a standardized protocol. It provides tools for:

- **Code Generation**: Sites, Pages, Flows, Elements, Services, Step Definitions
- **Analysis**: Framework structure, pattern compliance, metrics
- **Validation**: Schema validation, pattern validation, naming conventions

### 🎯 Platform Support

**Works with multiple IDEs and AI assistants:**
- ✅ **Claude Desktop** (recommended)
- ✅ **VS Code** (with GitHub Copilot, v1.102+)
- ✅ **Cursor IDE**
- ✅ **JetBrains IDEs** (IntelliJ, PyCharm, WebStorm, etc. - v2025.2+)
- ⚠️ **Antigravity** (likely supported)
- ✅ **Windsurf**

📖 **[Full IDE Compatibility Guide](./docs/IDE_COMPATIBILITY.md)**

## Installation

### Option 1: Download Pre-compiled Binaries (Recommended)
You can download the pre-compiled `saturday-mcp` and `saturday` CLI binaries for macOS, Linux, and Windows from the [GitHub Releases page](https://github.com/orieken/saturday-mcp/releases).

1. Download the archive for your operating system and architecture.
2. Extract the archive and move the binary to a location in your system `$PATH` (e.g., `/usr/local/bin`).

### Option 2: Build from Source
If you have Go 1.22+ installed, you can build the binaries from source:

```bash
cd saturday-mcp
go mod download
go build -o bin/saturday-mcp ./cmd/saturday-mcp
go build -o bin/saturday ./cmd/cli
```

## Running the Server

If installed via binaries (Option 1):
```bash
saturday-mcp
```

If built from source (Option 2):
```bash
./bin/saturday-mcp
```

The server will start and listen for MCP protocol messages on stdin/stdout.

## CLI Usage

The Saturday CLI provides command-line access to all MCP server functionality without requiring an MCP client.

### Installation

```bash
cd saturday-mcp
go build -o bin/saturday ./cmd/cli
```

### Quick Start

```bash
# Generate a Page Object
./bin/saturday generate page LoginPage --path /login --elements "username:#user:input"

# Analyze a project
./bin/saturday analyze framework ./my-project

# Validate patterns
./bin/saturday validate ./my-project

# Generate documentation
./bin/saturday docs ./my-project ./docs/API.md
```

📖 **[Full CLI Documentation](./CLI.md)**

## Development Status

### ✅ TODO-001: Project Setup & MCP Server Bootstrap (COMPLETED)
- [x] Go module initialization
- [x] MCP SDK integration
- [x] Server entry point (`cmd/saturday-mcp/main.go`)
- [x] Logging infrastructure
- [x] Server handler with tool registration
- [x] Basic `list_tools` implementation
- [x] Unit tests for logging and server components

### ✅ TODO-002: Template System Foundation (COMPLETED)
- [x] Template registry with thread-safe operations
- [x] Template loader with embedded filesystem support
- [x] Template processor with caching
- [x] Template cache with TTL and cleanup
- [x] Helper functions (pascalCase, camelCase, snakeCase, kebabCase, etc.)
- [x] Sample templates (page, site)
- [x] Comprehensive unit tests
- [x] Integration tests

### ✅ TODO-003: Input Validation & Schema System (COMPLETED)
- [x] Request models for all generation operations
- [x] Response models for results and errors
- [x] Validator with go-playground/validator
- [x] Custom validation rules (validName, validSelector)
- [x] User-friendly error messages
- [x] JSON schemas for requests
- [x] Comprehensive unit tests
- [x] 56.2% test coverage

### ✅ TODO-004: Site Generator (COMPLETED)
- [x] SiteGenerator implementation
- [x] Integration with template system
- [x] Integration with validation system
- [x] Request validation before generation
- [x] Metadata generation
- [x] Filename generation with kebab-case
- [x] Comprehensive unit tests
- [x] 83.3% test coverage

### ✅ TODO-004.5: MCP Server Integration (COMPLETED)
- [x] Wire up all generators to MCP server
- [x] Add `generate_page` tool handler
- [x] Add `generate_flow` tool handler
- [x] Add `generate_steps` tool handler
- [x] JSON request/response handling
- [x] Error handling and logging
- [x] Tool schema definitions
- [x] Integration with template + validation systems
- [x] File writing support for all tools
- [x] Comprehensive E2E integration tests

### 🎊 Phase 2 Complete!
**All core code generators (Site, Page, Flow, Steps) are implemented, tested, and integrated with MCP!**

### ✅ TODO-008 & 009: Framework Analyzer & Validation Tools (COMPLETED)
- [x] Framework Analyzer (`analyze_framework`)
- [x] Pattern Validator (`validate_patterns`)
- [x] Improvement Suggester (`suggest_improvements`)
- [x] Performance Analyzer (`analyze_performance`)
- [x] Integration with MCP server handlers
- [x] Comprehensive unit tests

### ✅ TODO-010: Resource & Prompt Providers (COMPLETED)
- [x] Prompt Provider (`plan_feature`, `explain_framework`)
- [x] New Prompts (`debug_error`, `generate_gherkin`)
- [x] Resource Provider foundation
- [x] Comprehensive unit tests

### 🚀 Phase 6: The AI Agent Evolution (In Progress)
- [x] **TODO-018: Knowledge Graph** (`analyze_impact`)
- [x] **TODO-019: Test Execution** (`run_tests`)
- [x] **TODO-020: Visual Intelligence** (`visual_page_object` prompt)
- [x] **TODO-021: Observability Integration** (`prioritize_tests`)

### 📋 Planned
- TODO-012-017: Advanced Features (Superseded by Phase 6)

## Architecture


```
saturday-mcp/
├── cmd/
│   └── saturday-mcp/
│       └── main.go              # Server entry point
├── internal/
│   ├── logging/
│   │   ├── logger.go            # Structured logging
│   │   └── logger_test.go
│   ├── server/
│   │   ├── handler.go           # MCP server handler
│   │   ├── testing.go           # Test wrappers
│   │   └── handler_test.go
│   ├── templates/               # ✅ TODO-002
│   │   ├── registry.go          # Template registration
│   │   ├── loader.go            # Template loading from embedded FS
│   │   ├── processor.go         # Template execution
│   │   ├── cache.go             # Template caching
│   │   ├── helpers.go           # Template helper functions
│   │   ├── data/                # Embedded template files
│   │   │   ├── page.tmpl
│   │   │   ├── site.tmpl
│   │   │   ├── flow.tmpl
│   │   │   └── steps.tmpl
│   │   └── *_test.go            # Comprehensive tests
│   ├── validators/              # ✅ TODO-003
│   │   ├── validator.go         # Request validation
│   │   └── validator_test.go
│   ├── models/                  # ✅ TODO-003
│   │   ├── requests.go          # Request models
│   │   ├── responses.go         # Response models
│   │   └── schemas/             # JSON schemas
│   │       ├── page-generation.json
│   │       └── site-generation.json
│   ├── generators/              # ✅ TODO-004-007
│   │   ├── generator.go         # Generator facade
│   │   ├── site_generator.go    # Site class generation
│   │   ├── page_generator.go    # Page class generation
│   │   ├── flow_generator.go    # Flow class generation
│   │   ├── step_generator.go    # Step definition generation
│   │   └── *_test.go            # Comprehensive tests
│   ├── filewriter/              # ✅ TODO-004.6
│   │   ├── filewriter.go        # File writing utilities
│   │   └── filewriter_test.go
│   ├── integration/             # ✅ TODO-004.7
│   │   └── e2e_test.go          # End-to-end tests
│   ├── analyzers/               # TODO-008
│   └── utils/                   # Utilities
├── docs/                        # Architecture & implementation docs
│   ├── IMPLEMENTATION_REVIEW.md
│   ├── TODO-004.5-COMPLETE.md
│   └── IDE_COMPATIBILITY.md
└── go.mod
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run integration tests only
go test ./internal/integration -v
```

### Test Coverage
- **filewriter**: 85.4%
- **generators**: 90.2%
- **logging**: 100.0%
- **templates**: 91.9%
- **validators**: 56.2%

## Available Tools

### ✅ Implemented Tools

#### `list_tools`
List all available Saturday framework tools with their implementation status.

**Example:**
```json
{
  "name": "list_tools",
  "arguments": {}
}
```

#### `generate_site`
Generate a Site class with page and flow registration.

**Example:**
```json
{
  "name": "generate_site",
  "arguments": {
    "name": "myApp",
    "baseUrl": "https://myapp.com",
    "pages": ["home", "dashboard", "profile"],
    "flows": ["login", "checkout"],
    "description": "My application site",
    "writeToFile": true,
    "outputPath": "/path/to/project"
  }
}
```

#### `generate_page`
Generate a Page class with element registration.

**Example:**
```json
{
  "name": "generate_page",
  "arguments": {
    "name": "loginPage",
    "path": "/login",
    "elements": [
      {
        "name": "usernameInput",
        "selector": "#username",
        "type": "input"
      },
      {
        "name": "passwordInput",
        "selector": "#password",
        "type": "input"
      },
      {
        "name": "submitButton",
        "selector": "button[type='submit']",
        "type": "button"
      }
    ],
    "description": "Login page",
    "writeToFile": true,
    "outputPath": "/path/to/project"
  }
}
```

#### `generate_flow`
Generate a Flow class for multi-step user journeys.

**Example:**
```json
{
  "name": "generate_flow",
  "arguments": {
    "name": "checkoutFlow",
    "steps": [
      "addItemToCart",
      "proceedToCheckout",
      "enterShippingInfo",
      "enterPaymentInfo",
      "confirmOrder"
    ],
    "description": "Complete checkout process",
    "writeToFile": true,
    "outputPath": "/path/to/project"
  }
}
```

#### `generate_steps`
Generate Cucumber step definitions from Gherkin patterns.

**Example:**
```json
{
  "name": "generate_steps",
  "arguments": {
    "steps": [
      {
        "type": "Given",
        "pattern": "I am on the login page"
      },
      {
        "type": "When",
        "pattern": "I enter {string} and {string}"
      },
      {
        "type": "Then",
        "pattern": "I should see the dashboard"
      }
    ],
    "language": "typescript",
    "description": "Login feature steps",
    "writeToFile": true,
    "outputPath": "/path/to/project"
  }
}
```

#### `generate_element`
Generate a custom Element/Component class.

**Example:**
```json
{
  "name": "generate_element",
  "arguments": {
    "name": "NavBar",
    "rootSelector": "nav.main-nav",
    "methods": ["clickHome", "search"],
    "description": "Navigation bar component",
    "writeToFile": true,
    "outputPath": "/path/to/project"
  }
}
```

#### `generate_service`
Generate an API Service class.

**Example:**
```json
{
  "name": "generate_service",
  "arguments": {
    "name": "User",
    "baseUrl": "https://api.example.com",
    "endpoints": [
      {
        "name": "GetUser",
        "method": "GET",
        "path": "/users/1"
      },
      {
        "name": "CreateUser",
        "method": "POST",
        "path": "/users"
      }
    ],
    "description": "User management service",
    "writeToFile": true,
    "outputPath": "/path/to/project"
  }
}
```

#### `analyze_framework`
Analyze existing framework structure and patterns.

**Example:**
```json
{
  "name": "analyze_framework",
  "arguments": {
    "projectPath": "/path/to/project"
  }
}
```

#### `suggest_improvements`
Suggest code improvements based on Saturday framework best practices.

**Example:**
```json
{
  "name": "suggest_improvements",
  "arguments": {
    "projectPath": "/path/to/project"
  }
}
```

#### `migrate_code`
Analyze and migrate legacy code to Saturday Framework patterns.

**Example:**
```json
{
  "name": "migrate_code",
  "arguments": {
    "sourceCode": "test('example', async ({ page }) => { await page.click('#submit'); });",
    "type": "page"
  }
}
```

#### `analyze_performance`
Analyze code for performance bottlenecks (concurrent scanning).

**Example:**
```json
{
  "name": "analyze_performance",
  "arguments": {
    "projectPath": "/path/to/project"
  }
}
```

#### `analyze_impact`
Analyze the impact of modifying a specific file (dependency graph).

**Example:**
```json
{
  "name": "analyze_impact",
  "arguments": {
    "projectPath": "/path/to/project",
    "targetFile": "lib/pages/LoginPage.ts"
  }
}
```

#### `run_tests`
Execute tests and capture output.

**Example:**
```json
{
  "name": "run_tests",
  "arguments": {
    "projectPath": "/path/to/project",
    "filter": "login"
  }
}
```

#### `parse_test_failure`
Parse standard output from Playwright tests to identify failing files and lines.

**Example:**
```json
{
  "name": "parse_test_failure",
  "arguments": {
    "output": "Error: ... at tests/login.spec.ts:15:5"
  }
}
```

#### `prioritize_tests`
Rank test coverage needs based on production usage metrics.

**Example:**
```json
{
  "name": "prioritize_tests",
  "arguments": {
    "metricsFile": "/path/to/metrics.json"
  }
}
```

#### `generate_documentation`
Generate markdown documentation for the project.

**Example:**
```json
{
  "name": "generate_documentation",
  "arguments": {
    "projectPath": "/path/to/project",
    "outputPath": "/path/to/output/docs.md"
  }
}
```

#### `validate_patterns`
Validate code against Saturday framework patterns (naming conventions, inheritance).

**Example:**
```json
{
  "name": "validate_patterns",
  "arguments": {
    "projectPath": "/path/to/project"
  }
}
```

### 📚 Available Resources

The server exposes internal templates for reference:

- **`saturday://templates/site`**: Template for Site class
- **`saturday://templates/page`**: Template for Page class
- **`saturday://templates/flow`**: Template for Flow class
- **`saturday://templates/steps`**: Template for Cucumber steps

### 💡 Available Prompts

The server provides pre-defined prompts to help you build faster:

- **`plan_feature`**: Helps you plan the implementation of a new feature (Pages, Flows, Steps).
- **`explain_framework`**: Explains the core architectural concepts of Saturday.
- **`debug_error`**: Analyzes test failures and suggests debugging steps.
- **`generate_gherkin`**: Generates structured BDD scenarios from requirements.
- **`visual_page_object`**: Generates a Page Object from a UI screenshot (requires attaching image).
- **`implement_feature`**: Orchestrates the "Autonomous QA" workflow to implement a feature from scratch.

## MCP Client Configuration

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "saturday": {
      "command": "/path/to/saturday-mcp/bin/saturday-mcp"
    }
  }
}
```

### VS Code

Configure in `.vscode/settings.json`:

```json
{
  "mcp.servers": {
    "saturday": {
      "command": "/path/to/saturday-mcp/bin/saturday-mcp"
    }
  }
}
```

## Contributing

See `docs/` for detailed architecture and implementation guides.

## License

MIT
