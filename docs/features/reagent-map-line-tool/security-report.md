---
feature: "reagent-map-line-tool"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Aurelis"]
files_touched: []
issue_refs: ["reagent-map-line-tool"]
linked_adrs: []
linked_kis: []
---

# Security Review: reagent-map-line-tool

## Threat Model Summary

This feature adds a two-click line-drawing tool to the existing Cesium globe view: a button arms the tool, a LEFT_CLICK handler resolves screen coordinates to a geographic point via pickEllipsoid + EllipsoidalOccluder visibility check, and up to two points are rendered as draft entities tagged with empty id/name placeholders. The diff touches exactly one file (ReagentMapView.vue), adds no new dependency, no new network call, no new persisted or transmitted data, and no auth/session/permission logic. Saving the drawn line is explicitly out of scope (reagent-map-save), so there is no trust boundary crossing in this change beyond the existing read-only fetchReagentSites/fetchExpeditions calls, which are unmodified. All STRIDE categories were assessed and none produced an actionable finding; this is a genuine no-security-surface feature, not an unreviewed one.

## Dependency Audit

- No package.json or lockfile changes in this diff (git diff --stat confirms only ReagentMapView.vue changed) — no new dependency introduced, so no CVE/audit scan was needed for this stage.

## STRIDE Analysis

None

### Spoofing

No auth, session, or identity logic touched. The existing Authorization header in fetchExpeditions (unchanged in this diff) is unaffected by the line tool.

### Tampering

Draft points/line exist only as component-local Vue refs and Cesium entities; nothing is sent to an API, written to storage, or shared between users. There is no data for an attacker to tamper with because none crosses a boundary yet (saving is out of scope for this order).

### Repudiation

Ephemeral client-side UI state with no persistence or audit requirement in the spec; not applicable until reagent-map-save introduces a mutation to audit.

### Information Disclosure

Checked both read and egress direction. Read: draft entities are tagged via tagEntity() with id:'' and name:'' — no site/user data is attached to drafts, so nothing sensitive is newly exposed via getEntities(). Egress: no new fetch/XHR/image/script URL is constructed from user click coordinates or any other input; pickLandPoint only calls Cesium math/geometry APIs (pickEllipsoid, EllipsoidalOccluder, Cartographic.fromCartesian) which stay in-process. No new outbound destination exists for data to leave through.

### Denial of Service

handleMapClick is bounded by the toolStarted flag and draftPoints.length >= MAX_LINE_POINTS (2) guard, so repeated clicks after two points are placed are no-ops — no unbounded entity accumulation. No loop over user-controlled collections, no new network call to exhaust, no missing timeout introduced (the file's existing fetches already have timeouts, unchanged here).

### Elevation of Privilege

No role, permission, or server-side authorization logic exists in this diff to elevate past. The tool draws locally for whoever is already viewing the page; it grants no new capability against the API.

## Findings

None

## Files Modified

None

## Security Checklist

- [x] No hardcoded secrets, API keys, or tokens introduced
- [x] No new network request destination added (NFR: no request leaves app/mock-api origins)
- [x] No PII or sensitive data added to entity metadata, logs, or telemetry — Draft entities carry id:'' and name:'' — empty placeholders, not real data
- [x] No new dependency requiring a CVE scan
- [x] Input from clicks (screen coordinates) only flows through Cesium's own geometry APIs, never into a DOM sink, template string, or request
- [x] DoS: click handling is bounded by MAX_LINE_POINTS guard
- [x] No auth/session/authorization code touched — N/A for this feature

## Notes for QA

- No security-driven test requirements beyond the functional acceptance criteria — there are no new auth paths, API calls, or data boundaries to add negative security test cases for.
- Worth confirming at the functional level (already covered by AC9/AC10) that a click off the globe or on the far side truly places nothing — this is also the DoS/robustness boundary, not just a UX one.

## Notes for Tech Writer

None
