# Accessibility Report: reagent-map-globe

## Overall Status

**APPROVED WITH FIXES APPLIED**

## Scope Reviewed

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` (new view, template + script)
- `apps/ye-olde-magic-shop/src/App.vue` (new nav link only — `<router-link to="/reagent-map">`)

`src/api/reagentSites.ts`, `src/featureFlags.ts`, `src/router/index.ts`, `vite.config.ts`, and
`main.ts` carry no markup and are out of this review's scope.

## Semantic HTML Audit

- The new nav entry in `App.vue` renders through `<router-link>`, which compiles to a real `<a>`
  with visible text ("Reagent Map") — correct interactive element, no `onClick`-on-`div` pattern.
- `ReagentMapView.vue`'s site roster uses a real `<ul>`/`<li>` list, not `div`-soup.
- The view's heading (`<h2>Aurelis Reagent Map</h2>`) follows the same heading-level convention
  already established by the sibling `LocationsView.vue` (page title as `h2`, nested under the
  app shell's single `h1` in `App.vue`'s header) — consistent, no new violation introduced.
- The Cesium mount point (`<div id="cesiumContainer">`) is unavoidably a `div`: Cesium's `Viewer`
  constructor requires a plain container element to inject its canvas into. This is not a
  semantic-HTML violation by itself — the issue was the missing accessible name (see Findings).

## Accessibility (WCAG) Findings

### 1. Fixed — Cesium globe container had no accessible name (WCAG 1.1.1 Non-text Content)
- **File**: `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue`
- **Before**: `<div id="cesiumContainer">` carried no `role` or `aria-label`. A screen reader
  landed on an empty, unlabeled node where the entire globe graphic lives — not even an
  announcement that something is there.
- **Fix applied**: added `role="img"` and a descriptive `aria-label` pointing a screen-reader
  user at the adjacent list as the full-content equivalent, instead of silently hiding the node.
  `role="img"` also stops assistive tech from trying to drill into Cesium's non-semantic canvas
  internals.
- **Why this shape and not `aria-hidden`**: hiding the container outright would be the simpler
  fix, but Cesium's `Viewer` sets `tabindex="0"` on its container for mouse/keyboard camera
  panning. An `aria-hidden` ancestor with a focusable descendant is itself a WCAG 4.1.2 violation
  (hidden-but-focusable), so `aria-hidden` was rejected in favor of giving the node a real name.

### 2. Fixed — site list had no standalone accessible label
- **File**: `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue`
- **Before**: `<ul data-testid="reagent-site-list">` relied entirely on visual proximity to the
  `<h2>` for context.
- **Fix applied**: added `aria-label="Reagent sites"` directly on the `<ul>` so the list is
  self-describing if a screen-reader user reaches it via list/landmark navigation rather than
  linear reading order.

## Keyboard Navigation Audit

- The new nav `<router-link>` is reachable and operable via Tab/Enter exactly like its sibling
  links (Shop, Locations, Cart) — no new keyboard trap.
- The site-name `<ul>` needs no keyboard interaction; it is static text, not actionable.
- **Escalation, not fixed**: Cesium's canvas itself accepts keyboard focus (`tabindex="0"`) for
  camera panning, but the app defines no visible `:focus-visible` style for that container, and
  Cesium's default pan/zoom/rotate keyboard bindings are undocumented in this view and not a
  substitute for a tested keyboard-operable interaction model. Making an arbitrary 3D camera
  genuinely keyboard-operable to WCAG 2.1.1 standard is a fundamental redesign of how the globe is
  driven, not a markup fix — out of this review's authority to apply directly. The adjacent
  site-name list already gives keyboard/screen-reader users the full data (every site, by name)
  without needing to drive the globe at all, so the feature's acceptance criteria are met by the
  list acting as the keyboard/AT-accessible equivalent. Flagging the globe's own keyboard
  operability as a follow-up for human/product judgment, not a blocker for this slice.

## Form Inputs

- No `<input>` or form elements are introduced by this feature. Not applicable.

## Fixes Applied Directly

1. Added `role="img"` + `aria-label` to `#cesiumContainer` in `ReagentMapView.vue`.
2. Added `aria-label="Reagent sites"` to the site-name `<ul>` in `ReagentMapView.vue`.

## Escalations Requiring Human/Product Judgment

- Cesium's camera is mouse-first with no verified, documented keyboard-operable equivalent beyond
  default library keybindings. Redesigning that is out of scope for this slice (spec explicitly
  defers fly-to/interaction tooling to the next three orders) and is noted here so it isn't lost.
- Color contrast of `text-gray-400` body copy against the app's dark theme was not verified —
  static markup analysis only, no automated contrast tooling was run per this agent's mandate.
  Flag for a visual/contrast pass if this theme is audited later.

## Notes for Code Reviewer

- Both fixes are additive `role`/`aria-label` attributes only — no change to component behavior,
  props, events, or the `saturdayMetadata`/`window.__SATURDAY_CESIUM__` seam contract. Safe to
  review as markup-only.