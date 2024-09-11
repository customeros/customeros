import { test } from '@playwright/test';

import { LoginPage } from './pages/loginPage/loginPage';
import { CustomersPage } from './pages/customers/customersPage';
import { OrganizationsPage } from './pages/organizations/organizationsPage';
import { OrganizationAboutPage } from './pages/organization/organizationAboutPage';
import { OrganizationsCmdKPage } from './pages/organizations/organizationsCmdKPage';
import { OrganizationPeoplePage } from './pages/organization/organizationPeoplePage';
import { OrganizationAccountPage } from './pages/organization/organizationAccountPage';
import { OrganizationSideNavPage } from './pages/organization/organizationSideNavPage';
import { OrganizationTimelinePage } from './pages/organization/organizationTimelinePage';

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
  const organizationId = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  await organizationsPage.checkNewOrganizationEntry(organizationId);

  // Go to Customers page and ensure no new org
  await organizationsPage.goToCustomersPage();
  await customersPage.ensureNumberOfCustomersExist(0);

  // Go back to All Orgs page
  await organizationsPage.goToAllOrgsPage();

  // Make the organization a customer
  await organizationsPage.updateOrgToCustomer(organizationId);

  // Go to Customers page and ensure we have a new customer
  await organizationsPage.goToCustomersPage();
  await customersPage.ensureNumberOfCustomersExist(1);
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
  const organizationId = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationId);

  // Go to People page
  await organizationSideNavPage.goToAbout();
  await organizationAboutPage.populateAboutFields();
  await organizationAboutPage.checkPopulatedAboutFields(
    organizationId,
    'customeros.fe.testing',
  );
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
  const organizationId = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationId);

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
  const organizationId = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationId);

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
  const organizationId = await organizationsPage.addNonInitialOrganization(
    testInfo,
  );

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await organizationsPage.goToOrganization(organizationId);

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
  await organizationsCmdKPage.verifyNavigationToSettings(page);
  await organizationsCmdKPage.verifyNavigationToCustomerMap(page);
});
