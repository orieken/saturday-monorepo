# DevOps Report: reagent-map-globe

## Status: PASS

## Summary

This slice adds a feature-flagged, offline Cesium 3D globe view to `apps/ye-olde-magic-shop`, pinned to the exact `cesium@1.146.0` dependency already resolved by `packages/saturday-playwright-cesium`. No backend, database, Docker, Kubernetes, Terraform, or GitHub Actions changes are required — the touched surface is one Vite config, one shell script, `package.json`/`pnpm-lock.yaml`, `.env`, and `playwright.config.ts`. All of it is additive and reversible: deleting the flag from `.env` or reverting the branch removes the feature entirely with no cleanup step.

The one open item is pre-existing, not introduced by this feature: `apps/ye-olde-magic-shop` has no GitHub Actions workflow at all today (`.github/workflows/packages.yml` only builds/tests `packages/*`; `saturday-mcp-release.yml` is unrelated). The app's only automated gate is the local pre-ship script, `scripts/ui-suite-in-chrome.sh`. Two fitness functions `architecture-notes.md` assigns to devops-engineer assume a CI job that runs `npm run build` for this app — that job does not exist, so both were verified manually below instead. See **Known Gaps** for the recommendation; wiring new CI is Approval Gate #7 and is not something this report does unilaterally.

## CI Pipeline

- [x] **New test commands**: No CI change needed — `tests/reagent-map-*.spec.ts` run under the existing `chrome` Playwright project, driven by the pre-existing `scripts/ui-suite-in-chrome.sh`. `playwright.config.ts`'s `chromium` project now carries `testIgnore: ['**/reagent-map-*.spec.ts']` so `npm run test:e2e` (chromium) is unaffected by their presence.
- [x] **Verified locally**: Ran `scripts/ui-suite-in-chrome.sh` (see `pre-ship-checks.md`) — exit 0, 40/40 passed in 31.1s, including all `reagent-map-*.spec.ts` specs. This confirms AC6 (reagent-map specs run only under `chrome`) end-to-end.
- [ ] **Not verified in this sandbox**: `npm run test:e2e` (chromium project) run in isolation to directly observe its exit code with the new `testIgnore` entry present. The `testIgnore` glob is config-level and inspected by code-reviewer; recommend CI (once it exists for this app, see Known Gaps) runs both projects and gates on both exit codes per `analysis.md`'s DevOps Task 2.
- [ ] **Build steps**: No new build artifacts beyond `dist/cesium` (see Deployment section) — no new build *step* added, the existing `closeBundle` hook in `vite.config.ts` handles it.
- [ ] **Lint/type-check**: No new tool configs added; nothing to wire.
- [ ] **IaC linting**: N/A — no Terraform/CloudFormation/K8s files touched.
- [x] **CI secrets documentation**: No new secret or credential introduced. `VITE_FEATURE_REAGENT_MAP` is a plain, non-sensitive build flag — no vault entry needed.

## Dependencies

- `apps/ye-olde-magic-shop/package.json` gains `cesium: "1.146.0"` (exact, no caret) and `@orieken/saturday-playwright-cesium: workspace:*` (devDependency).
- `pnpm-lock.yaml` updated accordingly.
- **Fitness function verified** (architecture-notes.md, owner: devops-engineer): *exactly one resolved `cesium` version exists across the workspace.* Confirmed by inspection — `pnpm-lock.yaml` has a single `cesium@1.146.0` resolution block (packages: section, line 2472) and a single matching snapshot (line 7290); `node_modules/.pnpm` at the workspace root shows one `cesium@1.146.0` directory, symlinked into `apps/ye-olde-magic-shop/node_modules/cesium`. No second version anywhere in the tree.
- No Docker base image exists for this app (Vite dev server / static build) — nothing to update.

## Database / Migrations

Not applicable. No schema, migration, or data-layer change in this slice.

## Environment Configuration & Observability

- **New env var**: `VITE_FEATURE_REAGENT_MAP` — build-time Vite flag read only by `isReagentMapEnabled()` in `src/featureFlags.ts`. `apps/ye-olde-magic-shop/.env` now sets it to `false` (off by default, matches Definition of Done). Values: unset / empty / anything other than the exact literal string `"true"` → feature off.
- `apps/ye-olde-magic-shop` has no `.env.example` file in this repo today — there is no existing convention to extend, so none was added. If one is introduced later, it should include `VITE_FEATURE_REAGENT_MAP=false`.
- No secret or credential involved — this is a plain boolean build flag, safe in any CI environment variable, no vault entry required.
- No OTel/observability instrumentation added by this slice — out of scope per `analysis.md`.

