// import { test } from '@playwright/test';
import { test } from './videoFixture';
import { flow, organizations } from './test-data';
import { FlowPage } from './pages/flows/flowPage';
import { FlowsPage } from './pages/flows/flowsPage';
import { LoginPage } from './pages/loginPage/loginPage';
import { ContactsPage } from './pages/contacts/contactsPage';
// import { ContactsPage } from './pages/contacts/contactsPage';
import { CustomersPage } from './pages/customers/customersPage';
import { WinRatesFor } from './pages/opportunitiesKanban/winRates';
import { KanbanColumns } from './pages/opportunitiesKanban/columns';
import { OrganizationsPage } from './pages/organizations/organizationsPage';
import { OrganizationAboutPage } from './pages/organization/organizationAboutPage';
import { OrganizationsCmdKPage } from './pages/organizations/organizationsCmdKPage';
import { OrganizationPeoplePage } from './pages/organization/organizationPeoplePage';
import { OrganizationAccountPage } from './pages/organization/organizationAccountPage';
import { OrganizationSideNavPage } from './pages/organization/organizationSideNavPage';
import { OrganizationTimelinePage } from './pages/organization/organizationTimelinePage';
import { OpportunitiesKanbanPage } from './pages/opportunitiesKanban/opportunitiesKanbanPage';

test.setTimeout(180000);

test('Convert an Organization to Customer', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const customersPage = new CustomersPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.goToAllOrgs();

  // Add organization and check new entry
  const organizationName = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  await organizationsPage.checkNewOrganizationEntry(organizationName);

  // Go to Customers page and ensure no new org
  await organizationsPage.goToCustomersPage();
  // await customersPage.ensureNumberOfCustomersExist(0);
  await customersPage.ensureCustomerExists(organizationName, false);

  // Go back to All Orgs page
  await organizationsPage.goToAllOrgsPage();

  // Make the organization a customer
  await organizationsPage.updateOrgToCustomer(organizationName);

  // Go to Customers page and ensure we have a new customer
  await organizationsPage.goToCustomersPage();
  // await customersPage.ensureNumberOfCustomersExist(1);
  await customersPage.ensureCustomerExists(organizationName, true);
});

test('Add About information to an Organization', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationAboutPage = new OrganizationAboutPage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.goToAllOrgs();

  // Add organization and check new entry
  const organizationName = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationName);

  await organizationSideNavPage.goToAbout();

  //Check enrichment
  await organizationAboutPage.enrichOrganization(organizations.create.website);
  organizations.create.name = organizationName;
  await organizationAboutPage.checkEnrichedAboutFields(organizations.create);

  //Check updates that override the enrichment
  await organizationAboutPage.populateAboutFields(organizations.update);
  await organizationAboutPage.checkPopulatedAboutFields(organizations.update);
});

test('Create People entry in an Organization', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationPeoplePage = new OrganizationPeoplePage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.goToAllOrgs();

  // Add organization and check new entry
  const organizationName = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationName);

  // Go to People page
  await organizationSideNavPage.goToPeople();
  await organizationPeoplePage.createContactFromEmpty();
});

test('Create Timeline entries in an Organization', async ({
  page,
}, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);
  const organizationTimelinePage = new OrganizationTimelinePage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.goToAllOrgs();

  // Add organization and check new entry
  const organizationName = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationName);

  // Go to Account page and update org
  await organizationSideNavPage.goToAccount();
  await organizationTimelinePage.ensureEmailPermissionPromptIsRedirecting();
  await page.goBack();
  await organizationTimelinePage.ensureLogEntryCanBeAdded();
  await organizationTimelinePage.ensureReminderCanBeAdded();
});

