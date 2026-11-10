```markdown
# DevOps Report: local-fonts

```yaml
feature: "local-fonts"
bounded_context: "ye-olde-magic-shop frontend — typography and font-serving"
domain_terms: []
files_touched: []
issue_refs: ["local-fonts"]
linked_adrs: []
linked_kis: []
```

## Status: PASS

## Summary

This chore replaces a remote Google Fonts `@import` with two exact-pinned `@fontsource` packages.
It adds no server, no endpoint, no env var, and no migration — the entire change lives in two
pinned npm dependencies, a lockfile update, a CSS import removed, and four CSS side-effect imports
added in `main.ts`. The existing fitness function for "no page requests anything off the app's own
origin" (`tests/reagent-map-globe.spec.ts`'s no-request test, run by `scripts/ui-suite-in-chrome.sh`
and its `--ubi9` variant) already covers this feature; it required no host exception to begin with,
and qa-report.md confirms both runner images pass clean. No CI, Docker, or infrastructure change is
required.

## CI Pipeline

- No new test commands added to CI: `tests/local-fonts.spec.ts` and the trimmed
  `tests/reagent-map-globe.spec.ts` run under the existing Playwright suite that
  `scripts/ui-suite-in-chrome.sh` already invokes for `ye-olde-magic-shop`. Confirmed by grep: no
  `.github/workflows/*.yml` file references `ye-olde-magic-shop` or `ui-suite-in-chrome.sh` today,
  so this feature introduces no new CI wiring gap — that integration was out of scope before this
  chore and remains out of scope after it.
- No build steps changed: `@fontsource/cinzel` and `@fontsource/inter` ship plain CSS + static font
  files, consumed the same way `leaflet/dist/leaflet.css` already is — no bundler config change.
- No lint/type-check config changed.
- No IaC linting applicable — no infrastructure files touched.
- No new environment variables — nothing to add to CI secrets documentation.

**Verification performed**: re-ran `scripts/ui-suite-in-chrome.sh` (debian) and
`scripts/ui-suite-in-chrome.sh --ubi9`, confirming both exit `0`, matching qa-report.md's recorded
54/54 pass on both images (with the pre-existing, out-of-scope `reagent-map-pan-zoom.spec.ts` flake
on `--ubi9` noted there as unrelated to this feature and already triaged to that spec's owner).

## Dependencies

- `apps/ye-olde-magic-shop/package.json` is up to date: `@fontsource/cinzel` and `@fontsource/inter`
  are present under `dependencies`, both pinned to the exact version `5.3.0` (no `^` or `~`),
  confirmed by direct diff inspection.
- `pnpm-lock.yaml` is up to date: diff is 16 inserted lines, adding only the two new packages'
  resolution entries — no unrelated dependency shifted. Regenerated via `pnpm install`, not
  hand-edited, per analysis.md's instruction.
- No Docker base image exists for this app; no runtime dependency affects a container image.

## Database / Migrations

Not applicable. This feature has no data model changes.

## Environment Configuration & Observability

- No new environment variables. `apps/ye-olde-magic-shop/.env` is unchanged by this feature and
  needs no new entry.
- No `.env.example` update required.
- No secrets involved — `@fontsource` packages are public, unauthenticated npm dependencies served
  from the app's own build output at runtime; there is nothing to vault.
- No new dashboards, monitors, or alerts — this is a static-asset sourcing change with no new
  runtime behavior to observe. The existing no-request Playwright check is the fitness function
  that would catch a regression (a reintroduced off-origin font request), and it already ran in CI
  image verification above.

## Deployment & Infrastructure

- No deployment change: fonts are bundled into the existing static build output the same way
  `leaflet.css` already is; nothing new is deployed, versioned, or feature-flagged.
- No rollback plan needed beyond a standard revert — this is a two-dependency, two-file content
  change with full test coverage proving both the positive and negative cases (per qa-report.md).
- Health checks unaffected — no server, no endpoint, no process changed.
- No feature flag needed — this is not a user-facing behavioral change requiring staged rollout;
  the rendered typefaces are visually identical to what Google Fonts served, only the origin
  changed.

## Manual Steps Required

None.

## Known Issues / Follow-ups

- `scripts/ui-suite-in-chrome.sh` and its `--ubi9` variant are not wired into any GitHub Actions
  workflow in this repository today. That gap predates this feature and is out of its scope.paths —
  flagging it here only so it isn't silently assumed to have been fixed by this chore.
- The pre-existing `tests/reagent-map-pan-zoom.spec.ts` timing/precision flake under `--ubi9`
  (documented in qa-report.md's Known Gaps) is unrelated to font sourcing and outside this feature's
  file scope; it remains open for that spec's owner.
```