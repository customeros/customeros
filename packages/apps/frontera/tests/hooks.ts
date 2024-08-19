import { test as base } from '@playwright/test';

import { LoginPage } from './pages/loginPage';
import { AddressBookPage } from './pages/addressBookPage';

base.beforeEach(async ({ page }) => {
  const loginPage = new LoginPage(page);
  const addressBookPage = new AddressBookPage(page);

  await loginPage.login();
  await addressBookPage.waitForPageLoad();

  let isSelectAllOrgsClicked = false;

  try {
    isSelectAllOrgsClicked = await addressBookPage.selectAllOrgs(); // Returns true if successful
  } catch (error) {
    console.warn('Select All Orgs button not found or visible:', error);
  }

  if (isSelectAllOrgsClicked) {
    await addressBookPage.archiveOrgs();
    await addressBookPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }
});

// Export the base test object
export const test = base;
