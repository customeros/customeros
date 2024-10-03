import { chromium } from '@playwright/test';

import { FlowsPage } from './pages/flows/flowsPage';
import { LoginPage } from './pages/loginPage/loginPage';
import { ContactsPage } from './pages/contacts/contactsPage';
import { WinRatesFor } from './pages/opportunitiesKanban/winRates';
import { OpportunitiesPage } from './pages/opportunities/opportunitiesPage';
import { OrganizationsPage } from './pages/organizations/organizationsPage';
import { OpportunitiesKanbanPage } from './pages/opportunitiesKanban/opportunitiesKanbanPage';

async function globalSetup() {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const contactsPage = new ContactsPage(page);
  const flowsPage = new FlowsPage(page);
  const opportunitiesPage = new OpportunitiesPage(page);
  const opportunitiesKanbanPage = new OpportunitiesKanbanPage(page);

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
    console.warn('Select All Contacts button not found or visible:', error);
  }

  if (isSelectAllContactsClicked) {
    await contactsPage.archiveOrgs();
    await contactsPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }

  // Archive flows
  await flowsPage.waitForPageLoad();

  let isSelectAllFlowsClicked = false;

  try {
    isSelectAllFlowsClicked = await flowsPage.selectAllFlows(); // Returns true if successful
  } catch (error) {
    console.warn('Select All Flows button not found or visible:', error);
  }

  if (isSelectAllFlowsClicked) {
    await flowsPage.archiveOrgs();
    await flowsPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }

  // Archive opportunities
  await opportunitiesPage.goToOpportunities();

  let isSelectAllOpportunitiesClicked = false;

  try {
    isSelectAllOpportunitiesClicked =
      await opportunitiesPage.selectAllOpportunities(); // Returns true if successful
  } catch (error) {
    console.warn('Select All Flows button not found or visible:', error);
  }

  if (isSelectAllOpportunitiesClicked) {
    await opportunitiesPage.archiveOrgs();
    await opportunitiesPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }

  // Reset opportunities win rate
  await opportunitiesKanbanPage.goToOpportunitiesKanban();
  await opportunitiesKanbanPage.setWinRates(WinRatesFor.Identified, -100);
  await opportunitiesKanbanPage.setWinRates(WinRatesFor.Qualified, -100);
  await opportunitiesKanbanPage.setWinRates(WinRatesFor.Committed, -100);

  // Create initial organization
  await organizationsPage.goToAllOrgs();

  await organizationsPage.addInitialOrganization();

  await browser.close();
}

export default globalSetup;
