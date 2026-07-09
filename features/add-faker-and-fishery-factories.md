---
type: feature
status: ready
created: 2026-07-09
---

# Add Faker and Fishery Test Data Factories to Saturday Core

> Verified against the real `saturday-monorepo` codebase before writing this, not assumed. Confirmed:
> all seven Site-Centric primitives (`BaseSite`, `BasePage`, `BaseElement`, `BaseFlow`, `BaseFilter`,
> `SiteManager`, `TabManager`) are real, working implementations in `packages/saturday-core/src/`;
> Cucumber.js is genuinely wired (real `.feature` files, real `SaturdayWorld` integration); OpenTelemetry
> is real and working for Cucumber scenarios (via `saturday-cucumber-otel-formatter`) and Playwright Test
> runs (via `saturday-playwright-otel-reporter`, fully complete per its own `TODO_PLAYWRIGHT_OTEL.md`).
> This spec touches none of that. Confirmed gap: zero `@faker-js/faker` or `fishery` usage anywhere in
> the monorepo, despite `shared/rules/typescript-conventions.md` (in `ai-assistant-dot-files`)
> specifying both as this stack's fake-data and factory convention.
>
> Explicitly NOT addressed here, on purpose: OpenTelemetry coverage for the framework's own Jest/Vitest
> unit tests is a real, separate gap — but it's already tracked by the maintainers themselves in
> `TODO_JEST_OTEL.md` and `TODO_VITEST_OTEL.md` (both 0% started). This spec doesn't duplicate that
> planning; it targets the one gap not already on your own roadmap.

## Summary
Give anyone writing a Saturday test two things they don't have today: a standard way to generate
realistic fake data instead of hand-typing string literals, and a factory pattern for building complex
test fixtures without repeating the same object-construction code in every test file.

## Acceptance Criteria
- [ ] Given a test author needs a realistic fake value (a name, email, address, etc.), when they use
      `@faker-js/faker`, then it's available as a real dependency in at least `saturday-core` — not
      something each package or app has to add for itself.
- [ ] Given a test author needs to build a complex fixture (e.g. test data for a `BasePage` flow), when
      they define a factory with `fishery`'s `Factory.define()`, then calling it produces a fully
      populated, realistic object with zero manual field-by-field construction.
- [ ] Given the existing test suite (Vitest in `saturday-core`/`saturday-cucumber`/etc., Jest in the OTel
      and k6 packages — confirmed mixed, not a single standard today), when this feature ships, then
      every existing test still passes unmodified — this is additive, not a rewrite.
- [ ] Given a new consumer of `saturday-core` wants an example to follow, when they look at the package,
      then there's at least one real, working factory (e.g. a `UserFactory`) demonstrating the pattern,
      not just the two libraries added to `package.json` with no usage anywhere.

## Out of Scope
- OpenTelemetry instrumentation for the framework's own Jest/Vitest unit test runs — already tracked in
  `TODO_JEST_OTEL.md` and `TODO_VITEST_OTEL.md`. Don't duplicate that planning here.
- Standardizing every package onto Vitest (currently a mixed Vitest/Jest split across packages) — a real
  inconsistency, but a separate concern from adding test-data tooling, and a bigger, more disruptive
  change that deserves its own spec if it's worth doing at all.
- Filling in `BaseFlow`'s implementation or building concrete `Filters` beyond what already exists —
  both already have real (if thin) implementations; not what this spec is about.
- The "Friday Platform" integration described in `PROJECT_OVERVIEW.md` — the referenced
  `friday-platform/`/`docs/friday/` paths don't exist in this repo; out of scope and unverifiable from
  here regardless.

## Domain Language
- **Factory**: per `shared/rules/typescript-conventions.md` (in `ai-assistant-dot-files`) — a
  `fishery`-based test-data builder. Not to be confused with the Gang-of-Four Factory pattern.
- No other new terms — this implements an already-documented convention, it doesn't introduce new
  domain concepts.

## Non-Functional Requirements
- Neither `@faker-js/faker` nor `fishery` may become a runtime dependency of any non-test code path —
  both are test-tooling only, and must not end up in a production bundle if `saturday-core` is ever
  consumed outside a test context.
- Factories must not import Playwright, Cucumber, or any adapter-specific package — a factory produces
  plain data; it has no business knowing which test runner or browser driver is in use.

## Trust Boundaries
None — test-tooling code, not production request-handling code. No user input, no external network
calls.

## Test Approach
Saturday/Playwright and Saturday/Cucumber both apply here — this is core library tooling consumed by
both. Use whichever test runner the specific package already uses (Vitest for `saturday-core`, matching
its existing `vitest.config.ts`) rather than introducing a third runner.

## Open Questions
- **Where does the example factory live?** Recommend: `packages/saturday-core/src/factories/`, as a new
  subdirectory alongside the existing `elements/` — but confirm with whoever owns `saturday-core`'s
  internal structure, since this is a new top-level concern for that package.
- **Should `saturday-core` re-export `faker` directly, or just depend on it internally for its own
  example factory?** Recommend: don't re-export — let consumers add their own `@faker-js/faker`
  dependency directly, since re-exporting a third-party library through your own package surface creates
  a version-pinning obligation you don't need to take on for a testing utility.

## Infrastructure / Deploy Notes
- New dependencies: `@faker-js/faker` and `fishery`, added to `packages/saturday-core/package.json`
  (dev dependency, given the Non-Functional Requirement above).
- No env vars, no migrations, no deploy-sequence changes.

## Definition of Done
- [ ] All acceptance criteria verified
- [ ] All existing tests still pass, plus new tests for the example `UserFactory`
- [ ] `packages/saturday-core`'s own README (if one exists) or the monorepo root README's package list
      updated to mention the new factory pattern is available
- [ ] CI green (note: this repo currently has no CI workflow computing a test-count/coverage badge —
      "CI green" here means the existing `pnpm -r test` passes, not a badge update)
