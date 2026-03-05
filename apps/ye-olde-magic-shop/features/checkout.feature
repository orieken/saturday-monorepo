Feature: Checkout Process

  As a shopper
  I want to complete my purchase
  So that I can receive my magical items

  Background:
     Given the demo web app is running
     Given I am on the login page
     When I log in with username "alice@example.com" and password "password123"
     Given I am on the demo shop home page
     When I add the first product to the cart
     And I go to the cart page

  @checkout
  Scenario: Complete a transaction successfully
     When I complete the transaction
     Then I should be redirected to the order history page
     And I should see the new order in the history
