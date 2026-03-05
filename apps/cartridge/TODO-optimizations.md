# Optimization & Feature Wishlist

## Code Quality & Refactoring
- [ ] **Component Splitting**: Break down `CucumberDemoWidget.vue` into smaller sub-components:
  - `FeatureRow.vue`: For individual feature rendering.
  - `ScenarioRow.vue`: For individual scenario rendering.
- [ ] **Type Safety**: Eliminate `any` types (e.g., in `toggleScenario`, `toggleFeature`). Use shared interfaces from a `types` module.
- [ ] **Style Extraction**: Extract repeated Tailwind classes (like button styles, pill badges) into reusable base components or CSS classes to reduce clutter.

## UX Enhancements
- [ ] **Global Controls**: Add "Expand All" and "Collapse All" buttons for features/scenarios.
- [ ] **Search & Filter**: Implement a client-side search bar to filter visible features/scenarios by name or tag.
- [ ] **Run History**: Improve visualization of run history. Currently, we can only see the "latest" status per scenario. A dedicated view or drawer for past runs would be valuable.
- [ ] **Keyboard Navigation**: Ensure full keyboard accessibility for all interactive elements (rows, buttons, toggles), confirming `tabindex` and `aria` attributes are correct.

## Performance
- [ ] **Virtualization**: Implement `vue-virtual-scroller` (or similar) for the main feature list to support large test suites efficiently.

## Testing
- [ ] **Unit Tests**: Add tests for Pinia stores (`runs.ts`, `cucumber.ts`) to ensure state logic (persistence, polling) works as expected.
- [ ] **Component Tests**: Add basic mount tests for `CucumberDemoWidget` to verify rendering state.