## Deployment & Infrastructure

- **Feature flag / dark deploy**: Already satisfied by design — `VITE_FEATURE_REAGENT_MAP=false` in `.env` means this ships dark by default; the route and nav link don't exist at all (not merely hidden) when the flag is off, so there's zero runtime surface to roll back if something is wrong post-deploy. Flipping the flag to `true` is the only activation step, and flipping it back to `false` is the complete rollback.
- **Static asset copy (build-time)**: `vite.config.ts`'s `cesiumStaticAssetsPlugin` copies `node_modules/cesium/Build/Cesium` → `dist/cesium` in its `closeBundle` hook, serving Cesium's prebuilt workers/widgets/assets same-origin without committing them to source control.
  - **Fitness function** (architecture-notes.md, owner: devops-engineer): *`dist/cesium` contains the Cesium `Build` tree after `npm run build`.* Could not be executed end-to-end in this sandbox — `npm run build` fails here on a pre-existing, platform-specific issue unrelated to this feature: the optional native binary `@rollup/rollup-darwin-arm64` is missing from this worktree's `node_modules` (a known npm optional-dependencies bug, https://github.com/npm/cli/issues/4828). This is a local macOS/ARM install-state problem, not a defect in the feature or the plugin, and would not reproduce on CI's Linux runners.
  - Verified instead by direct inspection: the plugin's `closeBundle` guard (`if (!existsSync(CESIUM_SOURCE_DIR)) return;` then `cpSync(..., recursive: true)`) is correct, and the source tree it copies — `apps/ye-olde-magic-shop/node_modules/cesium/Build/Cesium` — exists and contains `Assets/`, `Widgets/`, `Workers/`, `Cesium.js`, confirmed via `ls`. The copy will succeed on any environment where `npm run build` itself runs.
  - Recommend whoever next runs a clean `npm run build` for this app (CI, once it exists, or a fresh local install) confirms `dist/cesium/Widgets` exists per `analysis.md`'s DevOps Task — this report could not be the one to close that loop here.
- **Health checks**: No HTTP health-check endpoint exists or is needed for this static frontend app; unaffected.
- **Immutable infra**: N/A — this app has no server-side deployment artifact (container image, systemd unit) in this repo; it is built as static assets. No manual patching involved.

## Rollback Plan

Low risk, cheap rollback: the feature flag defaults to off, so no deploy-time rollback is anticipated. If the merged code itself needs reverting, `git revert` on this branch's commit(s) removes the route, nav link, dependency, and Vite plugin cleanly — no migration, no data to reconcile, no stateful resource to tear down.

## Known Gaps / Follow-ups

1. **No CI exists for `apps/ye-olde-magic-shop` today.** `.github/workflows/packages.yml` covers only `packages/*`; nothing in `.github/workflows/` runs this app's `npm run test:e2e`, `npm run test:chrome`, or `npm run build`. This predates this feature and is out of this order's `scope.paths`. Two fitness functions in `architecture-notes.md` (single resolved `cesium` version; `dist/cesium` present after build) name "the CI job that runs `npm run build` for this app" as their enforcement point — that job doesn't exist, so both were verified manually in this report instead of automatically. **Recommend**: a follow-up order to add a GitHub Actions workflow for `apps/ye-olde-magic-shop` running `npm run test:e2e`, `npm run test:chrome` (or the equivalent of `scripts/ui-suite-in-chrome.sh`), and `npm run build` with a `dist/cesium/Widgets` existence check — gated on both exit codes per `analysis.md`'s DevOps Task 2. This is a new fitness-function wiring and requires explicit human approval ("approve fitness function" or "add to CI") per `approval-gates.md` gate #7; it is not wired here.
2. **`npm run build` could not be executed in this sandbox** due to a pre-existing, environment-specific missing optional native dependency (`@rollup/rollup-darwin-arm64`), unrelated to this feature. The `closeBundle` copy step was verified by code inspection and source-directory existence instead of a live build. Recommend the next environment that successfully runs `npm run build` for this app confirms `dist/cesium/Widgets` exists, closing out architecture-notes.md's fitness function with a real execution rather than inspection.
3. Carried over from `tech-writer.md` / `security-report.md` (no DevOps action required, noted for completeness): one LOW-severity finding that `vite.config.ts`'s dev-only Cesium static-asset middleware can throw on malformed percent-encoding in a request path — dev-server-only, no production exposure, tracked as tech debt.

## Manual Steps Required

None. No secrets to set, no migration to run, no deployment to trigger.