import { test as base } from '@playwright/test';

import { LoginPage } from './pages/loginPage';
import { ContactsPage } from './pages/contactsPage';
import { OrganizationsPage } from './pages/organizationsPage';

// Define a custom test fixture
const test = base.extend({
  // Custom fixtures can be defined here if needed
});

test.beforeAll(async ({ browser }) => {
  // Create a new page instance
  const page = await browser.newPage();

  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const contactsPage = new ContactsPage(page);

  await loginPage.login();

  // Archive organizations
  await organizationsPage.waitForPageLoad();

  let isSelectAllOrgsClicked = false;

  try {
    isSelectAllOrgsClicked = await organizationsPage.selectAllOrgs(); // Returns true if successful
  } catch (error) {
    console.warn('Select All Orgs button not found or visible:', error);
  }

  if (isSelectAllOrgsClicked) {
    await organizationsPage.archiveOrgs();
    await organizationsPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }

  // Archive contacts
  await contactsPage.waitForPageLoad();

  let isSelectAllContactsClicked = false;

  try {
    isSelectAllContactsClicked = await contactsPage.selectAllContacts(); // Returns true if successful
  } catch (error) {
    console.warn('Select All Orgs button not found or visible:', error);
  }

  if (isSelectAllContactsClicked) {
    await contactsPage.archiveOrgs();
    await contactsPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }

  await organizationsPage.waitForPageLoad();
  await organizationsPage.addInitialOrganization();

  // Close the page after setup
  await page.close();
});

// Export the custom test object
export { test };
