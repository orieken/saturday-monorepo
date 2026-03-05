# Architecture

## High-Level Architecture
(See mermaid diagram in `docs/diagrams/high_level.md` or embed below)

```mermaid
graph TB
    subgraph "Client Layer"
        A[Claude Desktop]
        B[VS Code]
        C[Other MCP Clients]
    end

    subgraph "Saturday MCP Server"
        D[MCP Server Core]
        D --> E[Tool Handler]
        D --> F[Resource Provider]
        D --> G[Prompt Provider]

        E --> H[Generation Engine]
        E --> I[Analysis Engine]
        E --> J[Validation Engine]

        F --> K[Template Store]
        F --> L[Documentation Store]
        F --> M[Example Store]

        H --> N[Template Processor]
        H --> O[AST Generator]
        H --> P[File Writer]

        I --> Q[Code Parser]
        I --> R[Pattern Analyzer]

        J --> S[Schema Validator]
        J --> T[Pattern Validator]
    end

    subgraph "Data Layer"
        U[(Templates)]
        V[(Patterns)]
        W[(Examples)]
        X[(Schemas)]
    end

    A --> D
    B --> D
    C --> D

    N --> U
    R --> V
    M --> W
    S --> X
```

## Component Architecture
- `cmd/` — entry point
- `internal/server/` — MCP server, handlers, tools, resources, prompts
- `internal/templates/` — loader, processor, registry, cache
- `internal/generators/` — site, page, flow, element, service, steps
- `internal/analyzers/` — framework, patterns, metrics
- `internal/validators/` — schema, pattern, naming
- `internal/models/` — request/response/schemas
- `internal/utils/` — fs, ast, logging, errors

