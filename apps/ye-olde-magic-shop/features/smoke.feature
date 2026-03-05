Feature: Web App Smoke Test

  Scenario: Homepage loads
    Given the demo web app is running
    Then the page title should contain "Ye Olde Magic Shop"
