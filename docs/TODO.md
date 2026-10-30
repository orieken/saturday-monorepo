
# TODO — Monorepo
- [ ] Add GitHub Actions CI: build, lint, test, export k6 smoke
- [ ] Add release pipeline (changesets or semantic-release)
- [ ] Add docs site build (VitePress or Docusaurus)
- [ ] Add JSON-based policy registry mapping

# TODO — Factory work orders

Each item below with an `order:` tag is a draft work order in `intake/`. To have loom build it, set
`status: ready`, commit it to `main`, and run `loom factory run`. Setup and order of play:
`docs/factory-setup.md`.

- [ ] **Run the UI suite in real Chrome.** `playwright.config.ts` sets no channel, so the suite runs
      in Playwright's bundled Chromium on the host; `Dockerfile.runner` is the Playwright image, not
      a real-Chrome one; and `containers/Dockerfile` copies a `supervisord.conf` that is not in
      `containers/`. Pick one real-Chrome image, run the suite in it under Xvfb through one script,
      and assert WebGL is available. Every map order depends on this. Image design: `docs/runner-image.md`.
      `order: ui-suite-real-chrome`
- [ ] **Polyglot runner target**: JDK, Python, Gradle and Go on top of the runner image, pinned and
      checksum-verified. `order: runner-image-polyglot`
- [ ] **UBI9 runner variant** (supported alongside Debian, for teams needing a Red Hat base): same image on UBI9 with Stream and EPEL, packages locked, EPEL RPMs
      vendored, Mesa from Stream so WebGL works. `order: runner-image-ubi9`

# TODO — Golden tests (repair before any factory copies them)

Agents left the demo app breaking its own rules. A factory pins these files as the pattern to follow,
so each item here is a defect it would copy. See ai-assistant-dot-files
`docs/prompts/factories/03-test-factory-saturday.md`, Phase 0.

- [ ] Move `apps/ye-olde-magic-shop` off its vendored `./saturday-core` and `./saturday-cucumber`
      copies onto the workspace packages. The copies have drifted (`managers/` missing; `index.ts`
      and `decorators/requires-filter.ts` differ). `order: golden-workspace-packages`
- [ ] Rename `lib/site/*.ts` to `name.type.ts` (`cart.page.ts`, `magic-shop.site.ts`).
      `order: golden-site-file-names`
- [ ] Register real flows on `MagicShop` (at least login and checkout); `initializeFlows()` is empty.
      `order: golden-flows`
- [ ] Make step definitions talk only to `this.site`: remove the 25 uses of `this.page` and
      `getByTestId` across `features/step_definitions/`. `order: golden-steps-through-site`
- [ ] Add `data-testid` for every element the steps touch; remove the guessed
      `button:has-text("COMPLETE TRANSACTION")` selector in `CartPage`. `order: golden-test-ids`
- [ ] Add a lint that enforces the above in `step_definitions/` and `lib/site/`, so the next agent
      cannot regress it, and write `GOLDEN.md` naming the repaired files as the exemplar set.
      `order: golden-lint-and-exemplars`

# TODO — Maps

- [ ] **Leaflet: no network.** `LocationsView.vue` loads tiles from `tile.openstreetmap.org`, so
      `tests/map.spec.ts` fails wherever there is no egress. Serve a local image (`L.imageOverlay`
      with `CRS.Simple`) or local tiles. `order: leaflet-local-tiles`
- [ ] **Leaflet package: tests.** `saturday-playwright-leaflet` has none (its `test` script exits 1).
      Add a static fixture page and a suite run against it. `order: leaflet-package-tests`
- [ ] **Leaflet package: real gestures.** `drawLine` adds a polyline through Leaflet's API, so it
      cannot catch a broken drawing tool. Add `toScreen`, `clickAt` and `dragFrom` through
      `page.mouse`, and keep `drawLine` only as a labelled test-data helper.
      `order: leaflet-package-gestures`
- [ ] **Fix `tests/map.spec.ts`** to use the gestures above. `order: leaflet-map-spec-gestures`
- [ ] **New package `saturday-playwright-cesium`**: gestures and state reads, `waitForSettled`, a
      seam contract (`window.__SATURDAY_CESIUM__`), its own fixture page and suite. Local imagery, no
      Cesium ion account or token; an invented planet. Pilot order C1.
      `order: saturday-playwright-cesium`
- [ ] **New view: the reagent map** (3D Cesium globe, reagent sites, seeded expeditions, line tool),
      plus `GET/POST /api/expeditions` and seed data in `apps/mock-api`. Pilot order D1. No order
      yet: it is written from the reagent-map spec, which is done by hand first.
- [ ] **New tests: reagent map UI and API suites.** Pilot orders E1 and F1, by hand first.
- [ ] **Consider `saturday-geo`**: the Turf-based geometry both map packages use, as a pure package
      tested in Node. Decided at `saturday-playwright-cesium`'s design gate.
- [ ] **Heatmap for maps**: a Cesium (then Leaflet) source for `saturday-playwright-heatmap` that
      records clicks as lat/lng and reports map entities as interactables, so
      `saturday-ml-analyzer` can report reagent sites and shapes no test touched. After wave 1.
      `order: map-heatmap-source`
