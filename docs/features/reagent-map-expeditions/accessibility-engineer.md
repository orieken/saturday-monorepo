I already fixed the accessibility gap I found directly in the source file. Per this invocation's instructions, here is the complete artifact on stdout:

```markdown
---
feature: "reagent-map-expeditions"
bounded_context: "Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7"
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue"]
issue_refs: ["reagent-map-expeditions"]
linked_adrs: []
linked_kis: []
---

# Accessibility Report: reagent-map-expeditions

## Overall Status

**APPROVED**

## Scope Reviewed

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — the site-list button markup, the new expedition polyline entities, and the globe's `role="img"` accessible name.
- `apps/ye-olde-magic-shop/src/api/expeditions.ts` — no markup; out of scope for this review.

## Semantic HTML Review

- Each `<li>` in `reagent-site-list` now wraps its text in a native `<button type="button">` rather than attaching `@click` to the `<li>` or a `<span>`. This is the correct pattern — no div/li-soup click handler was introduced.
- `type="button"` is set explicitly, so the control cannot be mistaken for a form submit button.
- No new non-semantic container was introduced for the expedition polylines; they are added as Cesium scene entities, not DOM elements, so there is no DOM markup to review for them beyond the globe's existing `role="img"` container.

## Interactive Elements Review

- `flyToSite(site)` is bound via `@click` on a real `<button>`, not a `<div>` or `<li>` — satisfies the "no onClick on non-interactive elements" rule.
- The button carries no custom `tabindex`, `role`, or keyboard handler of its own — it relies entirely on the browser's native button semantics, which is correct and simpler than a hand-rolled keyboard handler would have been.

## Accessibility (WCAG) Review

- Each fly-to-site button has `aria-label="Fly to {{ site.name }}"`, giving it an accessible name that announces the action, not just the destination name — meets analysis.md's accessibility NFR exactly.
- **Finding (fixed)**: the globe's `role="img"` container's `aria-label` described only reagent sites ("3D globe of Aurelis with every reagent site marked...") and did not mention the newly drawn expedition polylines. Expeditions are rendered purely as visual Cesium scene geometry with no DOM text equivalent anywhere else in the view, so a screen-reader user had no way to learn expeditions exist at all — a WCAG 1.1.1 (Non-text Content) gap newly introduced by this order. Fixed directly: the label now reads "3D globe of Aurelis with every reagent site marked and every known expedition drawn as a line between its two sites; the full list of sites is in the adjacent list". This is a text-only attribute change; it does not alter any rendered pixel, so it does not affect the `globe-default-view.png` baseline or `reagent-map-a11y-check.spec.ts`'s existing `toHaveAttribute('aria-label', /.+/)` assertion.
- No `<input>` elements were added by this order — the existing-input label-association rule does not apply here.

## Keyboard Navigation Review

- The fly-to-site buttons are native `<button>` elements with no `outline-none` or other focus-suppressing class in their Tailwind utility list (`block w-full cursor-pointer text-left`), so the browser's default `:focus-visible` outline remains visible on Tab. Confirmed no project-wide CSS resets `outline` for buttons.
- Each button activates on both Enter and Space via native button behavior — no custom key handling was written that could get this wrong.
- Tab order follows document order (list items in their existing order) — no `tabindex` override was introduced that could disrupt it.

## Fixes Applied

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue`: updated the globe's `aria-label` to also describe the newly drawn expeditions, closing the non-text-content gap above.

## Escalations

None — no component required a fundamental redesign to be accessible.
```