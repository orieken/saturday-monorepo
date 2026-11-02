# Domain model: reagent-map-spec

Event storm draft. Gate `confirm-domain`, owner: analyst. Generated from `domain-model.json`.

## Bounded contexts

- **Expedition Planning** (new): owns Expedition, ReagentSite, Planet.
- **Identity**: owns Adventurer session. Existing sign-in (mock API sessions). Expedition Planning reads who is signed in; it does not change sign-in.

## Aggregates

### Expedition

Fields: id; owner (an adventurer, or the shop for seeded ones); shape: a line with exactly two endpoints, each a latitude and longitude on the planet; savedAt.

- Both endpoints are on the planet's surface.
- The two endpoints are not the same place (zero-length line rejected, Q10).

### ReagentSite

Fields: id; name; latitude; longitude.

- Published by the shop; adventurers cannot add or move one in wave 1.

### Planet

Fields: name; radius in km; surface image.

- A perfect sphere; distances on it use its own radius, never Earth's (D-13, D-14).

## Commands and events

| Command | Issued by | Events |
|---|---|---|
| OpenReagentMap | Adventurer (signed in or not) | ReagentMapOpened |
| FlyToSite | Adventurer | SiteFocused |
| StartLine | Adventurer | LineStarted |
| PlacePoint | Adventurer | PointPlaced, PointRejected |
| SaveExpedition | Signed-in adventurer | ExpeditionSaved, ExpeditionRejected |

| Event | Means |
|---|---|
| ReagentMapOpened | The globe is shown with the reagent sites and the expeditions this adventurer may see. |
| SiteFocused | The camera has moved so a chosen reagent site is on screen, on the near side of the globe. |
| LineStarted | The adventurer chose to draw a line; no point is placed yet. |
| PointPlaced | A point on the land became an endpoint of the line being drawn. The second one completes the line. |
| PointRejected | A click on space or on the far side of the globe placed nothing (Q5). |
| ExpeditionSaved | A completed line is stored for its owner and is shown again when the view opens. |
| ExpeditionRejected | Saving was refused: not signed in, or the line has zero length. |

## Policies

- When **ReagentMapOpened**, then Show every reagent site, the shop's seeded expeditions, and, for a signed-in adventurer, their own saved expeditions (Q6).
- When **PointPlaced (second point)**, then The line is complete and can be saved.
- When **ExpeditionSaved**, then The expedition is among that adventurer's visible expeditions from then on.

## Read models

- **ReagentSites**: Every site's name and position. (from the shop's reagent-site data)
- **VisibleExpeditions**: Seeded expeditions plus the signed-in adventurer's own; seeded only when signed out. (from the shop's expedition data, filtered by owner)

## Hotspots

| | Question | Assumption if unanswered |
|---|---|---|
| H1 | Does an expedition need a name or anything else besides its line when saved? | No: a line and its owner. A name can come in a later wave. |
| H2 | Can an adventurer abandon a line before saving it, and start again? | Yes: choosing the line tool again, or leaving the view, discards the unsaved line. Nothing is stored. |
| H3 | Can a signed-out adventurer draw a line at all, or only not save it? | They can draw; Save tells them to sign in. The line is not kept. |
| H4 | What does a reagent site carry besides its name and position (the reagent it yields, a description)? | Name and position only in wave 1. |
| H5 | How does the adventurer choose the site to fly to: from a list beside the globe, or by clicking its marker? | From a list of site names beside the globe. |
| H6 | Is 'Expedition' the right word for a saved shape, and 'Reagent site' for a place? These become domain terms in DOMAIN_DICTIONARY.md. | Yes, as used in prompt 05. |
