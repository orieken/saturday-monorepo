import { Given, Then, When } from '@cucumber/cucumber';
import { TestWorld } from '../support/world';
import { expect } from '@playwright/test';

Given('I am logged in as {string}', async function (this: TestWorld, email: string) {
  if (!this.site) throw new Error('Site not initialized');
  
  // Navigate to login page
  await this.site.loginPage.navigate();
  
  // Use the demo password for all test users
  await this.site.loginPage.login(email, 'password123');
  
  // Wait for redirect to account page
  await this.page!.waitForURL('**/account', { timeout: 5000 });
});


Given('I am on the Past Orders page', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  // Navigate to account page and switch to orders tab
  await this.site.accountPage.navigate();
  await this.site.accountPage.switchToOrderHistoryTab();
});

When('I view the details of my most recent order', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  // Click on the first order
  await this.site.accountPage.clickFirstOrder();
});

Then('I should see the order date, total, and line items', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  // Verify order details are visible
  const details = await this.site.accountPage.getOrderDetails();
  expect(details.date).toBeTruthy();
  expect(details.total).toBeTruthy();
  expect(details.items).toBeGreaterThan(0);
});

When('I log out', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  await this.site.accountPage.logout();
});

Then('I should be redirected to the home page', async function (this: TestWorld) {
  if (!this.page) throw new Error('Page not initialized');
  
  await this.page.waitForURL('**/', { timeout: 5000 });
  expect(this.page.url()).toContain('/');
});

Then('I should not see account information', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  const isVisible = await this.site.accountPage.isAccountInfoVisible();
  expect(isVisible).toBe(false);
});

Given('I am on the login page', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  await this.site.loginPage.navigate();
});



Then('I should see a login error message', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  const isVisible = await this.site.loginPage.isErrorVisible();
  expect(isVisible).toBe(true);
});


Then('I should remain on the login page', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  const onLoginPage = await this.site.loginPage.isOnLoginPage();
  expect(onLoginPage).toBe(true);
});


When('I log in with username {string} and password {string}', async function (this: TestWorld, email: string, password: string) {
  if (!this.site) throw new Error('Site not initialized');
  
  await this.site.loginPage.login(email, password);
});

Then('I should be redirected to my account dashboard', async function (this: TestWorld) {
  if (!this.page) throw new Error('Page not initialized');
  
  await this.page.waitForURL('**/account', { timeout: 5000 });
  expect(this.page.url()).toContain('/account');
});


Then('I should see a list of my past orders', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  // Switch to orders tab
  await this.site.accountPage.switchToOrderHistoryTab();
  
  // Verify orders are visible
  const orderCount = await this.site.accountPage.getOrderCount();
  expect(orderCount).toBeGreaterThan(0);
});