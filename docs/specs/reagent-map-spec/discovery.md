# Discovery: reagent-map-spec

Gate `confirm-problem`, owner: requester. Generated from `discovery.json`; edit that file, not this one.

## Request

> Ye Olde Magic Shop gets a second map view: a 3D globe of an invented planet. Reagent sites are marked on it, and saved expeditions are drawn on it when the view loads. An adventurer draws shapes to plan an expedition and saves it. Shapes arrive one wave at a time: a straight line first, then circles that enclose sites, then the rest.

## Goal ladder

Five Whys used only to climb from the feature to the outcome behind it.

1. **Why a reagent map?** Adventurers gather reagents at sites spread over the planet and need to see where those sites are relative to each other.
2. **Why draw on it?** An expedition is a route or an area over the sites; drawing it on the globe is how an adventurer plans it.
3. **Why save it?** A plan is worth keeping: the adventurer comes back to it, and the shop can show existing expeditions when the view opens.
4. **Why does this repository want the feature?** It is the pilot that takes one feature through every factory on the line, with a person at each gate, and it gives Saturday a Cesium testing package.

## Goal

An adventurer can plan a reagent-gathering expedition by drawing it on a 3D globe of the invented planet, save it, and find it again; and the pilot proves the factory line can deliver that feature with its spec, build and tests each signed by the role that owns them.

## Outcomes

| | Outcome | How we know |
|---|---|---|
| O1 | An adventurer draws a straight-line expedition between two points on the land, saves it, and sees the same line after reloading the view. | Both saved endpoints are within the agreed tolerance of where they were placed (tolerance: Q3). |
| O2 | Seeded expeditions are drawn in the right places when the view opens. | Every seeded line's endpoints match the seed data within the agreed tolerance. |
| O3 | The camera can be moved to any reagent site. | After flying to a site, the site is on screen and on the near side of the globe. |
| O4 | The view needs no network. | The view and its tests pass with egress blocked: local planet imagery, no Cesium ion. |
| O5 | Every UI and API scenario can fail. | Each scenario fails on at least one injected fault (05, Judges). |
| O6 | Each gate is answered by its owner and the corrections are recorded. | corrections.md has an entry for every gate. |

## Actors

| Actor | Kind | Note |
|---|---|---|
| Adventurer | user | Signed-in shop customer who plans expeditions. |
| Shop | system | Publishes reagent sites and serves saved expeditions (mock API). |
| Analyst | gate role | Owns the domain and the rules. |
| Developer | gate role | Owns the view, the API and the seam; answers the security gate. |
| Tester | gate role | Owns the scenarios and the Cesium package. |

## Impact map

| Actor | Impact | Deliverables |
|---|---|---|
| Adventurer | Plans an expedition by drawing it instead of noting coordinates by hand | reagent map view (D1), line tool, expeditions API |
| Adventurer | Finds saved and seeded expeditions where they were drawn | seeded expeditions, GET /api/expeditions |
| Tester | Automates a 3D globe the way a user does, and checks shapes without reading pixels | @orieken/saturday-playwright-cesium (C1), draw-line flow (E1), API client (F1) |
| Developer | Builds the view against a seam the tests rely on | window.__SATURDAY_CESIUM__ seam contract (C1), entity metadata |

## Questions

Each has the assumption the spec will use if it stays unanswered. `blocking` questions stop the run at their gate.

| | Gate | Blocking | Question | Assumption if unanswered |
|---|---|---|---|---|
| Q1 | `confirm-problem` | yes | Is the goal above right, with the product outcome (O1-O3) and the pilot outcome (O5-O6) both required for wave 1 to be done? | Yes, both. |
| Q2 | `confirm-problem` | yes | Must an adventurer be signed in to save an expedition, and to see the map at all? | Viewing needs no sign-in; saving does, as orders do today. |
| Q3 | `answer-questions` | yes | How close must a saved endpoint be to where it was placed, in kilometres on the planet? | Within 1 km. |
| Q4 | `answer-questions` | yes | Is a line the shortest path over the surface (geodesic) or a straight line on a flat map (rhumb)? | Geodesic, Cesium's default. |
| Q5 | `answer-questions` | yes | What does the adventurer see after clicking space, or the far side of the globe, while drawing? | Nothing is placed and the line tool stays waiting for a point on the land; no message. |
| Q6 | `answer-questions` | yes | Can an adventurer see other adventurers' expeditions? (Trust boundary: the developer's security gate needs it.) | No: each adventurer sees their own expeditions plus the shop's seeded ones. |
| Q7 | `answer-questions` | no | Which expeditions are seeded, and who owns them? | Three seeded lines owned by the shop, crossing well-known sites, one of them reaching the far side. |
| Q8 | `answer-questions` | no | The planet's name and radius, and how many reagent sites. | A sphere of radius 6,000 km; eight sites, at least one near each pole and one at 180 degrees longitude. |
| Q9 | `answer-questions` | no | Wave 2, raised now: does a site on a circle's edge count as enclosed, and do overlapping circles count a site once or twice? | Edge counts as enclosed; a site is counted once per circle. |
| Q10 | `confirm-problem` | no | Is a line with both endpoints at the same place, or a second click before the first is placed, a wave 1 concern? | A zero-length line is rejected; out of scope beyond that. |
