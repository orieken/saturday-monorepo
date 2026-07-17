---
type: feature
status: ready
created: 2026-07-17
---

# Add BasePartial for Cross-Page Shared UI

> Verified against the real `saturday-monorepo` codebase before writing this, not assumed. Confirmed:
> `packages/saturday-core/src/base/base-page.ts`, `base-element.ts`, `base-flow.ts`, `base-filter.ts`,
> `base-site.ts` all exist and are actively used. Specific elements (`ButtonElement`, `InputElement`,
> `LinkElement`) live in `packages/saturday-core/src/elements/`. Zero existing `partial`/`SharedHeader`/
> `GlobalHeader` concept anywhere in `packages/saturday-core/src/` — confirmed by grep. This is the gap
> `ai-assistant-dot-files`' updated `docs/patterns/saturday-framework-patterns.md` identified: `BaseElement`
> is documented as "within a single BasePage," leaving cross-page shared UI (headers, footers, global
> nav) with no clean home. This spec closes that gap in the TypeScript reference implementation.

## Summary
Give test authors a first-class abstraction for shared UI that appears across multiple pages (headers,
footers, global navigation) without duplicating locators in every `BasePage` subclass or awkwardly
stretching `BaseElement` semantics beyond a single page's scope.

## Acceptance Criteria
- [ ] Given a test author needs to model shared UI that appears on multiple pages, when they extend
      `BasePartial`, then they get a class shaped similarly to `BasePage` (owns locators, exposes
      interaction methods, can host its own `Filter`-gated behavior) but explicitly not tied to a
      single page's lifecycle.
- [ ] Given a `BasePage` subclass composes a `BasePartial`, when a test accesses that partial through
      the page (e.g. `page.header()`), then the same partial instance is available consistently across
      every page that composes it — no test needs to reconstruct the partial's locators per page.
- [ ] Given a `BasePartial` needs its own state gate (e.g. a `GlobalHeader` behaves differently when
      logged in vs. logged out), when the partial declares a `@RequiresFilter` decorator on a method,
      then that filter fires exactly like it does on a `BasePage` method — the partial's Filter behavior
      is not a special case.
- [ ] Given the existing test suite across `packages/saturday-core` and `apps/ye-olde-magic-shop`, when
      this feature ships, then every existing test still passes unmodified — this is additive, not a
      rewrite of `BasePage` or `BaseElement`.
- [ ] Given a developer looks for the pattern to follow, when they open `packages/saturday-core/src/`,
      then there's at least one real, working `BasePartial` subclass (e.g. `GlobalHeaderPartial`)
      demonstrating composition into a `BasePage`, not just the abstract base class with no consumer.

## Out of Scope
- Retrofitting existing `apps/ye-olde-magic-shop` pages to use partials — a separate migration once
  the base concept exists. Would balloon this spec into a per-page refactor otherwise.
- Adding `BasePartial` support to the ML/heatmap/OTel packages (`packages/saturday-core/src/ml`,
  `saturday-cucumber-otel-formatter`, etc.) — partials are a Site-Centric concept, not an observability
  concept. If a partial happens to be worth heatmapping, the existing element-level ML mechanisms already
  cover its underlying elements.
- Any framework port to another language (C#, Python, Java) — each has its own separate feature spec.
- A "PartialManager" analogous to `SiteManager` — partials are page-scoped-in-composition, not
  application-scoped. No manager needed.

## Domain Language
- **BasePartial** (per updated `shared/DOMAIN_DICTIONARY.md`): a shared UI section (header, footer,
  global nav) that appears across multiple `BasePage`s. Lives in a `partials/` subdirectory;
  deliberately does NOT follow the `FooPage` naming convention. Fills the gap `BaseElement` (which is
  "within a single BasePage") doesn't cover. Synonyms to AVOID: `SharedComponent`, `LayoutFragment`,
  `PartialPage`.
- No other new terms.

## Non-Functional Requirements
- `BasePartial` must not import or depend on any specific `BasePage` subclass — the composition goes
  page-to-partial only, never the other direction (a partial that knows which pages use it is a
  circular design smell).
- File placement: base class at `packages/saturday-core/src/base/base-partial.ts` (alongside
  `base-page.ts`); concrete partials at `packages/saturday-core/src/partials/` (parallel to the
  existing `src/elements/` convention).
- Cyclomatic complexity < 7 per method (framework rule), coverage ≥ 85% for the new class (framework
  rule).

## Trust Boundaries
None — test-tooling code, not production request-handling code. No user input, no external network
calls.

## Test Approach
Vitest for unit tests on `BasePartial` itself (mirroring `base-page.spec.ts` / `base-element.spec.ts`
if those exist, matching the package's established test style). Real Playwright test against a fixture
site (`apps/ye-olde-magic-shop`) exercising the worked-example `GlobalHeaderPartial` composed into two
different pages — proves the "same partial, multiple pages" claim end-to-end, not just at unit-test
level.

## Open Questions
- **Should `BasePartial` extend `BaseElement` or be its own root class?** Recommend: its own root
  class. A partial composes multiple elements internally and has its own lifecycle concerns (state
  filters that apply to the partial as a whole); making it a subclass of `BaseElement` would inherit
  the "wraps one Playwright locator" model that doesn't fit. Confirm by comparing to how `BasePage`
  itself relates to `BaseElement` today — if `BasePage` isn't a `BaseElement`, `BasePartial` shouldn't
  be either.
- **Should partials be composed into pages via a decorator, a base-class hook, or plain accessor
  methods?** Recommend plain accessor methods (`header(): GlobalHeaderPartial`) — matches how
  `BasePage` already exposes its elements today per its own convention, keeps the mechanism
  transparent, and doesn't introduce a new framework primitive just for composition.
- **Naming convention for concrete partials — `FooPartial` (mirrors `FooPage`/`FooElement`) or bare
  `Foo` (relies on directory placement for signal)?** Recommend `FooPartial` suffix for grep-ability
  and consistency with the framework's other naming conventions. The updated pattern doc says partials
  "deliberately do NOT follow the `FooPage` naming convention" — meaning don't call them `FooPage`;
  the `Partial` suffix is a distinct signal, not the same thing.

## Infrastructure / Deploy Notes
- No new dependencies — uses Playwright + existing `saturday-core` primitives.
- No env vars, no migrations, no deploy-sequence changes.
- `packages/saturday-core/package.json` may need a new entry-point export for `BasePartial` and any
  partials/ barrel — match how `BasePage` / `BaseElement` are currently exported.

## Definition of Done
- [ ] All acceptance criteria verified
- [ ] `BasePartial` and one worked-example concrete partial (e.g. `GlobalHeaderPartial`) exist and
      compile
- [ ] Unit tests written and passing, coverage ≥ 85% on the new files
- [ ] At least one `apps/ye-olde-magic-shop` page composes the worked-example partial, exercised by an
      existing or new Playwright test that passes
- [ ] `pnpm -r build` and `pnpm -r test` both pass with no regressions
- [ ] `packages/saturday-core`'s README (or the monorepo root README) updated to mention `BasePartial`
      alongside `BasePage` and `BaseElement`
