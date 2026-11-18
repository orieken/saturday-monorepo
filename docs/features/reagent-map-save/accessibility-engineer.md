---
feature: "reagent-map-save"
bounded_context: "Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7"
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue"]
issue_refs: ["reagent-map-save"]
linked_adrs: []
linked_kis: []
---

# Accessibility Report: reagent-map-save

## Overall Status

**APPROVED WITH FIXES APPLIED**

## Scope Reviewed

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — the new save-line control, its disabled state, and the message regions this order adds (`save-line-tool`, `save-line-message`).
- `apps/ye-olde-magic-shop/src/api/expeditions.ts` — reviewed for completeness only (no markup); no findings, it carries no UI.

## Semantic HTML Audit

- The save control is a native `<button type="button">`, not a `<div>`/`<span>` with a click handler — correct use of an interactive element, consistent with the existing `start-line-tool` button.
- The sign-in/zero-length/failure message is a `<p>` with `role="alert"`, matching the existing `reagent-map-error` convention for urgent, assertive text.
- No new `div`-soup or non-semantic interactive elements were introduced.

## ARIA & WCAG Compliance

- `save-line-tool` carries `aria-label="Save this expedition"`, matching its own visible text — redundant but harmless, and consistent with the existing `start-line-tool` button's same pattern. No change needed.
- `save-line-message` correctly uses `role="alert"` for the sign-in prompt, the server's zero-length text, and any failure message — all are time-sensitive, worth-interrupting text, matching `reagent-map-error`'s existing precedent.
- The native `:disabled` binding on `save-line-tool` is correct: it removes the button from the tab order and is announced as disabled by assistive technology, rather than using `aria-disabled` plus manual focus suppression. The reason the button is disabled (fewer than two draft points) is already communicated to assistive technology separately via `line-tool-status` (`role="status"`, `aria-live="polite"`), which is read before the button becomes enabled.
- **Finding (fixed): no accessible confirmation on a successful save.** On success, `promoteSavedExpedition()` cleared `saveMessage` to `null` and reset `toolStarted` to `false`, which also removes `line-tool-status` from the DOM (it is `v-if="toolStarted"`). The only feedback was a new polyline drawn inside the Cesium canvas — visible to a sighted adventurer, but the canvas is exposed to assistive technology only as a single static `role="img"` with a fixed `aria-label` that does not change. A screen-reader user clicking Save had **no programmatically determinable confirmation that the save succeeded** — a WCAG 2.2 §4.1.3 (Status Messages, AA) gap. Fixed directly: added a `save-line-status` paragraph (`role="status"`, `aria-live="polite"`, matching the existing polite-announcement pattern already used by `line-tool-status`) that is set to `'Expedition saved.'` in `promoteSavedExpedition()` and cleared alongside `saveMessage` in `startLine()`. `role="status"` (not `alert`) was chosen deliberately — a success confirmation is informational, not urgent, so it should not interrupt the way the error/sign-in message does.

## Keyboard Navigation

- `save-line-tool` is a native button inserted immediately after `start-line-tool` in document order, so it falls naturally into the existing tab sequence with no `tabindex` manipulation needed.
- Enter and Space both activate it via native button semantics; no custom `keyup`/`keydown` handling was added or needed.
- `:focus-visible` styling is inherited from the project's existing button treatment (shared with `start-line-tool`); no new focus style was introduced that would need separate verification.
- One residual UX note, not blocking: when a save succeeds, the button itself is removed from the DOM (`toolStarted` → `false`), so a keyboard user's focus falls back to the document body. This is the same focus-loss pattern the pre-existing `start-line-tool`/`toolStarted` toggle already has elsewhere in this component (e.g., starting a new line also conditionally removes/adds controls), so it is not a regression introduced by this order. Worth a follow-up to move focus to `start-line-tool` after a successful save, but it is not a new defect this order created, and fixing it would touch focus-management behavior beyond this order's scope — escalated as a non-blocking follow-up rather than fixed here.

## Fixes Applied

1. Added a `role="status"` / `aria-live="polite"` confirmation element (`save-line-status`, text "Expedition saved.") so a successful save is announced to assistive technology, not only drawn on the Cesium canvas. New `saveStatus` ref set in `promoteSavedExpedition()` and cleared in `startLine()`.

## Escalations

- **Focus management after a successful save** (see Keyboard Navigation above): focus is dropped to `<body>` when the save button disappears. Pre-existing pattern in this component, not introduced by this order — flagged for a future pass rather than fixed here, since remediation would mean deciding where focus should land across the whole tool-start/tool-end lifecycle, not just the save path.

## Verification

- No automated a11y tooling was run, per this agent's contract — findings are from static analysis of the template and script changes in `ReagentMapView.vue` against the NFR in `analysis.md` ("The save control and the sign-in/zero-length message MUST be reachable by keyboard and announced to a screen reader") and WCAG 2.2 AA.
- The fix adds no new interactive element, so no new keyboard-reachability or focus-order check was needed beyond what's covered above.
