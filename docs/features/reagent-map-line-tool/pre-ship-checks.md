# Pre-ship checks: reagent-map-line-tool

Run by hand with `scripts/ui-suite-in-chrome.sh`, one run at a time, with the final code.

| Check | Run | Result |
|---|---|---|
| `ui-suite` (debian) | 1 | 77 passed |
| `ui-suite-ubi9` | 1 | 70 passed, 7 failed |
| `ui-suite-ubi9` | 2 | 77 passed |

The 7 failures in ubi9's first run were every test in `tests/map.spec.ts`, all timing out in
`page.goto('/locations')`. That page is the Leaflet map, which this order does not touch. It loads
its tiles from `tile.openstreetmap.org`, so a slow response from that server holds up the page's load
event. Draft order `leaflet-local-tiles` serves the tiles locally. The re-run passed.
