---
feature: "golden-test-ids"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: ["golden-test-ids"]
linked_adrs: []
linked_kis: []
---

# Security Review: golden-test-ids

## Threat Model Summary

This change replaces brittle text/structure-based Playwright locators in CartPage.ts and HomePage.ts with data-testid-based locators, and renames one data-testid attribute in CartView.vue so CartPage's completeTransactionBtn locator resolves. No new interfaces, endpoints, auth paths, user input handling, secrets, or external calls are introduced. The spec itself declares trust boundaries as none (no auth, no user input, no external calls), and the diff confirms this: all changes are to test-locator strings and one static HTML attribute value. STRIDE categories were each considered and found not applicable to this change's attack surface; no findings were raised and no code modifications were required from this review.

## Dependency Audit

- No new dependencies were introduced by this change — dependency audit not applicable.

## STRIDE Analysis

None

### Spoofing

data-testid attributes are static, non-sensitive DOM markers with no bearing on identity or session verification; renaming one (checkout-btn -> complete-transaction-btn) and retargeting a Playwright locator creates no spoofing surface.

### Tampering

Changes are confined to test-only locator strings (lib/site/*.ts) and a presentational attribute in CartView.vue; no data is read, parsed, or mutated as a result of these values, so there is no tamperable data path introduced.

### Repudiation

No new user actions, audit events, or logging paths are added or altered; the checkout button's click handler (@click="checkout") and its audit trail, if any, are unchanged by this diff.

### Information Disclosure

data-testid values are already visible in rendered HTML/DOM to any client regardless of this change; no new data, selectors, or internal paths are exposed that weren't already present as plain CSS/text selectors before.

### Denial of Service

No loops, network calls, or resource-bound operations are touched; locator resolution is a standard Playwright DOM query with no unbounded iteration introduced.

### Elevation of Privilege

No authorization logic, role checks, or permission boundaries exist in or near the changed files; this is purely a test-infrastructure and template-markup change.

## Findings

None

## Files Modified

None

## Security Checklist

- [x] No hardcoded secrets, tokens, or credentials introduced
- [x] No new external calls or trust boundary crossings introduced
- [x] No new dependencies requiring vulnerability scan
- [x] No auth/session/authorization logic touched
- [x] No user-supplied input handling introduced

## Notes for QA

- Confirm CartPage's completeTransactionBtn locator ([data-testid="complete-transaction-btn"]) now resolves against CartView.vue's renamed attribute and the existing checkout scenario still passes end-to-end.
- Confirm HomePage.clickFirstProductDetails() still navigates correctly now that it targets [data-testid="item-name"] instead of the [data-testid="item-card"] a selector — verify item-name is rendered as a clickable/anchored element in the markup.
- No security-driven test cases are required for this change; it is out of scope for STRIDE-based test additions.

## Notes for Tech Writer

- No documentation changes required for this change — it is an internal test-locator refactor with no user-facing or API-facing behavior change.
