# Documentation Report: reagent-sites-api

```yaml
feature: "reagent-sites-api"
bounded_context: "Expedition Planning — Reagent Site catalog (apps/mock-api backend)"
domain_terms: ["Reagent Site", "Aurelis", "Adventurer"]
files_touched: []
issue_refs: ["reagent-sites-api"]
linked_adrs: []
linked_kis: []
```

## Status: PASS

## Summary

The mock API gains one new endpoint, `GET /api/reagent-sites`. It returns the eight seeded
reagent sites of the invented planet Aurelis. No request body, query parameter, or session is
required. `apps/mock-api` has no README, CHANGELOG, or OpenAPI document today, so this report and
the inline code comments are the full documentation record for this slice. The permanent record
persists to `docs/features/reagent-sites-api/` after the pipeline completes.

## API Documentation

No OpenAPI or Swagger document exists for `apps/mock-api` to update. The endpoint is documented
here, matching the style of the existing `/api/items` entry in `server.js`'s startup log.

| Method | Path | Auth | Response |
|---|---|---|---|
| GET | `/api/reagent-sites` | None — public by design | `200`, JSON array of 8 `{ id, name, lat, lng }` objects |

- The response is identical whether the request carries no `Authorization` header, a valid
  session token, or a garbage token. This is a deliberate contract decision (spec AC1), not an
  oversight. A future change MUST NOT add `requireAuth` to this route.
- `id` is a kebab-case slug of `name` (for example, `moonwell-spring`). `lat` and `lng` are
  decimal degrees, returned unchanged from `apps/mock-api/data/reagent-sites.json`.
- The endpoint takes no path or query parameters.
- On a missing or malformed `reagent-sites.json`, the endpoint returns `200` with `[]`, the same
  degrade-on-error behavior every other `readJSON()`-backed endpoint in this file already has.

## README / CHANGELOG

`apps/mock-api` has no `README.md` and the repository has no root `CHANGELOG.md`. Neither exists
to update, so no entry was added. If the project adopts a changelog later, this feature's entry
is: "Added `GET /api/reagent-sites`, a public endpoint returning the seeded reagent sites of
Aurelis."

## Architecture Decision Records

No ADR was written. The one structural change in this slice — exporting the Express `app` and
guarding `app.listen()` behind `require.main === module` for testability — is already captured in
`architecture-notes.md` and is not a decision that constrains future features; it is the standard
Express testability pattern.

## Code-Level Docs

- `server.test.js` carries `@issue` / `@ac` annotations on its test block, per
  `shared/rules/testing-conventions.md`'s Test Annotation Convention.
- `server.js`'s new code (`REAGENT_SITES_PATH`, `readReagentSites()`, the `GET /api/reagent-sites`
  handler, and the `============ REAGENT SITES ENDPOINTS ============` block) follows the file's
  existing naming and comment conventions exactly, so no additional inline documentation was
  needed beyond what the code itself already makes clear.

## Known Gaps / Follow-ups

- `apps/mock-api` still has no README or CHANGELOG. Out of scope for this order
  (`scope.paths` is `apps/mock-api/**` and `docs/features/reagent-sites-api/**`, and adding either
  file is a repository-wide documentation decision, not specific to this endpoint).
- Two pre-existing observability gaps noted by `sre-engineer.md` (unstructured logging in
  `POST /api/orders`, no OTel in `apps/mock-api`) are unrelated to this feature's `files_touched`
  and are not addressed here.

## Notes for DevOps

None.