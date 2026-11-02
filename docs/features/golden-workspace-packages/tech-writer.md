# Documentation Report: golden-workspace-packages

```yaml
feature: "golden-workspace-packages"
bounded_context: "Test Automation"
domain_terms: ["SaturdayConfig"]
files_touched: []
issue_refs: ["intake/golden-workspace-packages.md", "#7 (fix/cart-spec-magicshop-import)", "salvage/golden-workspace-packages @ f372a30"]
linked_adrs: []
linked_kis: []
```

## Status

PASS — the delivery record below is complete. One acceptance criterion (AC4) is not met; it is recorded as an open gap for human decision, not resolved by this stage.

## Summary

This feature has no user-facing documentation surface. It moves `apps/ye-olde-magic-shop` from three vendored, drifted copies of the framework to real `workspace:*` links against `packages/saturday-core`, `packages/saturday-cucumber`, and `packages/saturday-playwright-heatmap`, and fixes one confirmed behavior gap (`browser.headless`'s default source) inside `packages/saturday-cucumber`. No public API, CLI flag, or configuration surface changed in an end-user-visible way: `SaturdayConfig.browser.headless` was already a documented boolean field; only where its default value comes from changed, from a hardcoded literal to the same `process.env.HEADLESS` read the vendored copy already performed. Per the analyst's own instruction in `analysis.md`'s Tech Writer Tasks, this report states explicitly: no README, API doc, or ADR change is required by this feature.

## Documentation Checklist

### Always Update

- **CHANGELOG.md** — No `CHANGELOG.md` exists anywhere in this repository outside `node_modules`. There is no changelog convention to append to, and this feature does not introduce one; adding a repo-wide changelog file is a separate decision outside this feature's scope.
- **README.md** — Checked `apps/ye-olde-magic-shop/README.md` and `packages/saturday-cucumber/README.md` directly. Neither mentions the vendored directories, the removed Dockerfiles, or a `headless`/`HEADLESS` configuration section. Neither file needs a change: the app's quick-start steps already describe `pnpm install` plus `pnpm --filter` commands that work unchanged under `workspace:*`, and the cucumber package's README never documented `defaultConfig` or `browser.headless` in the first place, so there is no stale sentence to correct.

### Update if Applicable

- **API documentation** — Not applicable. No endpoint or public API surface changed.
- **Configuration docs** — Not applicable. `packages/saturday-cucumber/README.md` has no existing Configuration section describing `SaturdayConfig`'s fields or defaults, so there is no stale value to fix. Writing a new Configuration section is a larger documentation task than this feature authorizes; the spec's Out of Scope list ("Renaming files, adding flows, changing step definitions. Those are later orders.") and `analysis.md`'s own Tech Writer Task 2 both point to recording this explicitly rather than expanding scope.
- **Architecture Decision Record** — Not written. `architecture-notes.md` records both structural decisions (the `workspace:*` switch, the `HEADLESS` default fix) as following an already-established in-file pattern, with no new package, abstraction, or layer boundary. Neither decision meets the bar for an ADR.
- **Getting started / setup guides** — Not applicable. `pnpm install` from the repository root is unchanged; the three dependencies resolve as workspace symlinks automatically.
- **Migration guides** — Not applicable. No breaking change to any consumer-facing contract.

### Code-Level Docs

- **Module-level docstrings** — Not applicable. No new module was created.
- **Function/method docstrings** — Not applicable. No new public function or method was introduced; the one changed line is a default-value expression inside an existing object literal.
- **Type hints / JSDoc** — `code-review-report.md` raised one non-blocking item: `packages/saturday-cucumber/tests/config.spec.ts` carries no issue-ref/AC-ref JSDoc annotation per the Test Annotation Convention in `shared/rules/testing-conventions.md`. This is feedback directed at the developer and marked non-blocking; it is not addressed in this stage and is carried forward below as a follow-up.

## Delivery Record

### Before / After Suite Results

**Before** (documented starting condition, per the spec and corroborated independently by the developer): `scripts/ui-suite-in-chrome.sh` and the app's `test:bdd`/`test:chrome` suites fail at collection, because the vendored `saturday-core`/`saturday-cucumber`/`saturday-playwright-heatmap` copies have no `dist/` after a fresh install.

**After**, per `qa-report.md`:
- `apps/ye-olde-magic-shop` Cucumber suite (`test:bdd`): 16/16 scenarios, 85/85 steps passed.
- Dockerized Chrome Playwright suite, via `scripts/ui-suite-in-chrome.sh`: 24/24 specs passed (magic-shop, map, cart, heatmap, webgl, playwright-config). First attempt failed with an unrelated `Cannot find module @rollup/rollup-linux-x64-gnu` error, traced to a stale `ui-suite-pnpm-store` Docker volume predating this change; after removing that volume, the suite completed all 16/16 build steps and all 24/24 Playwright tests, exit 0.
- `packages/saturday-cucumber` unit suite: 22/22 passing, including the new `config.spec.ts` (97.0% statement coverage).
- `packages/saturday-core` unit suite: 166/166 passing (89.2% statement coverage).
- `packages/saturday-playwright-heatmap` unit suite: 6/6 passing (100.0% statement coverage).
- Combined test run total: 234 passed, 0 failed, 0 skipped.
- No `.feature` file or spec assertion was changed by this feature, confirmed by `git status`/`git diff`.

### Dockerfile Decision

- `apps/ye-olde-magic-shop/Dockerfile.runner` is removed. No workflow, script, or compose file in the repository references it; `scripts/ui-suite-in-chrome.sh` builds from `containers/saturday-runner/Containerfile` instead. This removal has no open issue.
- `apps/ye-olde-magic-shop/Dockerfile` is removed, matching the salvaged commit's own choice — but QA found a real consumer the developer's literal-string grep missed: `prototype-runner/docker-compose.yml`'s `web-app` service builds with `context: ../apps/ye-olde-magic-shop` and no explicit `dockerfile:` key, so Compose looks for the now-deleted implicit default `Dockerfile`. QA reproduced the failure directly (`docker compose build web-app` in `prototype-runner/` fails with `open Dockerfile: no such file or directory`).
  - Fixing this the right way — rewriting `apps/ye-olde-magic-shop/Dockerfile` to build from the workspace — requires changing `prototype-runner/docker-compose.yml`'s build context, a file outside this run's `scope.paths`. Making that change would itself violate AC6's "no file outside `scope.paths` changes" constraint.
  - This is recorded as an **open gap**, not resolved by this pipeline run. It needs a human or architect decision: either approve a `scope.paths` extension to include `prototype-runner/docker-compose.yml` so the Dockerfile can be rewritten against the workspace, or confirm `prototype-runner` is already stale — its own `WORKFLOW.md` describes a `final/` directory layout removed by the "extract starter kits to standalone repos" refactor (commit `3f650ab`) — and accept this as pre-existing drift the ship gate should note but not block on.

### Open Questions Resolution

1. **"A vendored file has behaviour the package lacks"** — Resolved. The one confirmed gap, `defaultConfig.browser.headless`'s default source in `packages/saturday-cucumber/src/config/index.ts`, was fixed before any vendored file was deleted, with a new test (`config.spec.ts`) shown to fail against the pre-fix hardcoded value and pass after. Beyond that one case, all 29 deleted vendored test files were diffed against their `packages/*/tests/**` counterparts: 27 are byte-identical, and the remaining two (`saturday-cucumber/tests/coverage.spec.ts`, `world.spec.ts`) are strict supersets in `packages/` — no vendored coverage was lost, so nothing further needed to move.
2. **"A spec fails through the script for a reason unrelated to the vendored packages"** — Not encountered. All 24 Chrome Playwright specs passed, including `map.spec.ts`, so this assumption path was never triggered and the spec/script were not touched.

### Acceptance Criteria Status (for the record)

- AC1 (workspace:* dependencies) — PASS
- AC2 (vendored directories removed, no coverage lost) — PASS
- AC3 (every prior-passing scenario still passes) — PASS
- AC4 (Dockerfiles build from workspace, or removed with reason recorded) — **NOT MET** — `Dockerfile.runner`'s removal is clean; `Dockerfile`'s removal broke `prototype-runner/docker-compose.yml`'s `web-app` service (see Dockerfile Decision above).
- AC5 (`scripts/ui-suite-in-chrome.sh` exits 0 on a fresh clone) — PASS
- AC6 (complexity, coverage, scope.paths discipline) — PASS

## Known Gaps / Follow-ups

- **AC4 gap** — `apps/ye-olde-magic-shop/Dockerfile`'s removal broke `prototype-runner/docker-compose.yml`'s `web-app` build. Carried forward for a human/architect decision: extend `scope.paths` to rewrite the Dockerfile against the workspace, or accept `prototype-runner` as already-stale drift. Do not treat this run as closing AC4 until that decision is made.
- **Test annotation** — `packages/saturday-cucumber/tests/config.spec.ts` has no issue-ref/AC-ref JSDoc per the Test Annotation Convention. Non-blocking per code review; still open.

## Files Updated by This Stage

None. No README, CHANGELOG, ADR, or code-level doc met the bar for a change — see Documentation Checklist above for why each was considered and skipped. The delivery-record content above is written into this report for persistence to `docs/features/golden-workspace-packages/` by the orchestrator, per this feature's own scope exception for that path.

## Notes for DevOps

- `pnpm-lock.yaml` was regenerated via a clean `pnpm install` against current `main`, not reused byte-for-byte from the salvaged branch. Confirm no existing CI step pins to the prior lockfile hash.
- The AC4 gap affects `prototype-runner/docker-compose.yml`'s `web-app` build only. Confirm whether this environment is active in CI before the ship gate, since this run does not resolve it.