# Documentation Report: expeditions-api

```yaml
feature: "expeditions-api"
bounded_context: "Expedition Planning — Expedition storage (apps/mock-api backend)"
domain_terms: ["Adventurer", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: []
issue_refs: ["expeditions-api", "reagent-map-spec"]
linked_adrs: []
linked_kis: []
```

## Status: PASS

## Summary

The mock API gains two new endpoints: `GET /api/expeditions` and `POST /api/expeditions`. An
adventurer saves an expedition (a line between two points on Aurelis) and finds it again later.
The shop's three seeded expeditions are visible to every caller. No adventurer ever sees another
adventurer's expedition. `apps/mock-api` has no README, CHANGELOG, or OpenAPI document today, so
this report and the inline code comments are the full documentation record for this slice, the
same convention the sibling `reagent-sites-api` delivery established. The permanent record
persists to `docs/features/expeditions-api/` after the pipeline completes.

## API Documentation

No OpenAPI or Swagger document exists for `apps/mock-api` to update. Both endpoints are documented
here, matching the style of the existing entries in `server.js`'s startup log.

| Method | Path | Auth | Response |
|---|---|---|---|
| GET | `/api/expeditions` | Optional — reads a session if present, never requires one | `200`, JSON array of expeditions |
| POST | `/api/expeditions` | Required — `Bearer` session token | `200` with the stored expedition, or `400`/`401`/`422` |

### GET /api/expeditions

- With no `Authorization` header, or one that does not match an active session, the response is
  exactly the shop's three seeded expeditions (`owner: "shop"`).
- With a valid session, the response is the seeded expeditions plus the caller's own. It never
  includes another adventurer's expedition — this is the feature's trust boundary.
- An invalid or expired token behaves exactly like no session: seeded expeditions only, never a
  `401`. This is a deliberate contract decision, not an oversight — `GET` uses an optional-auth
  lookup, unlike `POST`.
- Each item in the array has the shape `{ id, owner, from: {lat, lng}, to: {lat, lng}, savedAt }`.
  Seed records also carry a human-readable `name`; the API contract does not require it.
- The endpoint takes no path or query parameters, and the response is not paginated. This matches
  the existing, pre-existing behavior of `GET /api/items` and `GET /api/reagent-sites` and is a
  known sitewide gap, not something this feature introduces or fixes (see Known Gaps below).

### POST /api/expeditions

- Requires a valid session. With no session, or an invalid token, the response is `401` and
  nothing is stored.
- Request body: `{ from: { lat, lng }, to: { lat, lng } }`. Both points MUST have finite
  `lat`/`lng` values, with `lat` between -90 and 90 and `lng` between -180 and 180. A missing
  field, a non-numeric value, `NaN`, `Infinity`, or an out-of-range coordinate is rejected with
  `400` and `{ error: "from and to are required as { lat, lng }" }`.
- The two points MUST be at least 1 km apart on Aurelis (a sphere of radius 6,000 km), measured by
  great-circle (haversine) distance. Two points closer than 1 km are rejected with `422` and
  `{ error: "Zero-length line: from and to are the same place" }`, naming the zero-length line.
  Exactly 1.0 km apart is accepted — the rule is "closer than," not "at or closer than."
- The stored and returned `owner` is always the session's adventurer (`req.userId`). Any
  `owner` or `userId` field present in the request body is ignored entirely — ownership can never
  be forged from the client.
- On success, the response is the full stored expedition, including a generated `id` and a
  `savedAt` ISO-8601 timestamp, and it is immediately visible to that adventurer's own subsequent
  `GET /api/expeditions` call.
- The distance check is pure local arithmetic. No request this endpoint makes leaves the app's and
  mock API's own origins — there is no external geocoding or mapping call.

## README / CHANGELOG

`apps/mock-api` has no `README.md` and the repository has no root `CHANGELOG.md`. Neither exists
to update, so no entry was added. If the project adopts a changelog later, this feature's entry
is: "Added `GET /api/expeditions` and `POST /api/expeditions` — adventurers save and retrieve
expeditions on Aurelis; the shop's seeded expeditions are visible to everyone."

## Architecture Decision Records

No ADR was written. This slice adds no new package, base class, or layer boundary — it follows
the existing `readJSON`/`writeJSON`, `requireAuth`, and owner-filtering patterns already
established by `GET`/`POST /api/orders`, with one addition: `resolveSessionUserId`, an optional-
auth session lookup for the one route (`GET /api/expeditions`) that personalizes its response
without requiring a session. That decision is captured in `architecture-notes.md` and is not a
constraint future features must honor; it is the standard pattern for an optional-auth read.

## Code-Level Docs

- `server.test.js` carries `@issue` / `@ac` annotations on every new test case, per
  `shared/rules/testing-conventions.md`'s Test Annotation Convention, covering AC3, AC4 (and its
  mirror for Borin), AC11, AC12, AC13 (including the exact-1km boundary), the invalid-token case,
  malformed-body cases, and the two security-hardening fixes (non-finite and out-of-range
  coordinates).
- `server.js`'s new code (`EXPEDITIONS_PATH`, `readExpeditions()`/`writeExpeditions()`,
  `resolveSessionUserId()`, `geodesicDistanceKm()`, `isValidPoint()`, and the
  `============ EXPEDITIONS ENDPOINTS ============` block) follows the file's existing naming and
  comment conventions exactly. Two comments explain non-obvious constraints: why
  `resolveSessionUserId` must never reject a request (unlike `requireAuth`), and why `isValidPoint`
  uses `Number.isFinite` instead of `typeof` (a security fix — see below). No additional inline
  documentation was needed beyond what the code already makes clear.

## Security Fixes Reflected in Behavior Above

`security-reviewer` found and fixed two validation gaps during this delivery; the API
documentation above already reflects the corrected, final behavior rather than the pre-fix
behavior:

- `isValidPoint` now rejects `NaN`/`Infinity` (via `Number.isFinite`, not `typeof`) and
  out-of-range latitude/longitude. Before the fix, a crafted body could defeat the 1 km
  zero-length-line check and persist corrupted coordinates.
- The `GET` owner filter now requires `userId !== null` before matching an expedition's `owner`,
  closing a latent null-coincidence path on the feature's stated trust boundary.

## Known Gaps / Follow-ups

- `apps/mock-api` still has no README or CHANGELOG. Out of scope for this slice
  (`scope.paths` is `apps/mock-api/**` and `docs/features/expeditions-api/**`, and adding either
  file is a repository-wide documentation decision, not specific to these two endpoints).
- `GET /api/expeditions` returns an unbounded, unpaginated array — a pre-existing sitewide gap
  shared with `GET /api/items` and `GET /api/reagent-sites`, flagged by `code-reviewer` as a
  follow-up, not addressed here.
- `isValidPoint` (apps/mock-api/server.js:99) measures cyclomatic complexity 7 against the
  framework's `< 7` budget, per `qa-report.md`. Flagged for the next touch to that function to
  extract a separate range-check helper; does not block shipping this slice and needs no test
  change.

## Notes for DevOps

None.