test('Create Contracts in an Organization', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationAccountPage = new OrganizationAccountPage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.goToAllOrgs();

  // Add organization and check new entry
  const organizationName = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationName);

  // Go to Account page and update org
  await organizationSideNavPage.goToAccount();
  await organizationAccountPage.updateOrgToCustomer();
  await organizationAccountPage.addNoteToOrg();

  // Add the first contract to organization and check new entry
  await organizationAccountPage.addContractEmpty();
  await organizationAccountPage.addBillingAddress(0);
  await organizationAccountPage.checkContractsCount(1);
  await organizationAccountPage.addSLIsToContract(0);
  await organizationAccountPage.checkSLIsInAccountPanel();

  // Add the second first contract to organization
  await organizationAccountPage.addContractNonEmpty();
  await organizationAccountPage.checkContractsCount(2);

  // Delete a contract
  await organizationAccountPage.deleteContract(1);
  await organizationAccountPage.checkContractsCount(1);
});

test('CmdK global menu', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationsCmdKPage = new OrganizationsCmdKPage(page);

  await loginPage.login();
  await organizationsPage.goToAllOrgs();

  await organizationsCmdKPage.accessCmdK();
  await organizationsCmdKPage.verifyFinder();
  await organizationsCmdKPage.verifyOrganizationCreation(page, testInfo);
  await organizationsCmdKPage.verifyNavigationToTargets(page);
  await organizationsCmdKPage.verifyNavigationToOpportunities(page);
  await organizationsCmdKPage.verifyNavigationToCustomers(page);
  await organizationsCmdKPage.verifyNavigationToContacts(page);
  await organizationsCmdKPage.verifyNavigationToInvoices(page);
  await organizationsCmdKPage.verifyNavigationToContracts(page);
  await organizationsCmdKPage.verifyNavigationToFlows(page);
  await organizationsCmdKPage.verifyNavigationToSettings(page);
});

test('Assign contact to flow', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const flowsPage = new FlowsPage(page);
  const flowPage = new FlowPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationPeoplePage = new OrganizationPeoplePage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);
  const contactsPage = new ContactsPage(page);

  //
  await loginPage.login();
  await flowsPage.goToFlows();

  const flowName = await flowsPage.addFlow();

  await flowPage.checkNewFlowEntry(flowName);
  await flowPage.goToFlows();
  await flowsPage.checkNewFlowEntry(flowName, flow.create);

  await organizationsPage.goToAllOrgs();

  // Add organization and check new entry
  const organizationName = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationName);

  // Go to People page
  await organizationSideNavPage.goToPeople();

  const contact = await organizationPeoplePage.createContactFromEmpty();

  await organizationSideNavPage.goBack();
  await contactsPage.waitForPageLoad();
  await contactsPage.updateContactFlow(contact, flowName);
  await flowsPage.goToFlows();
  // await flowsPage.checkNewFlowEntry(flowName, flow.update);
});

test('Create opportunities', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const opportunitiesKanbanPage = new OpportunitiesKanbanPage(page);

  await loginPage.login();
  await organizationsPage.goToAllOrgs();

  const organizationName = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  await opportunitiesKanbanPage.goToOpportunitiesKanban();
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    0,
    0,
    0,
    0,
  );
  await opportunitiesKanbanPage.addOpportunity(organizationName);

  const opportunityName = await opportunitiesKanbanPage.updateOpportunityName(
    organizationName,
  );

  await opportunitiesKanbanPage.setOpportunityArrEstimate(opportunityName);
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0,
  );
  await opportunitiesKanbanPage.setWinRates(WinRatesFor.Identified, 10);
  await opportunitiesKanbanPage.setWinRates(WinRatesFor.Qualified, 30);
  await opportunitiesKanbanPage.setWinRates(WinRatesFor.Committed, 55);
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Qualified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    1.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Identified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Committed,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    2.75,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Identified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Won,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Identified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Lost,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Identified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Qualified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    1.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Committed,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    2.75,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Qualified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    1.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Won,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Qualified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    1.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Lost,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Qualified,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    1.5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Committed,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    2.75,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Won,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Committed,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    2.75,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Lost,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Committed,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    2.75,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Won,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    5,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Lost,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    0,
  );

  await opportunitiesKanbanPage.moveOpportunityCard(
    opportunityName,
    KanbanColumns.Won,
  );
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    1,
    1,
    5,
    5,
  );
});
