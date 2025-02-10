import { chromium } from '@playwright/test';

// import { FlowsPage } from './pages/flows/flowsPage';
import { LoginPage } from './pages/loginPage/loginPage';
import { ContactsPage } from './pages/contacts/contactsPage';
import { CompaniesPage } from './pages/companies/companiesPage';
import { WinRatesFor } from './pages/opportunitiesKanban/winRates';
import { OpportunitiesPage } from './pages/opportunities/opportunitiesPage';
import { OpportunitiesKanbanPage } from './pages/opportunitiesKanban/opportunitiesKanbanPage';

async function globalSetup() {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const loginPage = new LoginPage(page);
  const organizationsPage = new CompaniesPage(page);
  const contactsPage = new ContactsPage(page);
  // const flowsPage = new FlowsPage(page);
  const opportunitiesPage = new OpportunitiesPage(page);
  const opportunitiesKanbanPage = new OpportunitiesKanbanPage(page);

  await loginPage.login();

  // Archive organizations
  await organizationsPage.goToCompanies();

  if ((await page.locator(organizationsPage.allOrgsAddOrg).count()) === 0) {
    await organizationsPage.selectAllOrgs();
    await organizationsPage.archiveOrgs();
    await organizationsPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }

  // Archive contacts
  await contactsPage.waitForPageLoad();

  if (
    (await page
      .locator(contactsPage.finderTableContacts)
      .locator('[data-index]')
      .count()) > 0
  ) {
    await contactsPage.selectAllContacts(); // Returns true if successful
    await contactsPage.archiveOrgs();
    await contactsPage.confirmArchiveOrgs();
    await new Promise((resolve) => setTimeout(resolve, 1500));
  }

  // Archive flows
  // await flowsPage.waitForPageLoad();
  //
  // if (
  //   (await page
  //     .locator(flowsPage.finderTableFlows)
  //     .locator('[data-index]')
  //     .count()) > 0
  // ) {
  //   await flowsPage.selectAllFlows();
  //   await flowsPage.archiveOrgs();
  //   await flowsPage.confirmArchiveOrgs();
  //   await new Promise((resolve) => setTimeout(resolve, 1500));
  // }

  // Archive opportunities
  await opportunitiesPage.goToOpportunitiesList();

  if (
    (await page
      .locator(opportunitiesPage.finderTableOpportunities)
      .locator('[data-index]')
      .count()) > 0
  ) {
    await opportunitiesPage.selectAllOpportunities();
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
  await organizationsPage.goToCompanies();

  await organizationsPage.addInitialOrganization();

  await browser.close();
}

export default globalSetup;
