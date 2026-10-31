# Performance Report: ui-suite-real-chrome

**Stage**: Shift-left performance & reliability review
**Reviewed artifacts**: `analysis.md`, `architecture-notes.md`
**Gate verdict**: ✅ PASS — no critical blockers. Two recommendations for the developer below.

## Scope Note

This feature has no API endpoints, no database access, and no request-handling loop — it is a
Playwright config addition, a new spec file, a Bash orchestration script, and a `.dockerignore`
fix. The four governing principles below are evaluated against that surface: the orchestration
script's Docker lifecycle and health-check loop are the only places where network calls or
repeated operations occur.

## Idempotency on Mutations

No API or data mutation exists in this feature — N/A for the standard "retried POST charges
twice" case.

One operational idempotency risk in `scripts/ui-suite-in-chrome.sh` is worth flagging before the
developer writes it: if the script is run twice in a row (exactly the "pre-ship check" use case
named in the spec's Infrastructure/Deploy Notes), and the `docker build`/`docker run` invocations
use a fixed container name, the second run fails with "name already in use" instead of cleanly
replacing the first run's container. Recommendation: either run the container unnamed/ephemeral
(`docker run --rm` without `--name`), or have `main` remove any stale container from a prior run
before starting a new one. Not a blocker — the 034adbd draft should be checked for this before
the developer treats it as finished, since a factory pre-ship check that can't be re-run cleanly
defeats its own purpose.

## Strict Timeouts

The architecture already sets the one timeout that matters most here: the 180s bound on the
app/mock-API health-check loop (architecture-notes.md, Stability Design), with a distinct exit
code (4) on timeout. That satisfies the outer bound.

**Finding — the inner `curl` calls need their own explicit timeout, not just the outer counter.**
A bounded `until curl ...; do sleep N; done` loop is only as bounded as its individual iterations.
If a single `curl` call to `localhost:8000` or `localhost:8001` hangs (e.g., the process accepted
the TCP connection but never responds), that one call can block past its intended slice of the
180s budget — in the worst case, if the app wedges mid-handshake, a single hung `curl` could
consume most or all of the remaining budget in one iteration instead of failing fast and retrying.
Require `curl --max-time 2 --connect-timeout 2` (or equivalent) on every health-check call inside
the loop, so each iteration has its own short, explicit timeout and the outer 180s bound is backed
by fast-failing iterations rather than relying on one call to misbehave quickly. This is
`architecture-guardrails.md` #5's "every network call MUST have an explicit timeout" applied one
level deeper than the architecture notes already went — loopback calls are not exempt just because
they're local.

No other network calls exist in this feature's scope (the `docker build` pull of base-image layers
is inside the already-verified, out-of-scope Containerfile).

## Prevent N+1 Queries

No immediate risks identified. There is no loop in this feature's scope that performs a database
call or a per-item network request — `webgl.spec.ts` makes one page load and one
`getContext('webgl2')` check; the orchestration script's only loop is the bounded health-check
poll addressed above, which is a readiness gate on a process the script itself started, not a
per-record fetch pattern.

## Caching Strategies

Already correctly identified by the architect: the named pnpm store volume for
`pnpm install --frozen-lockfile` is exactly the right caching strategy for this script's hot
path — a pre-ship check that re-runs frequently should not re-download the dependency tree on
every invocation. No change needed; confirming this decision stands rather than re-deriving it.

No other hot path exists in this feature (the Docker image build itself is out of scope and
already has its own verified build-time characteristics per `docs/runner-image.md`).

## Blocking / Gate Decision

No critical N+1, idempotency, or missing-timeout defect blocks the developer from proceeding.
The one finding that should be addressed before this is marked done is the missing per-call
`curl` timeout inside the health-check loop — small, mechanical, and consistent with the explicit-
timeout discipline the architecture notes already applied to the loop's outer bound. The container
re-run idempotency point is a recommendation, not a gate.

## Developer Handoff Notes

- Add `--max-time 2 --connect-timeout 2` (or project-appropriate values) to every `curl` call
  inside the health-check loop in `scripts/ui-suite-in-chrome.sh`, so each polling iteration has
  its own explicit timeout independent of the outer 180s counter.
- Verify the 034adbd draft's container lifecycle allows a clean second run (no stale named
  container collision) before treating the script as finished.