Feature: Drawing a line
  An adventurer draws an expedition by placing two points on the land of Aurelis.
  The line follows the shortest path over the surface. Points off the land are not placed.

  Background:
    Given an adventurer has opened the reagent map
    And the adventurer has started a line

  # UI: drawing is the behaviour under test; nothing reaches the shop until a save.
  @AC6 @level-ui @runner-playwright
  Scenario: Two points on the land make a line that can be saved
    When the adventurer places a point at Moonwell Spring
    And the adventurer places a point at Ember Crag
    Then a completed line joins Moonwell Spring and Ember Crag
    And the line can be saved

  @AC7 @level-ui @runner-playwright
  Scenario: One point is not a line
    When the adventurer places a point at Moonwell Spring
    Then no line is complete
    And the line cannot be saved

  # UI: the shop holds only the two ends; the path between them is only in the drawing.
  @AC8 @level-ui @runner-playwright
  Scenario: A line follows the shortest path over the surface
    When the adventurer places a point at Thornveil Glade
    And the adventurer places a point at Glimmerdeep Caves
    Then the line follows the shortest path over the surface of Aurelis between them

  @AC9 @level-ui @runner-playwright
  Scenario: A point on space beside the globe is not placed
    When the adventurer tries to place a point on space beside the globe
    Then no point is placed
    And the line tool is still waiting for its first point

  @AC10 @level-ui @runner-playwright
  Scenario: A point on the edge of the globe, where the land faces away, is not placed
    When the adventurer tries to place a point on the edge of the globe
    Then no point is placed
