
# ADR-0001 — CLI Language Choice

## Decision
Implement the initial CLI in **TypeScript/Node** within this monorepo.

## Status
Accepted — 2025-11-04

## Context
- Tight integration with Playwright and existing TS codebase.
- Faster iteration for developer tooling; easy distribution via npm.
- Rich ecosystem (yargs/commander, execa, chalk) and shared types.

## Consequences
- Simpler developer experience now. If/when we need a static, single-binary
  distribution for air-gapped environments, we can either:
  1) Add **pkg**/**nexe** bundling for Node, or
  2) Re-implement the CLI in **Go**, calling the same exporter package over a stable protocol/contract.
