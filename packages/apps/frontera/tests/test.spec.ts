import { test } from './hooks';
import { LoginPage } from './pages/loginPage';
import { CustomersPage } from './pages/customersPage';
import { OrganizationsPage } from './pages/organizationsPage';
import { OrganizationPeoplePage } from './pages/organization/organizationPeoplePage';
import { OrganizationAccountPage } from './pages/organization/organizationAccountPage';
import { OrganizationSideNavPage } from './pages/organization/organizationSideNavPage';
import { OrganizationTimelinePage } from './pages/organization/organizationTimelinePage';

test.setTimeout(180000);

test('convert org to customer', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const customersPage = new CustomersPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.waitForPageLoad();

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

test('create people in organization', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationPeoplePage = new OrganizationPeoplePage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.waitForPageLoad();

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

test('create timeline entries in organization', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);
  const organizationTimelinePage = new OrganizationTimelinePage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.waitForPageLoad();

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

test('create contracts in organization', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const organizationsPage = new OrganizationsPage(page);
  const organizationAccountPage = new OrganizationAccountPage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await organizationsPage.waitForPageLoad();

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
