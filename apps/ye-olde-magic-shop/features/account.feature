Feature: Account login and past orders
  As a returning customer
  I want to log into my account
  So that I can see my past orders

  Background:
    Given the demo web app is running

  @happy_path
  Scenario: Successful login shows the Past Orders page
    Given I am on the login page
    When I log in with username "alice@example.com" and password "password123"
    Then I should be redirected to my account dashboard
    And I should see a list of my past orders

  @negative
  Scenario: Invalid password shows an error
    Given I am on the login page
    When I log in with username "alice@example.com" and password "wrong-password"
    Then I should see a login error message
    And I should remain on the login page

  @logout
  Scenario: Logout hides account information
    Given I am logged in as "alice@example.com"
    When I log out
    Then I should be redirected to the home page
    And I should not see account information

  @orders
  Scenario: View details of a past order
    Given I am logged in as "alice@example.com"
    And I am on the Past Orders page
    When I view the details of my most recent order
    Then I should see the order date, total, and line items
