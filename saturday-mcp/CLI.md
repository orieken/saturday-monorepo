# Saturday CLI

Command-line interface for the Saturday Framework MCP server.

## Installation

### From Source

```bash
cd saturday-mcp
go build -o bin/saturday ./cmd/cli
```

### Add to PATH (Optional)

```bash
# Add to your shell profile (.bashrc, .zshrc, etc.)
export PATH="$PATH:/path/to/saturday-mcp/bin"
```

## Usage

### Generate Commands

Generate code artifacts for your test automation project.

#### Generate a Page Object

```bash
# Basic page
saturday generate page LoginPage --path /login

# Page with elements
saturday generate page LoginPage \
  --path /login \
  --elements "username:#user:input,password:#pass:input,submit:#login-btn:button"

# Write to file
saturday generate page LoginPage \
  --path /login \
  --elements "username:#user:input" \
  --write \
  -o ./my-project
```

#### Generate a Flow

```bash
saturday generate flow CheckoutFlow \
  --steps "Navigate to cart,Review items,Enter payment,Confirm order"

# Write to file
saturday generate flow CheckoutFlow \
  --steps "Step 1,Step 2,Step 3" \
  --write
```

#### Generate a Service

```bash
saturday generate service UserService \
  --base-url "https://api.example.com" \
  --endpoints "getUser:GET:/users/:id,createUser:POST:/users,deleteUser:DELETE:/users/:id"

# Write to file
saturday generate service UserService \
  --base-url "https://api.example.com" \
  --endpoints "getUser:GET:/users/:id" \
  --write
```

#### Generate a Site

```bash
saturday generate site MyAppSite \
  --url "https://myapp.com" \
  --pages "Home,Login,Dashboard" \
  --flows "LoginFlow,CheckoutFlow"

# Write to file
saturday generate site MyAppSite \
  --url "https://myapp.com" \
  --pages "Home,Login" \
  --write
```

### Analyze Commands

Analyze existing code for framework usage and issues.

#### Analyze Framework Structure

```bash
# Scan project
saturday analyze framework ./my-project

# JSON output
saturday analyze framework ./my-project --json
```

#### Analyze Performance

```bash
# Find performance issues
saturday analyze performance ./my-project

# JSON output
saturday analyze performance ./my-project --json
```

### Suggest Improvements

```bash
# Get improvement suggestions
saturday suggest ./my-project

# JSON output
saturday suggest ./my-project --json
```

### Validate Patterns

```bash
# Validate code patterns
saturday validate ./my-project

# Strict mode (fail on warnings)
saturday validate ./my-project --strict

# JSON output
saturday validate ./my-project --json
```

### Migrate Legacy Code

```bash
# Migrate a legacy test file to Page Object
saturday migrate page ./legacy-test.ts

# Write to file
saturday migrate page ./legacy-test.ts --write -o ./my-project
```

### Generate Documentation

```bash
# Generate project documentation
saturday docs ./my-project ./docs/PROJECT_DOCS.md
```

## Command Reference

### Global Flags

- `-v, --verbose` - Enable verbose output
- `-o, --output` - Output directory (default: current directory)
- `--help` - Show help for any command
- `--version` - Show version information

### Generate Subcommands

- `page [name]` - Generate a Page Object
- `flow [name]` - Generate a Flow class
- `service [name]` - Generate a Service class
- `site [name]` - Generate a Site class

### Analyze Subcommands

- `framework [path]` - Analyze framework structure
- `performance [path]` - Analyze performance issues

### Other Commands

- `suggest [path]` - Suggest code improvements
- `validate [path]` - Validate code patterns
- `migrate page [file]` - Migrate legacy code
- `docs [project-path] [output-file]` - Generate documentation

## Examples

### Complete Workflow

```bash
# 1. Generate a new site
saturday generate site MyShop \
  --url "https://shop.example.com" \
  --pages "Home,Products,Cart,Checkout" \
  --flows "PurchaseFlow,SearchFlow" \
  --write

# 2. Generate pages
saturday generate page ProductsPage \
  --path /products \
  --elements "search:#search:input,filter:.filter-btn:button" \
  --write

# 3. Validate the code
saturday validate ./lib

# 4. Analyze performance
saturday analyze performance ./lib

# 5. Generate documentation
saturday docs ./lib ./docs/API.md
```

### CI/CD Integration

```bash
#!/bin/bash
# validate.sh - Run in CI pipeline

set -e

echo "Validating Saturday framework patterns..."
saturday validate ./lib --strict

echo "Checking for performance issues..."
saturday analyze performance ./lib --json > performance-report.json

echo "Generating documentation..."
saturday docs ./lib ./docs/GENERATED.md

echo "✓ All checks passed!"
```

## Output Formats

### Default (Pretty Print)

Human-readable output with colors and formatting.

### JSON

Machine-readable output for scripting and CI/CD:

```bash
saturday analyze framework ./lib --json | jq '.stats'
```

## Tips

1. **Use `--write` carefully** - Always review generated code before committing
2. **Combine with version control** - Generate to a branch and review diffs
3. **Automate validation** - Add `saturday validate` to your CI pipeline
4. **Document as you go** - Run `saturday docs` regularly to keep docs updated
5. **Migrate incrementally** - Use `saturday migrate` to gradually adopt patterns

## Troubleshooting

### Command not found

Make sure the `bin` directory is in your PATH or use the full path:

```bash
/path/to/saturday-mcp/bin/saturday --help
```

### Permission denied

Make the binary executable:

```bash
chmod +x bin/saturday
```

### Template errors

Ensure all templates are loaded correctly. The CLI uses the same template system as the MCP server.

## Development

### Building

```bash
go build -o bin/saturday ./cmd/cli
```

### Testing

```bash
# Test a command
./bin/saturday generate page TestPage --path /test

# Test with actual project
./bin/saturday analyze framework ../ye-olde-magic-shop
```

## See Also

- [Saturday MCP Server README](../README.md)
- [MCP Protocol Documentation](https://modelcontextprotocol.io)
- [Saturday Framework](https://github.com/orieken/saturday-core)
