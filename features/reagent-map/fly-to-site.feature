Feature: Flying to a reagent site
  An adventurer brings any reagent site into view, even one on the far side of the globe.

  # UI: where the camera is and what faces it exist only in the view.
  @AC5 @level-ui @runner-playwright
  Scenario: Flying to a site on the far side brings it into view
    Given an adventurer has opened the reagent map
    And Tidewrack Shoals is on the far side of the globe
    When the adventurer flies to Tidewrack Shoals
    Then Tidewrack Shoals is on screen, on the near side of the globe
