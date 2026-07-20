# Architectural Decision Records (ADRs)

This directory contains records of architectural decisions made for the `saturday-mcp` project.

## Index

* [ADR-001: Use invopop/jsonschema for Tool Output Schemas](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-001-use-invopop-jsonschema-tool-output-schemas.md) (Accepted, 2026-07-20)
* [ADR-002: Default to OTLP over gRPC for OTel Trace Export](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-002-default-otlp-grpc-otel-trace-export.md) (Accepted, 2026-07-20)

---

## ADR Template

To create a new ADR, copy the template below and save it as `docs/adrs/ADR-[NNN]-[kebab-case-title].md`.

```markdown
# ADR-[NNN]: [Title — active voice: "Use X" not "Decision to use X"]

Date: YYYY-MM-DD
Status: Accepted | Proposed | Deprecated | Superseded by [ADR-XXX](link)
Deciders: [List of deciders]
Technical Story: [Link to PR, ticket, or phase of work]

## Context
[What situation, constraint, or force prompted this decision?]

## Decision
[What was decided — specific, concrete, active voice]

## Alternatives Considered
| Option | Why Rejected |
|---|---|
| [Alternative 1] | [Reason for rejection] |

## Consequences
- **Easier**: [what this unlocks]
- **Harder**: [what this constrains]
- **Changed**: [what is different going forward]

## Fitness Function
[How this decision is enforced (e.g., CI checks, lint rules, unit tests) — or "Judgment-only: [reason]"]
```
