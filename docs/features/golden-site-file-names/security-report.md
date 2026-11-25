---
feature: "golden-site-file-names"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Security Review: golden-site-file-names

## Threat Model Summary

This change renames Saturday site/page model files to kebab-case name.type.ts (e.g. MagicShop.ts -> magic-shop.site.ts) and updates every import to match. Class names, public interfaces, and all test/feature assertions are unchanged — confirmed by reading every modified diff hunk: each is either a pure git rename (0 content lines changed) or a one-to-six-line import-path update only. No new dependency was introduced (no package.json changed), no new endpoint, no auth/session code, no user-input handling, and no code crosses a trust boundary. The feature spec itself declares 'Trust Boundaries: None — no auth, no user input, no external calls,' which matches what the diff shows. STRIDE analysis below is correspondingly brief: there is no new attack surface to model.

## Dependency Audit

- No new or modified dependencies — package.json files are unchanged in this diff. Dependency/CVE audit (pnpm audit) not applicable for a file-rename-only change.

## STRIDE Analysis

None

### Spoofing

No identity, session, or auth code touched; MagicShop/LoginPage/AccountPage class logic is byte-identical, only their file paths and intra-package import paths changed.

### Tampering

No data handling logic changed; only import statements and the module filenames they reference were rewritten.

### Repudiation

No logging, audit-trail, or OTel instrumentation code was added or removed.

### Information Disclosure

No new egress paths, no new file reads, no change to what any component can send or receive — this is a rename, not a data-flow change.

### Denial of Service

No loops, network calls, or resource-handling code changed; nothing new that could be exhausted.

### Elevation of Privilege

No authorization checks exist in this feature's scope before or after the change.

## Findings

None

## Files Modified

None

## Security Checklist

- [x] Verified every renamed/modified file's diff hunk is import-path-only or a pure rename with no content change
- [x] Verified no new dependency was introduced (package.json unchanged)
- [x] Verified no auth, session, input-handling, or external-call code exists in the touched files
- [x] Confirmed class names and public interfaces are unchanged per developer's interfaceDesign notes

## Notes for QA

- No security-driven test additions needed — this is a rename with zero new behavior and zero new attack surface.
- Standard regression focus applies: confirm every existing Cucumber/Playwright scenario still passes unchanged (per spec AC3) since that is the only risk this change carries.

## Notes for Tech Writer

None
