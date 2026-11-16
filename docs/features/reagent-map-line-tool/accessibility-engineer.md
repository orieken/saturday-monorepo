# Accessibility Report: reagent-map-line-tool

## Scope

Reviewed `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` as modified by this order: the new
"Draw a new expedition" line-tool control, its click-driven point-placement flow on the Cesium
globe, and the draft-point/draft-line entities it adds. No other file in `files_touched` carries
markup.

## Semantic HTML Audit

- The new control is a real `<button type="button">` with visible text content, not a `<div>` or
  `<span>` with a click handler — correct per the Interactive Elements rule.
- It sits in its own `<div class="mb-4">`, a sibling of the existing header, error banner, globe,
  and site list — not nested inside any of them, so it does not change the DOM structure those
  elements already had.
- No other new markup was added by the developer beyond this one button.

## ARIA & Labeling

- `aria-label="Draw a new expedition"` duplicates the button's own visible text. This is
  compliant (WCAG 2.5.3 Label in Name is satisfied trivially when they match) but adds nothing;
  harmless, left as-is to match this file's established button convention (every `fly-to-site-*`
  button also carries both).
- **Fixed**: the line tool had no way to tell a sighted user what state it was in — "armed and
  waiting for point 1," "one point placed, waiting for point 2," or "line complete" were only
  observable by inspecting Cesium entities (`getEntities()`), which is a test-only path, not a
  user-facing one. Added a `role="status" aria-live="polite"` paragraph
  (`data-testid="line-tool-status"`) beside the button, driven by a new `lineToolStatus` computed
  property, so the state change is announced to assistive technology and visible to sighted users
  without requiring a redesign. It renders only once `toolStarted` is true, so it changes nothing
  about the page's default (tool-not-started) appearance.
  - Placed in the same sibling `<div>` as the button — outside the four regions
    `reagent-map-visual.spec.ts` screenshots (`header`, `reagent-site-list`, `reagent-map-error`,
    `#cesiumContainer`) — so it cannot affect any `@visual` baseline.

## Keyboard Navigation

- The "Draw a new expedition" button itself is a native `<button>`: reachable by Tab, activates
  on Enter/Space, and keeps the browser's default `:focus-visible` outline (no CSS in this file
  suppresses it). No violation here.
- **Escalation — not fixed, needs a human decision**: placing a point is only possible by a
  pointer click on `#cesiumContainer` (`viewer.screenSpaceEventHandler.setInputAction(...,
  LEFT_CLICK)`). There is no keyboard path to place a point, complete a line, or reach either of
  the two points once the tool is armed. This is not a markup defect fixable with an ARIA
  attribute or a semantic-element swap — a 3D globe's click-to-geographic-coordinate picking has
  no native keyboard equivalent, and building one (e.g., an arrow-key-driven crosshair, or a
  coordinate-entry fallback) is a real interaction-design decision, not a one-line fix. The
  developer's own analysis already flagged this and deferred it ("Placing a point is a
  mouse-only gesture with no keyboard equivalent... flagged for accessibility-engineer, not fixed
  here"), and this order's acceptance criteria describe only the mouse-driven behavior. Recorded
  here as a known gap for a human to decide whether it blocks ship or is tracked as follow-up
  work (most plausibly scoped into `reagent-map-save`, once the line tool's final interaction
  shape is settled).
- Related: `#cesiumContainer` carries `role="img"`, which tells assistive technology this element
  is a static, non-interactive image. That was already slightly inaccurate before this order (the
  globe already moved via `flyToSite`, triggered indirectly from the site list), but it is more
  inaccurate now that the same element also owns a live `LEFT_CLICK` handler that mutates state.
  Not changed here: swapping it for something interactive (e.g. `role="application"` plus real
  keyboard handling) is the same fundamental-redesign work as the keyboard-equivalent gap above,
  and doing it without the keyboard path behind it would make AT behavior worse, not better.
  Flagging alongside the escalation above rather than fixing in isolation.

## Fixes Applied

- Added a `role="status" aria-live="polite"` line-tool status message (`lineToolStatus` computed
  + `data-testid="line-tool-status"` paragraph) so the tool's armed/in-progress/complete state is
  announced and visible, not only inspectable through `getEntities()`.

## Escalated (Not Fixed — Needs Human/Design Decision)

- The line tool's point-placement gesture has no keyboard equivalent. Fixing this requires
  designing an actual alternate interaction (not a markup change) and is out of this order's
  acceptance criteria. Carried forward from the developer's own analysis; confirmed here rather
  than fixed.

## Verdict

No blocking semantic-HTML or markup-level violations. One accessibility gap (keyboard equivalent
for point placement) is real but requires a design decision beyond this review's scope, and is
explicitly out of this order's stated acceptance criteria — escalated, not fixed. One in-scope fix
applied (live status region) with no risk to existing `@visual` baselines.