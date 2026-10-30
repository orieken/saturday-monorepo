---
intake: 1
id: runner-image-polyglot
type: chore
plan: deliver-feature
status: draft
priority: 4
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [ui-suite-real-chrome]
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [containers/saturday-runner/**, docs/**]
---

# Add the polyglot toolchains to the runner image

## Request

> The container should have all the tools we would need. Saturday's ubi-vnc image carries Java,
> Python and Gradle for the framework's Java and Python ports.

## Summary

A `polyglot` target on top of each runner variant, for suites written in other languages.
The variants stay lean for the factory's TypeScript runs.

## Acceptance Criteria

- [ ] Given `containers/saturday-runner/Containerfile`, then a `polyglot` target builds `FROM debian` (and `ubi9-polyglot` `FROM ubi9` once that variant exists)
      and adds a JDK, Python, Gradle and Go, each at an exact version, from its release archive,
      verified against a recorded SHA-256 before it is unpacked.
- [ ] Given the `polyglot` image, then `runner-check` also records each toolchain's version in
      `/work/evidence/`, and fails if any is missing.
- [ ] Given `docs/runner-image.md`, then it lists the toolchain versions and how to bump one.
- [ ] Given the change is complete, then the `debian` and `ubi9` package lists are unchanged.

## Out of Scope

- Building any toolchain from source.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- None

## Trust Boundaries

- Toolchain archives are accepted only when their SHA-256 matches the recorded value.

## Test Approach

`runner-check` in the polyglot image.

## Open Questions

- Which versions: Assumption if unresolved: the current LTS JDK (Temurin), Python 3.12, the Gradle
  version Saturday's Java port uses, and Go matching `saturday-mcp/go.mod`.

## Infrastructure / Deploy Notes

- None
