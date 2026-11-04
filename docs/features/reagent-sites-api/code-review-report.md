---
feature: "reagent-sites-api"
bounded_context: ""
domain_terms: ["Reagent Site", "Aurelis", "Adventurer"]
files_touched: []
issue_refs: ["reagent-sites-api"]
linked_adrs: []
linked_kis: []
---

# Code Review: reagent-sites-api

## Overall Status

**APPROVED**

## Design Narrative

The implementation adds one unauthenticated, read-only endpoint (GET /api/reagent-sites) via a thin readReagentSites() helper that mirrors the existing readItems()/GET /api/items pattern exactly, in the same flat Express file with no new abstraction. The only structural change is a minimal testability fix — exporting the app and guarding app.listen under require.main === module — which is additive and preserves npm start's existing runtime behavior. Behaviour change: yes (adds a new public endpoint; all prior behavior is unchanged).

## Design Score

- **Clarity** 5/5
- **Cohesion** 5/5
- **Coupling** 5/5
- **Craft** 5/5

## Security Surface

- New endpoint is intentionally unauthenticated (no requireAuth) — verified AC requires identical response with/without a session; confirmed test asserts garbage Authorization header returns the same 200 payload.
- No outbound network call introduced — handler only reads a local static JSON file via the existing readJSON() helper.

## Performance Surface

- readReagentSites() re-reads and re-parses the JSON file from disk on every request, same as the pre-existing readItems()/readUsers()/readOrders() pattern — not a new smell, consistent with this file's established (unbounded-cache) convention at its current scale.

## Test Design Review

- 'returns all eight seeded sites with unchanged positions, no session' — YES, would fail: deep-equals response against the spec's seed.reagentSites; verified a mutated lat/lng breaks it (confirmed 4/4 passing on current code, and developer-reported mutation test in implementation-notes.md).
- 'returns the identical array with a garbage Authorization header' — YES, would fail: if requireAuth were mistakenly added to the route, the garbage-token call would return 401 instead of 200, failing assert.equal(status, 200).
- 'every site has an id' — YES, would fail: iterates the data file directly and asserts site.id is truthy for each of the 8 entries; a missing id fails immediately.

## Verification of Developer Self-Review

- Checklist claims matched the code on direct inspection: readReagentSites() and the GET handler are both single-line (complexity 1), no new class/abstraction was introduced, and the file's existing flat-object convention (items/users/orders) was followed for the new endpoint.
- Complexity, duplication, and 'fewest elements' claims verified directly against the diff — no discrepancies found.

## Feedback for the Developer

None
