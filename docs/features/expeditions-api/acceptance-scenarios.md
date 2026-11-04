---
feature: "expeditions-api"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: expeditions-api

## Seed data contains the three shop expeditions and the two adventurers' own expeditions

Verifies: The seed data MUST contain exactly the three shop-owned expeditions in seed.seededExpeditions, and adventurers Elara and Borin MUST exist as users with the expeditions seed.adventurers assigns them.

```gherkin
Scenario: Seed data contains the three shop expeditions and the two adventurers' own expeditions
  Given The mock API's data files hold the expeditions-api seed and have not been hand-edited
  When The expeditions and users seed data are loaded
  Then Exactly three expeditions owned by the shop are present: Spring to Crag (Moonwell Spring to Ember Crag), Glade to Caves (Thornveil Glade to Glimmerdeep Caves), and Cairn to Shoals (Saltmarsh Cairn to Tidewrack Shoals)
  And A user named Elara exists and owns an expedition from Moonwell Spring to Frostcap Hollow
  And A user named Borin exists and owns an expedition from Ember Crag to Saltmarsh Cairn
```

## An anonymous visitor sees only the shop's seeded expeditions

Verifies: GET /api/expeditions with no session MUST return exactly the three seeded expeditions (spec AC3).

```gherkin
Scenario: An anonymous visitor sees only the shop's seeded expeditions
  Given No Authorization header is sent with the request
  When The visitor requests GET /api/expeditions
  Then The response contains exactly the three seeded shop expeditions
  And No adventurer-owned expedition appears in the response
```

## Elara sees the seeded expeditions plus her own, and never sees Borin's

Verifies: GET /api/expeditions with Elara's session MUST return the seeded expeditions plus Elara's own, and MUST NOT ever return Borin's (spec AC4). This is the trust boundary the security gate checks.

```gherkin
Scenario: Elara sees the seeded expeditions plus her own, and never sees Borin's
  Given Elara is signed in with a valid session
  When Elara requests GET /api/expeditions
  Then The response contains the three seeded shop expeditions
  And The response contains Elara's own expedition from Moonwell Spring to Frostcap Hollow
  And The response does not contain Borin's expedition from Ember Crag to Saltmarsh Cairn
```

## A signed-in adventurer's new expedition is stored under their own identity, never a forged owner from the request body

Verifies: POST /api/expeditions with a session and a line whose two endpoints are at least 1 km apart on Aurelis (sphere, radius 6,000 km) MUST store the expedition for the session's adventurer, never an owner the request body supplies, and the stored expedition MUST then be returned by GET (spec AC11).

```gherkin
Scenario: A signed-in adventurer's new expedition is stored under their own identity, never a forged owner from the request body
  Given An adventurer is signed in with a valid session
  And The submitted line's two endpoints are at least 1 km apart on Aurelis
  And The request body also supplies a different owner or userId field
  When The adventurer submits POST /api/expeditions
  Then The expedition is stored under the signed-in adventurer's own identity, not the owner supplied in the body
  And A subsequent GET /api/expeditions by the same adventurer includes the newly stored expedition
```

## An anonymous visitor cannot create an expedition

Verifies: POST /api/expeditions with no session MUST be refused with 401 and MUST NOT store anything (spec AC12).

```gherkin
Scenario: An anonymous visitor cannot create an expedition
  Given No Authorization header is sent with the request
  When The visitor submits POST /api/expeditions with an otherwise valid line
  Then The response status is 401
  And No expedition is stored
```

## An expedition whose endpoints are less than 1 km apart is refused as a zero-length line

Verifies: POST /api/expeditions with endpoints less than 1 km apart MUST be refused with 422 naming a zero-length line, and MUST NOT store anything (spec AC13).

```gherkin
Scenario: An expedition whose endpoints are less than 1 km apart is refused as a zero-length line
  Given An adventurer is signed in with a valid session
  And The submitted line's two endpoints are less than 1 km apart on Aurelis
  When The adventurer submits POST /api/expeditions
  Then The response status is 422
  And The error names a zero-length line
  And No expedition is stored
```

## Running the expeditions test suite twice leaves the repository unchanged

Verifies: Tests MUST run through node:test via the mock-api's existing `test` script (node --test, added by reagent-sites-api) and MUST leave the repository unchanged: running the suite twice MUST show no git status change, including under apps/mock-api/data/. POST writes to a data file, as registering a user does, so tests MUST isolate or restore what they write.

```gherkin
Scenario: Running the expeditions test suite twice leaves the repository unchanged
  Given The mock API's node:test suite, including the expeditions tests, has been run once already
  When The suite is run a second time via the existing `test` script
  Then git status shows no changes under apps/mock-api/data/ or anywhere else in the repository after either run
```

## The implementation stays within its complexity budget, test-coverage bar, and file scope

Verifies: Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths (apps/mock-api/**, docs/features/expeditions-api/**) MUST change.

```gherkin
Scenario: The implementation stays within its complexity budget, test-coverage bar, and file scope
  Given The expeditions-api implementation is complete
  When The changed code is reviewed for complexity, coverage, and file scope
  Then Every new or changed function has cyclomatic complexity below 7
  And Every changed line is covered by a test that fails if the behavior it checks is broken
  And No file outside apps/mock-api/** or docs/features/expeditions-api/** has been modified
```
