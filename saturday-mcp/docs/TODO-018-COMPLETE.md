# TODO-018: Knowledge Graph & Impact Analysis

## Status: ✅ Complete

We have implemented the foundational "Brain" of the Saturday MCP agent by adding a Dependency Graph Analyzer. This allows the agent to look beyond single files and understand how components relate to each other.

## Implemented Components

### 1. Graph Analyzer (`internal/analyzers/graph_analyzer.go`)
- **Parser**: Scans TypeScript/JavaScript files for `import ... from ...` statements.
- **Node Linking**: Builds a directed graph of Nodes (Files) and Edges (Imports).
- **Resolver**: Intelligently resolves relative paths (e.g., `../pages/LoginPage`) to absolute project paths.
- **Heuristics**: Classification of Nodes into Pages, Flows, Steps, and Elements based on directory structure and naming.

### 2. Impact Analysis Tool (`analyze_impact`)
- **Algorithm**: Uses Breadth-First Search (BFS) to traverse "Incoming" edges (upstream dependencies).
- **Function**: Given a target file (e.g., `LoginPage.ts`), it finds all files that depend on it directly or transitively.
- **Use Case**: "If I change the 'submit' selector in `LoginPage.ts`, which Flow strings and Tests might break?"

## Integration

The tool is accessible via the MCP protocol as `analyze_impact`.

### Example Request
```json
{
  "name": "analyze_impact",
  "arguments": {
    "projectPath": "/path/to/repo",
    "targetFile": "lib/pages/LoginPage.ts"
  }
}
```

### Example Response
```json
{
  "target": "lib/pages/LoginPage.ts",
  "impacted": [
    "lib/flows/LoginFlow.ts",
    "tests/steps/login_steps.ts"
  ],
  "count": 2
}
```

## Testing

We verified the logic with `go test ./internal/analyzers/graph_analyzer_test.go`, which builds a mock project structure and asserts correct dependency resolution and impact tracing.
