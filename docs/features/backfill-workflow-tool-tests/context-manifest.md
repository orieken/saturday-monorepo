# Context Manifest: backfill-workflow-tool-tests

## 1. Scope and Boundaries
- **Target Component**: saturday-mcp/internal/tools/workflow_tool.go
- **Relevant Layers**: Adapter (Tool interface wrapping Domain Workflow interface)
- **Bounded Context**: MCP Server Tooling & Workflows

## 2. Pinpoint Files (To Keep Open)
- [workflow_tool.go](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/internal/tools/workflow_tool.go#L1-L47) -- Source under test: WorkflowTool struct, NewWorkflowTool, Name, Description, InputSchema, OutputSchema, Execute
- [testfixtures_test.go](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/internal/tools/testfixtures_test.go#L1-L116) -- Package-level test helpers (buildRequest, extractText, silentLogger)
- [run_tests_workflow_test.go](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/internal/workflows/run_tests_workflow_test.go#L1-L300) -- Reference fake pattern and workflow testing style
- [workflow.go](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/internal/domain/workflow.go#L1-L50) -- domain.Workflow interface contract

## 3. Global Rules and Constraints
- [ARCHITECTURE_RULES.md](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/ARCHITECTURE_RULES.md)
- [DOMAIN_DICTIONARY.md](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/DOMAIN_DICTIONARY.md)

## 4. Knowledge Items & ADRs (To Load)
- [ADR-001](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-001-use-invopop-jsonschema-tool-output-schemas.md) -- OutputSchema contract
- [ADR-002](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-002-default-otlp-grpc-otel-trace-export.md) -- Tracing middleware context

## 5. Prior Deliveries in This Bounded Context
- No prior feature archive deliveries in this bounded context.

## 6. Prune Recommendations (To Close)
- None

## 7. Token Budget
- **Estimated total tokens for pinned files**: ~500 tokens
- **Target agent tier**: Analyst/Developer ≤80%
- **Status**: OK
