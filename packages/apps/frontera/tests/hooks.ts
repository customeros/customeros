import { test as base } from '@playwright/test';

import { LoginPage } from './pages/loginPage';
import { AddressBookPage } from './pages/addressBookPage';

base.beforeEach(async ({ page }) => {
  // Teardown logic after each test
  const loginPage = new LoginPage(page);
  const addressBookPage = new AddressBookPage(page);

  await loginPage.login();
  await addressBookPage.waitForPageLoad();
  await addressBookPage.selectAllOrgs();
  await addressBookPage.archiveOrgs();
  await addressBookPage.confirmArchiveOrgs();
  await new Promise((resolve) => setTimeout(resolve, 1500));
});

// Export the base test object
export const test = base;
