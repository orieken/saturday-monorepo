---
feature: "reagent-map-visual-baselines"
bounded_context: "Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)"
domain_terms: ["Reagent Site", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/tests/", "apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/playwright.config.ts", "apps/ye-olde-magic-shop/.playwright-baselines/", "scripts/ui-suite-in-chrome.sh", "scripts/ui-suite-in-chrome.test.sh", "docs/testing-cesium.md, docs/runner-image.md"]
issue_refs: ["reagent-map-visual-baselines", "reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# Architecture Notes: reagent-map-visual-baselines

## Structural Decisions

### reagent-map-visual.spec.ts lives alongside the existing reagent-map-*.spec.ts siblings in apps/ye-olde-magic-shop/tests/, not in a new directory or package.
- **Rationale**: playwright.config.ts's chromium project already testIgnores '**/reagent-map-*.spec.ts', and the 'chrome' project (via scripts/ui-suite-in-chrome.sh) already picks up every file matching that glob. Matching the naming convention gets the new spec into both projects' routing for free — no config change needed for routing.
- **Fitness Function**: judgment-only — playwright.config.ts's chromium project already testIgnores '**/reagent-map-*.spec.ts', and the 'chrome' project (via scripts/ui-suite-in-chrome.sh) already picks up every file matching that glob. Matching the naming convention gets the new spec into both projects' routing for free — no config change needed for routing.
### Add a visible, stable error-state element to ReagentMapView.vue: a div with data-testid="reagent-map-error" and role="alert", set via a new ref immediately before the existing `throw sitesResult.reason` in onMounted, rendered in the template guarded by that ref.
- **Rationale**: No error-state element exists today (onMounted only rethrows). This is additive DOM in the existing template — it does not touch buildViewer(), the camera, newViewer.destroy(), or the Promise.allSettled control flow reagent-map-resilience.spec.ts already asserts against. Naming matches the existing data-testid="reagent-site-list" convention; role="alert" matches the accessibility pattern already used on #cesiumContainer (role="img"/aria-label).
- **Trade-offs**:
  - Touches a file reagent-map-resilience.spec.ts already covers in detail — any edit here re-risks that spec's zero-canvas/zero-list-item assertions if not placed carefully before the throw.
  - Reversibility: Cheap — it's additive markup with no behavioral coupling; removing the element later does not affect the throw/destroy control flow.
- **Fitness Function**: reagent-map-resilience.spec.ts's existing assertions (seam absent, zero canvas, zero list items) still pass unchanged after the element is added, and a new assertion confirms the element is visible with role=alert when the sites fetch rejects. — verified by Playwright run of reagent-map-resilience.spec.ts (existing) plus the new error-state @visual test in reagent-map-visual.spec.ts, both under --project=chrome.
### playwright.config.ts gets a new exported resolveRunnerVariant(): string | undefined that reads process.env.RUNNER_VARIANT at call time, and the config's snapshotPathTemplate is set to '.playwright-baselines/{arg}{ext}' where {arg} composes resolveRunnerVariant() with Playwright's own snapshot name token.
- **Rationale**: This mirrors the exact shape of resolveRunnerSlowMo()/resolveRunnerVideo() already in this file — exported, read-env-at-call-time, unit-tested in tests/playwright-config.spec.ts. No new config pattern is introduced; RUNNER_VARIANT just joins the two existing per-call env resolvers.
- **Fitness Function**: resolveRunnerVariant() resolves to 'debian' or 'ubi9' when RUNNER_VARIANT is set to that value, and to undefined (or a value that makes the @visual skip condition true) when unset — read fresh per call, not cached at config load. — verified by tests/playwright-config.spec.ts gains cases for resolveRunnerVariant() in the same block structure as the existing resolveRunnerSlowMo/resolveRunnerVideo test cases (lines ~68-98).
### apps/ye-olde-magic-shop/.playwright-baselines/{debian,ubi9}/ is a new committed directory holding only baseline PNGs, written exclusively by a --update-baselines-triggered Playwright run. No new package.
- **Rationale**: Keeps baselines colocated with the app they belong to, per-runner-image, consistent with how playwright-report and other app-local test artifacts are already scoped under apps/ye-olde-magic-shop/.
- **Fitness Function**: judgment-only — Keeps baselines colocated with the app they belong to, per-runner-image, consistent with how playwright-report and other app-local test artifacts are already scoped under apps/ye-olde-magic-shop/.
### scripts/ui-suite-in-chrome.sh's parse_args must be restructured from a single-positional `case "${1:-}"` into a loop over "$@" so --update-baselines and --ubi9 compose in either order, each setting an independent flag (VARIANT/IMAGE_TAG/PNPM_STORE_VOLUME for --ubi9; a new UPDATE_BASELINES=1 for --update-baselines).
- **Rationale**: The current parse_args only inspects $1 and treats any second value as an error — it cannot accept two flags today. This is a required correctness change, not an optional cleanup: without it, '--ubi9 --update-baselines' hits the existing unknown-argument branch and exits 64.
- **Trade-offs**:
  - Widens parse_args's surface — must keep it under the complexity-6 budget (design-principles.md); split into one small helper per flag if the loop body grows past a simple case per token.
  - Reversibility: Cheap — purely internal to one script, covered immediately by ui-suite-in-chrome.test.sh's existing exit-code-contract pattern.
- **Fitness Function**: Both orderings (--ubi9 --update-baselines and --update-baselines --ubi9) set VARIANT=ubi9/IMAGE_TAG=saturday-runner:ubi9 AND enable baseline updates; --update-baselines alone keeps the debian default; an unrecognized third argument still exits 64 with the existing usage message. — verified by bash scripts/ui-suite-in-chrome.test.sh — new test_update_baselines_composes_with_ubi9 case(s) using the existing FAKE_DOCKER_LOG assertion pattern.
### start_container() exports RUNNER_VARIANT into the container (alongside the existing RUNNER_SLOW_MO/RUNNER_VIDEO -e flags), and run_chrome_suite() takes the update-baselines flag as a parameter, appending --update-snapshots to the `npx playwright test --project=chrome --workers=3` invocation only when it is set.
- **Rationale**: Matches the existing -e RUNNER_SLOW_MO=... / -e RUNNER_VIDEO=... pattern in start_container() exactly. Gating --update-snapshots behind an explicit parameter (not a free env-var read inside run_chrome_suite) keeps the 'never write a baseline without the flag' invariant checkable by reading one call site.
- **Fitness Function**: RUNNER_VARIANT=debian|ubi9 appears in the docker run env flags for every invocation; --update-snapshots appears in the docker exec'd playwright command if and only if --update-baselines was passed. — verified by bash scripts/ui-suite-in-chrome.test.sh — FAKE_DOCKER_LOG assertions, extending the existing test_ubi9_variant/test_default_variant_is_debian pattern with a RUNNER_VARIANT check and a new case for --update-snapshots presence/absence.
### Every @visual test in reagent-map-visual.spec.ts calls test.skip(!process.env.RUNNER_VARIANT, '<reason>') (or an equivalent beforeEach guard) as its first statement — never a conditional inside the assertion body — so a missing baseline fails loudly whenever RUNNER_VARIANT is set, and the whole file no-ops cleanly on a developer Mac.
- **Rationale**: AC5 requires skip-when-unset and never-skip-when-set as a hard behavioral contract, not an incidental side effect of some other check. A single, consistently-placed skip call per test is the only shape that is easy to verify by reading the file and easy for code-reviewer to catch a violation of.
- **Fitness Function**: Running the suite with RUNNER_VARIANT unset reports every @visual test as 'skipped' with a stated reason; running it with RUNNER_VARIANT set and a baseline deliberately removed reports that test as 'failed', never 'skipped'. — verified by Two manual Playwright runs recorded in the delivery record: one with RUNNER_VARIANT unset (expect all @visual tests skipped), one with RUNNER_VARIANT set and a baseline temporarily removed (expect that test failed, not skipped).
### The globe canvas test's tolerance is a single named constant (e.g. GLOBE_MAX_DIFF_PIXEL_RATIO = 0.01) passed into toHaveScreenshot(), never a literal repeated at the call site, and this value is never edited upward — per AC3, a flaky globe test is removed and masked instead of its tolerance being raised.
- **Rationale**: A named constant makes the 1% ceiling a single line code-reviewer can check on every future diff to this file, rather than a convention that has to be remembered. Making this a wired CI fitness function (e.g. a grep-based lint rule) is not proposed here — Approval Gate #7 would require wiring a new CI check for a one-constant invariant that code-reviewer can already catch by inspection every time this file changes, which is a worse cost/benefit than the judgment-only route.
- **Fitness Function**: judgment-only — A named constant makes the 1% ceiling a single line code-reviewer can check on every future diff to this file, rather than a convention that has to be remembered. Making this a wired CI fitness function (e.g. a grep-based lint rule) is not proposed here — Approval Gate #7 would require wiring a new CI check for a one-constant invariant that code-reviewer can already catch by inspection every time this file changes, which is a worse cost/benefit than the judgment-only route.

## Component Placement

- `reagent-map-visual.spec.ts` → test layer, `apps/ye-olde-magic-shop/tests/` — Matches the existing reagent-map-*.spec.ts sibling files' location and naming; picked up automatically by both Playwright projects' existing glob rules.
- `Error-state element (data-testid="reagent-map-error")` → framework layer, `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Vue view template markup — UI/framework layer, not domain logic. Lives in the same single-file component that already owns the onMounted error path it renders for.
- `resolveRunnerVariant() + snapshotPathTemplate` → config layer, `apps/ye-olde-magic-shop/playwright.config.ts` — Same file and export pattern as the existing resolveRunnerSlowMo()/resolveRunnerVideo() env resolvers it mirrors.
- `.playwright-baselines/{debian,ubi9}/ baseline PNGs` → test layer, `apps/ye-olde-magic-shop/.playwright-baselines/` — App-local Playwright test fixtures, colocated with the app and keyed per runner image; not a shared or reusable artifact.
- `--update-baselines flag, RUNNER_VARIANT export, run_chrome_suite extra-args` → tooling layer, `scripts/ui-suite-in-chrome.sh` — Same script that already owns --ubi9 and the container lifecycle; this is an additive flag on existing lifecycle functions, not a new script.
- `New --update-baselines / RUNNER_VARIANT test cases` → test layer, `scripts/ui-suite-in-chrome.test.sh` — Same stubbed-docker shell test harness that already covers --ubi9 and every other flag/exit-code branch of the script under test.
- `Baseline-regeneration procedure documentation` → docs layer, `docs/testing-cesium.md, docs/runner-image.md` — Closes testing-cesium.md's existing 'No visual baselines' Known Gaps entry and extends runner-image.md's existing 'Updating the browsers (both locks)' section with the matching baseline-regeneration step.

## Bounded Context

- **Owning Context**: Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)
- **Context Crossings**: Test Automation (Saturday Framework / Playwright @visual baselines), Infrastructure (Saturday runner image's debian/ubi9 variants — docs/runner-image.md)

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [PASS] ReagentMapView.vue (framework/UI layer) must not gain new imports of adapter or domain-layer code as part of the error-state element. — The decision adds only a ref and template markup inside the existing component; no new import is introduced. Verified by direct read of ReagentMapView.vue (lines 1-158) rather than the loom graph — mcp__loom__dependencies_of/describe were not reachable in this session (no 'not_found'/'fresh' response was returned at all, suggesting the loom MCP server is not registered here), so this is a manual import-reading check, not a graph query.
- [PASS] playwright.config.ts's new RUNNER_VARIANT resolver must not introduce config-load-time env coupling (must read process.env fresh per call, like its two existing siblings). — Confirmed by reading playwright.config.ts lines 14-21: resolveRunnerSlowMo/resolveRunnerVideo are already exported, called-at-use-site functions, not top-level consts. resolveRunnerVariant() must follow the identical shape.
- [PASS] scripts/ui-suite-in-chrome.sh's new flag must not change the suite's run_chrome_suite invocation (--workers=3) or start_container's --shm-size=1g, per the feature's stated NFR. — The structural decisions above only add an env export and a conditional extra argument to run_chrome_suite; --workers=3 and --shm-size=1g remain untouched literals in the script as read.

## Anti-Pattern Check

- God object: not found — ReagentMapView.vue stays a single-purpose view component; the error-state addition is a few lines of template/ref, not new responsibilities.
- Shotgun surgery: not found — Each change is isolated to one file per concern (view component, config, script, two docs) with no change rippling across unrelated packages.
- Leaky abstraction: not found — snapshotPathTemplate and RUNNER_VARIANT stay inside playwright.config.ts and the script that sets the env var; no test file needs to know the path template's internal shape, only toHaveScreenshot()'s own name argument.
- Premature generalization: not found — No new abstraction framework introduced — RUNNER_VARIANT is one more call-time env resolver added next to two that already exist, not a new generalized settings layer.

## Fitness Functions

- **resolveRunnerVariant() resolves correctly and is read fresh per call, not cached at config load**
  - **Verification**: tests/playwright-config.spec.ts new cases mirroring the existing resolveRunnerSlowMo/resolveRunnerVideo tests
  - **Owner**: Developer
- **reagent-map-resilience.spec.ts's existing zero-canvas/zero-list-item/seam-absent assertions remain unchanged and passing after the error-state element is added**
  - **Verification**: Playwright run of reagent-map-resilience.spec.ts under --project=chrome
  - **Owner**: QA
- **--update-baselines composes with --ubi9 in either argument order; without --update-baselines, --update-snapshots is never passed to playwright**
  - **Verification**: bash scripts/ui-suite-in-chrome.test.sh
  - **Owner**: DevOps
- **RUNNER_VARIANT is exported into the container in start_container() for every invocation (debian and ubi9)**
  - **Verification**: bash scripts/ui-suite-in-chrome.test.sh FAKE_DOCKER_LOG assertion
  - **Owner**: DevOps
- **Every @visual test skips with a stated reason when RUNNER_VARIANT is unset, and never skips (missing baseline fails) when it is set**
  - **Verification**: Two manual Playwright runs recorded in the delivery record, per AC5's own acceptance wording
  - **Owner**: QA
- **Cyclomatic complexity below 7 on every changed/new function, including the restructured parse_args**
  - **Verification**: analyze-complexity tooling against the diff
  - **Owner**: Architect

## Refactoring Opportunities

- Extract Function / loop-ify conditional on `scripts/ui-suite-in-chrome.sh:parse_args` — Single-positional case statement cannot express two composable flags; needs a loop over "$@", ideally with one small helper per flag to stay under the complexity-6 budget.
- Introduce Parameter on `scripts/ui-suite-in-chrome.sh:run_chrome_suite` — Currently a hardcoded command string; needs to accept whether to append --update-snapshots rather than reading a global/env flag inline, so the 'never write a baseline without the flag' invariant is checkable at the call site in main().

## Developer Handoff Notes

- playwright.config.ts: add `export function resolveRunnerVariant(): string | undefined` reading process.env.RUNNER_VARIANT at call time (same shape as resolveRunnerSlowMo/resolveRunnerVideo at lines 14-21), and set `snapshotPathTemplate: '.playwright-baselines/{arg}{ext}'`-style config using it. Add matching cases to tests/playwright-config.spec.ts right after the existing resolveRunnerSlowMo/resolveRunnerVideo blocks (lines ~68-98).
- ReagentMapView.vue: add the error ref and its guarded template element immediately around the existing `throw sitesResult.reason` at line 142 — do not touch buildViewer(), addSiteEntities(), or onBeforeUnmount's destroy/seam-clear logic, since reagent-map-resilience.spec.ts asserts against all three.
- scripts/ui-suite-in-chrome.sh: parse_args (lines 25-38) must become a loop over "$@" so --update-baselines and --ubi9 compose in either order; start_container (lines 57-66) needs RUNNER_VARIANT added to its existing -e flag list; run_chrome_suite (lines 133-136) needs a parameter controlling whether --update-snapshots is appended.
- ui-suite-in-chrome.test.sh: extend write_fake_docker's FAKE_DOCKER_LOG capture (already used by test_ubi9_variant/test_default_variant_is_debian) to add update-baselines/RUNNER_VARIANT coverage — reuse the existing assert_contains helper, don't add a new assertion mechanism.
- docs: close testing-cesium.md's 'No visual baselines' Known Gaps entry (lines 193-194) and cross-reference docs/runner-image.md's existing 'Updating the browsers (both locks)' section with the regeneration procedure — both per shared/rules/writing-style.md.

## Open Architectural Questions

- Carried from analysis.md: human confirmation is still needed that adding a visible error-state element to ReagentMapView.vue (none exists today) is in scope — the architecture decision above treats it as additive UI chrome, not a change to the Cesium view's camera or behavior, but the spec's 'Changing the view itself is out of scope' wording could be read either way and that reading is a product/spec call, not a structural one.
- Whether the #cesiumContainer @visual test survives the five-run flakiness check in both runner images is unresolved until the developer actually runs it; AC3 already specifies the fallback (remove it, mask the canvas in the other @visual tests, record why) if it doesn't.
- The loom MCP graph tools (mcp__loom__dependencies_of/dependents_of/describe) were not reachable in this session — no response, not even a 'not_found'. Layer-boundary checks above were done by direct file reading instead, as the fallback this agent's own instructions describe; this should be verified against the graph if/when the loom MCP server is registered.
