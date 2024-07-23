import { test } from '@playwright/test';

import { LoginPage } from './pages/loginPage';
import { AllOrgsPage } from './pages/allOrgsPage';
import { CustomersPage } from './pages/customersPage';

test.setTimeout(180000);

test('get started link', async ({ page }) => {
  const loginPage = new LoginPage(page);
  const allOrgsPage = new AllOrgsPage(page);
  const customersPage = new CustomersPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await allOrgsPage.waitForPageLoad();

  // Add organization and check new entry
  await allOrgsPage.addOrganization();
  await allOrgsPage.checkNewEntry();

  // Go to Customers page and ensure no new org
  await allOrgsPage.goToCustomersPage();
  await customersPage.ensureNumberOfCustomersExist(0);

  // Go back to All Orgs page
  await allOrgsPage.goToAllOrgsPage();

  // Make the organization a customer
  await allOrgsPage.updateOrgToCustomer();

  // Go to Customers page and ensure we have a new customer
  await allOrgsPage.goToCustomersPage();
  await customersPage.ensureNumberOfCustomersExist(1);
});
