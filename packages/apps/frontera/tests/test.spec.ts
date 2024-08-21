import { test } from './hooks';
import { LoginPage } from './pages/loginPage';
import { CustomersPage } from './pages/customersPage';
import { AddressBookPage } from './pages/addressBookPage';
import { OrganizationPeoplePage } from './pages/organization/organizationPeoplePage';
import { OrganizationAccountPage } from './pages/organization/organizationAccountPage';
import { OrganizationSideNavPage } from './pages/organization/organizationSideNavPage';
import { OrganizationTimelinePage } from './pages/organization/organizationTimelinePage';

test.setTimeout(180000);

test('convert org to customer', async ({ page }) => {
  const loginPage = new LoginPage(page);
  const addressBookPage = new AddressBookPage(page);
  const customersPage = new CustomersPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await addressBookPage.waitForPageLoad();

  // Add organization and check new entry
  await addressBookPage.addOrganization();
  await addressBookPage.checkNewEntry();

  // Go to Customers page and ensure no new org
  await addressBookPage.goToCustomersPage();
  await customersPage.ensureNumberOfCustomersExist(0);

  // Go back to All Orgs page
  await addressBookPage.goToAllOrgsPage();

  // Make the organization a customer
  await addressBookPage.updateOrgToCustomer();

  // Go to Customers page and ensure we have a new customer
  await addressBookPage.goToCustomersPage();
  await customersPage.ensureNumberOfCustomersExist(1);
});

test('create and delete contracts', async ({ page }) => {
  const loginPage = new LoginPage(page);
  const addressBookPage = new AddressBookPage(page);
  const organizationPeoplePage = new OrganizationPeoplePage(page);
  const organizationAccountPage = new OrganizationAccountPage(page);
  const organizationSideNavPage = new OrganizationSideNavPage(page);
  const organizationTimelinePage = new OrganizationTimelinePage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load All Orgs page
  await addressBookPage.waitForPageLoad();

  // Add organization and check new entry
  await addressBookPage.addOrganization();

  //Access newly created organization
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await addressBookPage.goToOrganization();

  // Go to People page
  await organizationSideNavPage.goToPeople();
  await organizationPeoplePage.createContactFromEmpty();

  // Go to Account page and update org
  await organizationSideNavPage.goToAccount();
  await organizationAccountPage.updateOrgToCustomer();
  await organizationAccountPage.addNoteToOrg();
  await organizationTimelinePage.ensureEmailPermissionPromptIsRedirecting();
  await page.goBack();
  await organizationTimelinePage.ensureLogEntryCanBeAdded();
  await organizationTimelinePage.ensureReminderCanBeAdded();

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
