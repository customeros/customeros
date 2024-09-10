import { chromium } from '@playwright/test';

import { LoginPage } from './pages/loginPage';
import { ContactsPage } from './pages/contactsPage';
import { OrganizationsPage } from './pages/organizations/organizationsPage';

async function globalSetup() {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const contactsPage = new ContactsPage(page);

  await loginPage.login();

  // Archive organizations
  await organizationsPage.goToAllOrgs();

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

  await organizationsPage.goToAllOrgs();

  // Create initial organization
  await organizationsPage.addInitialOrganization();

  await browser.close();
}

export default globalSetup;
