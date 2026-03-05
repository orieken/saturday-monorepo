# Using Saturday MCP Server with Ye Olde Magic Shop

This guide shows you how to use the Saturday MCP Server to generate code for this project using AI assistants like Claude.

## What is Saturday MCP?

The Saturday MCP (Model Context Protocol) Server is an AI-powered code generation tool that helps you create Saturday framework code faster and more consistently. It integrates with Claude Desktop to generate:

- Site classes
- Page Objects
- Flow classes
- Cucumber step definitions
- And more!

## Setup

### 1. Build the MCP Server

```bash
# From the monorepo root
cd saturday-mcp
go build -o bin/saturday-mcp ./cmd/saturday-mcp
```

### 2. Configure Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "saturday": {
      "command": "/Users/YOUR_USERNAME/Projects/Rieken/saturday-monorepo/saturday-mcp/bin/saturday-mcp"
    }
  }
}
```

**Important**: Replace `YOUR_USERNAME` with your actual username!

### 3. Restart Claude Desktop

Close and reopen Claude Desktop to load the MCP server.

## Usage Examples

### Generate a Site Class

**You ask Claude:**
> "Use the Saturday MCP to generate a Site class for the Magic Shop with pages: home, products, cart, checkout"

**Claude will:**
1. Call the `generate_site` tool
2. Return the generated TypeScript code
3. You can save it to `lib/sites/magic-shop-site.ts`

**Example Request:**
```json
{
  "name": "magicShop",
  "baseUrl": "http://localhost:3000",
  "pages": ["home", "products", "cart", "checkout"],
  "description": "Magic shop e-commerce site"
}
```

**Generated Code:**
```typescript
import { BaseSite } from '@orieken/saturday-core';
import { Page } from 'playwright';
import { HomePage } from './pages/home-page';
import { ProductsPage } from './pages/products-page';
import { CartPage } from './pages/cart-page';
import { CheckoutPage } from './pages/checkout-page';

export class MagicShopSite extends BaseSite {
  constructor(page: Page, baseUrl: string) {
    super(page, baseUrl);
  }

  protected initializePages(): void {
    this.registerPage('home', new HomePage(this.page));
    this.registerPage('products', new ProductsPage(this.page));
    this.registerPage('cart', new CartPage(this.page));
    this.registerPage('checkout', new CheckoutPage(this.page));
  }

  protected initializeFlows(): void {
    // Add flows as needed
  }
}
```

### Generate a Page Object

**You ask Claude:**
> "Generate a ProductDetailsPage with elements: productTitle, price, addToCartButton, quantity"

**Claude will:**
1. Use the `generate_page` tool (coming in TODO-005)
2. Return a complete Page Object class

### Generate Step Definitions

**You ask Claude:**
> "Generate Cucumber step definitions for the checkout.feature file"

**Claude will:**
1. Read your feature file
2. Generate matching step definitions
3. Include proper imports and type safety

## Available Tools

Check what tools are available:

**You ask Claude:**
> "List all available Saturday MCP tools"

**Response:**
```json
[
  {
    "name": "list_tools",
    "description": "List all available Saturday framework tools",
    "status": "implemented"
  },
  {
    "name": "generate_site",
    "description": "Generate a new Site class with page and flow registration",
    "status": "implemented"
  },
  {
    "name": "generate_page",
    "description": "Generate a new Page class with element registration",
    "status": "planned"
  }
]
```

## Tips

### 1. Be Specific
Instead of: "Generate a page"  
Try: "Generate a ProductPage for /products with elements: searchBox (#search), filterButton (.filter-btn), productGrid (.products)"

### 2. Include Context
"Generate a Site class for the Magic Shop e-commerce application with home, products, cart, and checkout pages. The base URL is http://localhost:3000"

### 3. Iterate
Start with basic generation, then ask Claude to:
- Add more elements
- Add custom methods
- Add validation
- Add error handling

### 4. Review Generated Code
Always review the generated code before committing:
- Check naming conventions
- Verify selectors
- Add any custom logic
- Run tests

## Troubleshooting

### MCP Server Not Appearing in Claude

1. Check the path in `claude_desktop_config.json` is correct
2. Ensure the binary exists: `ls -la /path/to/saturday-mcp/bin/saturday-mcp`
3. Make it executable: `chmod +x /path/to/saturday-mcp/bin/saturday-mcp`
4. Restart Claude Desktop completely

### "Command Not Found" Error

Use absolute paths in the configuration, not relative paths like `~/`.

### Generation Errors

Check the Claude Desktop logs:
```bash
tail -f ~/Library/Logs/Claude/mcp*.log
```

## Next Steps

1. Generate your Site class
2. Generate Page Objects for each page
3. Generate Flows for common user journeys
4. Generate Step Definitions for your features
5. Run your tests!

## Learn More

- [Saturday MCP Documentation](../../saturday-mcp/README.md)
- [Saturday Core Documentation](../../packages/saturday-core/README.md)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
