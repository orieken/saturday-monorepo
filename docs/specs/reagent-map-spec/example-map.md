# Example map: reagent-map-spec

Rules → examples → questions. Gate `answer-questions`, owner: analyst. Generated from `example-map.json`.

Every rule has at least one negative example.

## R1. Opening the reagent map shows every reagent site at its position on the planet.

Events: ReagentMapOpened · Questions: Q8

- *happy*: The shop publishes 8 sites; the map shows 8 markers, each at its site's latitude and longitude.
- *happy*: A site at 180 degrees longitude and one near the north pole are shown where they are, not clipped or wrapped.
- *negative*: A site the shop has not published is not shown.

## R2. Opening the map shows the shop's seeded expeditions, and a signed-in adventurer's own saved expeditions, each where it was drawn.

Events: ReagentMapOpened · Questions: Q6, Q7, Q3

- *happy*: Signed out: the 3 seeded lines are shown and nothing else.
- *happy*: Signed in as Elara, who saved one line: the 3 seeded lines and Elara's line are shown.
- *negative*: Signed in as Elara: a line saved by Borin is not shown.
- *negative*: A seeded line is not drawn more than the agreed tolerance away from its seeded endpoints.

## R3. Flying to a reagent site brings it on screen, on the near side of the globe.

Events: SiteFocused

- *happy*: With the camera over the opposite side of the planet, choosing 'Moonwell Spring' turns the globe until the site is on screen and facing the camera.
- *negative*: Before choosing a site, a site on the far side is not on screen.

## R4. A line is drawn by placing two points on the land; it follows the agreed kind of line between them.

Events: LineStarted, PointPlaced · Questions: Q4, Q3

- *happy*: After choosing the line tool, clicking near Moonwell Spring and then near Ember Crag shows a line between those two places.
- *happy*: Two points far apart give a line that follows the agreed path between them (geodesic or rhumb, Q4).
- *negative*: After only one click, no line is complete and Save is not offered.

## R5. A click that is not on the land of the near side places nothing.

Events: PointRejected · Questions: Q5

- *happy*: With the line tool waiting, a click on a point on the land is placed.
- *negative*: A click on space beside the globe places nothing, and the tool still waits for a point.
- *negative*: A click on the limb of the globe, where the far side would be, places nothing.

## R6. A signed-in adventurer can save a completed line, and sees it again after reopening the map, within the agreed tolerance.

Events: ExpeditionSaved · Questions: Q3

- *happy*: Elara draws a line from Moonwell Spring to Ember Crag, saves it, reloads, and the line is shown with both ends within tolerance of where she placed them.
- *negative*: An endpoint saved 5 km from where it was placed counts as wrong.

## R7. Saving is refused when the adventurer is signed out or the line has zero length; nothing is stored.

Events: ExpeditionRejected · Questions: Q11

- *happy*: Signed in, a line between two different places is saved.
- *negative*: Signed out, pressing Save asks the adventurer to sign in, and after signing in and reopening the map the line is not there.
- *negative*: Both points placed at the same place: Save is refused and nothing is stored.

## R8. An unsaved line can be abandoned; nothing is stored.

Events: LineStarted

- *happy*: Elara places one point, chooses the line tool again, and starts over from no points.
- *negative*: Elara draws a line and leaves the map without saving; reopening it shows no new line.

## Questions

Blocking questions stop the run at `answer-questions`. Names used in examples (Moonwell Spring, Ember Crag, Elara, Borin) are placeholders until Q7 and Q8 are answered.

| | Blocking | Question | Assumption if unanswered |
|---|---|---|---|
| Q3 | yes | How close must a saved or seeded endpoint be to where it was placed, in km on the planet? | Within 1 km. |
| Q4 | yes | Is a line the shortest path over the surface (geodesic) or straight on a flat map (rhumb)? | Geodesic. |
| Q5 | yes | After a click off the land, does the adventurer see anything (a message, a shake), or does nothing happen? | Nothing happens; the tool keeps waiting. |
| Q6 | yes | Can an adventurer see other adventurers' expeditions? | No: own plus seeded only. |
| Q7 | no | Which expeditions are seeded, and who owns them? | 3 shop-owned lines crossing named sites, one reaching the far side from the default camera. |
| Q8 | no | Planet name, radius, and the reagent sites. | Planet 'Aurelis', radius 6,000 km, 8 named sites including one near each pole and one at 180 degrees longitude. |
| Q9 | no | Wave 2: does a site on a circle's edge count as enclosed, and do overlapping circles count a site once or twice? | Edge counts; once per circle. |
| Q11 | no | New: when two points are 'the same place', is that exactly equal, or closer than the tolerance? | Closer than the tolerance (Q3). |
