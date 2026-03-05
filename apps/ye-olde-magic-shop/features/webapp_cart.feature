Feature: Demo web-app cart
  As a shopper
  I want to add items to the cart on the demo web app
  So that the cart reflects the added items and totals

  Background:
    Given the demo web app is running

  Scenario: Add first product to cart and verify it appears in cart
    Given I am on the demo shop home page
    When I add the first product to the cart
    And I go to the cart page
    Then I should see the product in the cart with quantity 1
    And the total price should equal the product price

