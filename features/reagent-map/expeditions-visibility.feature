Feature: Which expeditions an adventurer sees
  Adventurers see the shop's seeded expeditions and their own, never another adventurer's.
  Signed out, they see the seeded expeditions only.

  Background:
    Given Elara has saved the expedition "Moonwell Spring to Frostcap Hollow"
    And Borin has saved the expedition "Ember Crag to Saltmarsh Cairn"

  @AC3 @level-api @runner-playwright
  Scenario: Signed out, only the seeded expeditions are visible
    Given no adventurer is signed in
    When the visible expeditions are requested from the shop
    Then exactly the three seeded expeditions are returned

  @AC4 @level-api @runner-playwright
  Scenario: An adventurer sees their own expeditions and not another adventurer's
    Given Elara is signed in
    When Elara requests her visible expeditions from the shop
    Then the three seeded expeditions and "Moonwell Spring to Frostcap Hollow" are returned
    And "Ember Crag to Saltmarsh Cairn" is not returned

  # UI: the API cannot show that each expedition is drawn between the right places.
  @AC3 @level-ui @runner-playwright
  Scenario: Seeded expeditions are drawn where they were seeded
    Given no adventurer is signed in
    When the reagent map is opened
    Then "Spring to Crag", "Glade to Caves" and "Cairn to Shoals" are shown
    And each has both ends within 1 km of its seeded sites

  # UI: the API cannot show what the map draws for a signed-in adventurer.
  @AC4 @level-ui @runner-playwright
  Scenario: A signed-in adventurer's own expedition is drawn with the seeded ones
    Given Elara is signed in
    When Elara opens the reagent map
    Then the three seeded expeditions and "Moonwell Spring to Frostcap Hollow" are shown
    And "Ember Crag to Saltmarsh Cairn" is not shown
