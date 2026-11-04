---
feature: "reagent-sites-api"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: reagent-sites-api

## The seeded reagent sites match the specification exactly

Verifies: apps/mock-api/data/reagent-sites.json MUST hold exactly the eight sites in docs/specs/reagent-map-spec/acceptance-criteria.json's seed.reagentSites, each with an id, name, lat and lng.

```gherkin
Scenario: The seeded reagent sites match the specification exactly
  Given The reagent sites seed data file has been created for the mock API
  When The seed data is compared against the specification's seed list of reagent sites
  Then Exactly eight reagent sites are present
  And Each site has an id, a name, a latitude, and a longitude
  And Every site's id, name, latitude, and longitude matches the specification's seed list
```

## A request with no session still receives all reagent sites at their original positions

Verifies: GET /api/reagent-sites MUST return all eight sites with their positions unchanged, whether or not the request carries a session (spec AC1; features/reagent-map/reagent-sites.feature @AC1 @level-api).

```gherkin
Scenario: A request with no session still receives all reagent sites at their original positions
  Given No session is present on the request
  When A client requests the list of reagent sites
  Then The response contains all eight reagent sites
  And Every site's position is unchanged from the seed data
```

## A request with a valid session receives all reagent sites at their original positions

Verifies: GET /api/reagent-sites MUST return all eight sites with their positions unchanged, whether or not the request carries a session (spec AC1; features/reagent-map/reagent-sites.feature @AC1 @level-api).

```gherkin
Scenario: A request with a valid session receives all reagent sites at their original positions
  Given A client holds a valid session
  When The client requests the list of reagent sites
  Then The response contains all eight reagent sites
  And Every site's position is unchanged from the seed data
```

## A request with an invalid or garbage session still receives all reagent sites at their original positions

Verifies: GET /api/reagent-sites MUST return all eight sites with their positions unchanged, whether or not the request carries a session (spec AC1; features/reagent-map/reagent-sites.feature @AC1 @level-api).

```gherkin
Scenario: A request with an invalid or garbage session still receives all reagent sites at their original positions
  Given A client holds an invalid or garbage session credential
  When The client requests the list of reagent sites
  Then The response contains all eight reagent sites
  And Every site's position is unchanged from the seed data
  And The request is not rejected on account of the session credential
```

## The mock API's test suite runs on Node's built-in test runner with no added dependency

Verifies: apps/mock-api MUST run its tests through Node's built-in node:test, with no new dependency, via a test script in apps/mock-api/package.json.

```gherkin
Scenario: The mock API's test suite runs on Node's built-in test runner with no added dependency
  Given The mock API project as delivered
  When The test script declared in apps/mock-api/package.json is run
  Then The suite executes using Node's built-in node:test runner
  And No new test-framework dependency has been added to run it
```

## The implementation stays within complexity, coverage, and scope limits

Verifies: Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change, excepting the run's own delivery record under docs/features/reagent-sites-api/.

```gherkin
Scenario: The implementation stays within complexity, coverage, and scope limits
  Given The feature has been implemented
  When The change is reviewed for complexity, test coverage, and file scope
  Then Every new or changed function has cyclomatic complexity below 7
  And Every changed line is covered by a test that can fail
  And No file outside apps/mock-api/** has changed, other than the delivery record under docs/features/reagent-sites-api/
```
