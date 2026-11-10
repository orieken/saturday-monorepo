---
feature: "local-fonts"
bounded_context: "ye-olde-magic-shop frontend — typography and font-serving"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/src/main.ts", "apps/ye-olde-magic-shop/src/styles.css", "apps/ye-olde-magic-shop/package.json", "apps/ye-olde-magic-shop/tests/local-fonts.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts"]
issue_refs: ["local-fonts"]
linked_adrs: []
linked_kis: []
---

# Accessibility Report: local-fonts

## Overall Status

**APPROVED — NO UI SURFACE**

## Scope Reviewed

This chore swaps the cross-origin `@import` of Cinzel and Inter in
`apps/ye-olde-magic-shop/src/styles.css` for four self-hosted `@fontsource` CSS imports added to
`apps/ye-olde-magic-shop/src/main.ts`, pins the two new dependencies in `package.json`, adds a
Playwright test (`tests/local-fonts.spec.ts`) that checks font weight through
`document.fonts.load`, and removes a host exception from `tests/reagent-map-globe.spec.ts`.

I read `analysis.md` and `implementation-notes.md`, then read the actual diff for every touched
file (`git diff --cached`) to confirm no file outside that list changed.

## Semantic HTML Review

Not applicable. No `.vue`, `.html`, `.tsx`, or `.jsx` file is touched by this change. No new
markup, element, or template is introduced or modified.

## Interactive Elements Review

Not applicable. No event listener, `onClick`, `keyup` handler, `<button>`, or `<a>` is added,
removed, or modified.

## ARIA and Labeling Review

Not applicable. No `aria-*` attribute or form control is touched.

## Keyboard Navigation Review

Not applicable. No interactive element or focus order changes. The font swap does not alter tab
order, focus rings, or `:focus-visible` styling — `styles.css`'s only change is deleting the
cross-origin `@import` line; no selector touching focus state is modified.

## Visual / Typography Impact

The fonts rendered (Cinzel at 400/700, Inter at 400/600) are unchanged — only their origin moves
from `fonts.googleapis.com` to the app's own bundle via `@fontsource`. Font-family, weight, and
fallback stack are identical to before, so there is no color-contrast, text-sizing, or
reflow impact to assess.

## Violations Found

None.

## Fixes Applied

None — no accessibility-relevant code exists in this diff to fix.

## Escalations

None.

---
*Reviewed by accessibility-engineer. This feature has no frontend markup surface — the pipeline's
UI/a11y gate is satisfied by confirming that fact against the actual diff, not by inspecting
components that don't exist in it.*
