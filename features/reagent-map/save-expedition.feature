Feature: Saving an expedition
  A signed-in adventurer saves a completed line and finds it again.
  Signed out, or with a zero-length line, nothing is saved. An unsaved line is not kept.

  # UI with an API read-back: the only check that the saved shape is right that does not
  # depend on rendering (05, the cross-check).
  @AC11 @level-ui @runner-playwright
  Scenario: A saved expedition is shown again and held by the shop where it was drawn
    Given Elara is signed in and has opened the reagent map
    When Elara draws a line from Moonwell Spring to Ember Crag and saves it
    And Elara opens the reagent map again
    Then her new expedition is shown
    And the shop holds it for Elara with both ends within 1 km of where she placed them

  @AC11 @level-api @runner-playwright
  Scenario: The shop keeps an expedition for the adventurer who saved it
    Given Elara is signed in
    When Elara saves an expedition from Moonwell Spring to Ember Crag
    Then Elara's visible expeditions include it, with both ends within 1 km of those sites

  @AC12 @level-api @runner-playwright
  Scenario: Saving while signed out is refused
    Given no adventurer is signed in
    When an expedition from Moonwell Spring to Ember Crag is saved
    Then the save is refused because no adventurer is signed in
    And the shop holds no new expedition

  # UI: the API cannot show that the adventurer is asked to sign in.
  @AC12 @level-ui @runner-playwright
  Scenario: A signed-out adventurer is asked to sign in when saving
    Given no adventurer is signed in and a completed line is on the reagent map
    When they try to save it
    Then they are asked to sign in
    And the shop holds no new expedition

  @AC13 @level-api @runner-playwright
  Scenario: A zero-length expedition is refused
    Given Elara is signed in
    When Elara saves an expedition whose two ends are less than 1 km apart
    Then the save is refused as a zero-length line
    And the shop holds no new expedition

  # UI: the API cannot show what the adventurer is told.
  @AC13 @level-ui @runner-playwright
  Scenario: Placing both points at the same place cannot be saved
    Given Elara is signed in and has opened the reagent map
    When Elara places both points of a line less than 1 km apart and tries to save it
    Then the save is refused as a zero-length line
    And the shop holds no new expedition

  # UI: abandoning a line is a view behaviour; the shop is checked for no new expedition.
  @AC14 @level-ui @runner-playwright
  Scenario: Starting a new line abandons the unfinished one
    Given Elara is signed in and has opened the reagent map
    And Elara has placed one point of a line
    When Elara starts a new line
    Then the new line has no points
    And the shop holds no new expedition

  @AC15 @level-ui @runner-playwright
  Scenario: An unsaved line is not kept after leaving the map
    Given Elara is signed in and has drawn a line without saving it
    When Elara leaves the reagent map and opens it again
    Then the unsaved line is not shown
    And the shop holds no new expedition
