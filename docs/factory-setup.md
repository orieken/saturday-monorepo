# Running the loom factory in this repo

The work orders in `intake/` are drafts. Nothing runs until a human promotes one. This page covers
the one-time setup, the order to promote them in, and what to watch.

## One-time setup (a human does this; agents cannot write under `.claude/`)

1. **Update the framework.** Done: v3.7.0, PR #1 (2026-10-29). For later updates, see
   `docs/updating-loom.md`.
2. **Ignore local and runtime files.** Add `.loom/` (factory worktrees and index) and
   `.claude/settings.local.json` to `.gitignore`, so the factory's own runs do not dirty the tree.
3. **Write `.claude/delivery-policy.yaml`**, starting with no pre-ship checks. Add each check after
   the order that creates it merges, as each order's Infrastructure section says:
   ```yaml
   maxContractRetries: 3
   preShipChecks: []
   # after ui-suite-real-chrome:      - {name: ui-suite, command: scripts/ui-suite-in-chrome.sh}
   # after golden-lint-and-exemplars: - {name: lint, command: pnpm --filter @orieken/ye-olde-magic-shop lint}
   ```
4. **Run one order first.** This is the factory's first run on a TypeScript repo. Promote
   `ui-suite-real-chrome` alone, run `loom factory run`, and read the whole run before promoting
   more. Note in `ai-assistant-dot-files/docs/prompts/factories/05-pilot-decisions.md` what went
   wrong.

## Order of play

```
ui-suite-real-chrome ─┬─► leaflet-local-tiles ───────────────┐
                      ├─► leaflet-package-tests ─► leaflet-package-gestures ─┴─► leaflet-map-spec-gestures
                      └─► saturday-playwright-cesium (also needs leaflet-package-tests)
                                    └─► map-heatmap-source (also needs leaflet-package-gestures)

golden-workspace-packages ─► golden-site-file-names ─┬─► golden-test-ids ─┐
                                                     └─► golden-flows ────┴─► golden-steps-through-site ─► golden-lint-and-exemplars
```

`depends-on` enforces this: promoting everything at once is safe, and the supervisor claims each
order once its dependencies have shipped. Two chains can run in parallel, one worktree each. Promote
in small batches anyway: every gate is yours, and a batch of thirteen is a long inbox.

Not here yet: the reagent map view and its tests (pilot orders D1, E1 and F1). They are written from
the reagent-map spec, which is being done by hand first.

## Watching it

```bash
loom factory run                       # supervise; Ctrl-C checkpoints the running stages
loom factory inbox                     # what waits on you
loom factory approve <id> <gate>
loom factory stop <id>                 # stop spending now
```

Budgets are $20–$40 an order. The six loom orders delivered so far cost $5–$14 each, but those ran
in Go in the repo the factory was built in. Expect the first TypeScript orders to cost more.
