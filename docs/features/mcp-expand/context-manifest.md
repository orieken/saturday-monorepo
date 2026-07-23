# Context Manifest: mcp-expand

## 1. Scope and Boundaries
- **Target Component**: saturday-mcp framework expansion
- **Relevant Layers**: Domain (`internal/domain/`), Tools (`internal/tools/`), Workflows (`internal/workflows/`), Prompts (`internal/prompts/`), Adapters (`internal/adapters/`), Server (`internal/server/`)
- **Bounded Context**: Framework-wide MCP Capabilities (Analytical, Search, Personas, Workflows)

## 2. Pinpoint Files (To Keep Open)
- [mcp-expand.md](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/prompts/mcp-expand.md#L1-L98) -- Feature specification prompt
- [mcp-add-plan.md](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/mcp-add-plan.md#L1-L312) -- Prior retrofit plan reference
- [tool.go](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/internal/domain/tool.go#L1-L30) -- domain.Tool interface
- [tool_provider.go](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/internal/server/tool_provider.go#L1-L50) -- Tool registration slice
- [testfixtures_test.go](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/internal/tools/testfixtures_test.go#L1-L116) -- Shared tool test helpers

## 3. Global Rules and Constraints
- [ARCHITECTURE_RULES.md](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/ARCHITECTURE_RULES.md)
- [DOMAIN_DICTIONARY.md](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/DOMAIN_DICTIONARY.md)

## 4. Knowledge Items & ADRs (To Load)
- [ADR-001](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-001-use-invopop-jsonschema-tool-output-schemas.md)
- [ADR-002](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-002-default-otlp-grpc-otel-trace-export.md)

## 5. Prior Deliveries in This Bounded Context
- [backfill-workflow-tool-tests](docs/features/backfill-workflow-tool-tests/)
- [refresh-docs-post-retrofit](docs/features/refresh-docs-post-retrofit/)

## 6. Prune Recommendations (To Close)
- None

## 7. Token Budget
- **Estimated total tokens for pinned files**: ~3000 tokens
- **Target agent tier**: Analyst/Developer ≤80%
- **Status**: OK
