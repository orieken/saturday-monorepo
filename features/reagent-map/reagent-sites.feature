Feature: Reagent sites on the reagent map
  Adventurers see where every reagent site on Aurelis is.
  Positions are correct within 1 km on the planet's surface.

  @AC1 @level-api @runner-playwright
  Scenario: The shop publishes every reagent site
    Given the shop publishes the eight reagent sites of Aurelis
    When an adventurer asks the shop for its reagent sites
    Then all eight sites are returned with their published positions

  # UI: the API cannot show that each site is drawn at its position on the globe.
  @AC1 @level-ui @runner-playwright
  Scenario: Every reagent site is shown at its position
    Given the shop publishes the eight reagent sites of Aurelis
    When an adventurer opens the reagent map
    Then all eight sites are listed beside the globe
    And each site is shown within 1 km of its published position

  # UI: clipping and repetition at the poles and the 180th meridian only happen when drawing.
  @AC2 @level-ui @runner-playwright
  Scenario Outline: Sites near a pole and on the 180th meridian are shown once, where they are
    Given an adventurer has opened the reagent map
    When the adventurer turns the globe to <site>
    Then <site> is shown once, within 1 km of its published position

    Examples:
      | site             |
      | Frostcap Hollow  |
      | Tidewrack Shoals |
