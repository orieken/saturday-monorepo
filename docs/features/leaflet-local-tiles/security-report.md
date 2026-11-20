---
feature: "leaflet-local-tiles"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: ["leaflet-local-tiles"]
linked_adrs: []
linked_kis: []
---

# Security Review: leaflet-local-tiles

## Threat Model Summary

This feature replaces a Leaflet OSM tile layer (external network call to tile.openstreetmap.org) with a static, hand-authored, in-repo SVG loaded via L.imageOverlay from the app's own origin. No auth, no user-supplied input, no API endpoints, no secrets, and no new dependencies are introduced. The change reduces attack/disclosure surface by eliminating the view's only outbound network call (which previously leaked the user's IP and request pattern to tile.openstreetmap.org on every pan/zoom). Verified the new SVG asset contains no <script>, event-handler attributes, or external xlink:href/href references — it is inert static markup. A new Playwright test (map.spec.ts) asserts that blocking every off-origin request still allows the view to load and the image overlay to render, enforcing the zero-egress requirement at the test level. No findings were raised; no code changes were made by this review.

## Dependency Audit

- No new dependencies added by this feature (implementation-notes.md confirms 'Dependencies Added: None').
- pnpm audit on the monorepo shows pre-existing findings (8 critical, 61 high, 53 moderate, 8 low) in unrelated packages (vite, body-parser, picomatch, etc.) via transitive paths not touched by this feature — these predate this change and are out of scope for this review.

## STRIDE Analysis

None

### Spoofing

No authentication, session, or identity concept exists in this view; not applicable.

### Tampering

The map image is a static asset committed to the repo (public/sword-coast-map.svg), not user- or runtime-supplied; no path for in-transit or at-rest tampering by an outside party.

### Repudiation

No state-changing or sensitive action occurs in this view; not applicable.

### Information Disclosure

Removing the OSM tile fetch eliminates the only mechanism by which the user's IP/UA and viewing pattern left the app's origin. The SVG was inspected for embedded script, event handlers, or external references and is confirmed inert — it cannot exfiltrate anything when rendered as an <img> via L.imageOverlay.

### Denial of Service

Single static local image with no per-tile fetch loop; no unbounded iteration, no external dependency whose outage or latency could degrade this view. Risk reduced relative to the prior OSM-tile-layer implementation.

### Elevation of Privilege

No roles, permissions, or privilege boundaries exist in this feature; not applicable.

## Findings

None

## Files Modified

None

## Security Checklist

- [x] No hardcoded secrets, tokens, or credentials introduced
- [x] No new external network calls introduced; existing one removed
- [x] New static asset (SVG) inspected for embedded scripts/external references — grep for script|onload|onerror|xlink:href|href="http|foreignObject found nothing
- [x] No new third-party dependencies added
- [x] Dependency vulnerability scan run (pnpm audit) — All findings are pre-existing and unrelated to this feature's changed files
- [x] User-supplied input validated — N/A — feature has no user input

## Notes for QA

- No security-driven test additions are required beyond what the developer already added — the off-origin-blocking Playwright test in map.spec.ts already enforces the zero-egress NFR.

## Notes for Tech Writer

- No security findings to record for this feature; the delivery summary's Security Findings section should read 'None identified'.
