// import { test } from '@playwright/test';
import { test } from './videoFixture';
import { companies } from './test-data';
// import { FlowPage } from './pages/flows/flowPage';
// import { ContactsPage } from './pages/contacts/contactsPage';
// import { FlowsPage } from './pages/flows/flowsPage';
import { LogoPage } from './pages/logoPage/logoPage';
import { LoginPage } from './pages/loginPage/loginPage';
import { SettingsPage } from './pages/settings/settingsPage';
import { MailboxesPage } from './pages/settings/mailboxesPage';
// import { ContactsPage } from './pages/contacts/contactsPage';
import { CustomersPage } from './pages/customers/customersPage';
import { CompaniesPage } from './pages/companies/companiesPage';
import { WinRatesFor } from './pages/opportunitiesKanban/winRates';
import { KanbanColumns } from './pages/opportunitiesKanban/columns';
import { CompanyCmdKPage } from './pages/companies/companyCmdKPage';
import { CompanyAboutPage } from './pages/company/companyAboutPage';
import { CompanyPeoplePage } from './pages/company/companyPeoplePage';
import { CompanySideNavPage } from './pages/company/companySideNavPage';
import { CompanyAccountPage } from './pages/company/companyAccountPage';
import { CompanyTimelinePage } from './pages/company/companyTimelinePage';
import { SettingsProductsPage } from './pages/settings/settingsProductsPage';
import { SkuType } from '../src/routes/src/types/__generated__/graphql.types';
import { OpportunitiesKanbanPage } from './pages/opportunitiesKanban/opportunitiesKanbanPage';

test.setTimeout(300000);

test('Convert a Company to Customer', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const companiesPage = new CompaniesPage(page);
  const customersPage = new CustomersPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load Companies page
  await companiesPage.goToCompanies();

  // Add company and check new entry
  const companyName = await companiesPage.addNonInitialCompany(testInfo);

  await companiesPage.checkNewCompanyEntry(companyName);

  // Go to Customers page and ensure no new company
  await companiesPage.goToCustomersPage();
  // await customersPage.ensureNumberOfCustomersExist(0);
  await customersPage.ensureCustomerExists(companyName, false);

  // Go back to Companies page
  await companiesPage.goToCompaniesPage();

  // Make the company a customer
  await companiesPage.updateCompanyToCustomer(companyName);

  // Go to Customers page and ensure we have a new customer
  await companiesPage.goToCustomersPage();
  // await customersPage.ensureNumberOfCustomersExist(1);
  await customersPage.ensureCustomerExists(companyName, true);
});

test('Add About information to a Company', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const companiesPage = new CompaniesPage(page);
  const companyAboutPage = new CompanyAboutPage(page);
  const companySideNavPage = new CompanySideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load Companies page
  await companiesPage.goToCompanies();

  // Add company and check new entry
  const companyName = await companiesPage.addNonInitialCompany(testInfo);

  //Access newly created company
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await companiesPage.goToCompany(companyName);

  await companySideNavPage.goToAbout();

  //Check enrichment
  await companyAboutPage.enrichCompany(companies.create.domain);
  companies.create.name = companyName;
  await companyAboutPage.checkEnrichedAboutFields(companies.create);

  //Check updates that override the enrichment
  // await companyAboutPage.populateAboutFields(companies.update);
  // await companyAboutPage.checkPopulatedAboutFields(companies.update);
});

test('Create People entry in Company', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const companiesPage = new CompaniesPage(page);
  const companyPeoplePage = new CompanyPeoplePage(page);
  const companySideNavPage = new CompanySideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load Companies page
  await companiesPage.goToCompanies();

  // Add company and check new entry
  const companyName = await companiesPage.addNonInitialCompany(testInfo);

  //Access newly created company
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await companiesPage.goToCompany(companyName);

  // Go to People page
  await companySideNavPage.goToPeople();
  await companyPeoplePage.createContactFromEmpty();
});

test('Create Timeline entries in an Company', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const companiesPage = new CompaniesPage(page);
  const companySideNavPage = new CompanySideNavPage(page);
  const companyTimelinePage = new CompanyTimelinePage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load Companies page
  await companiesPage.goToCompanies();

  // Add company and check new entry
  const companyName = await companiesPage.addNonInitialCompany(testInfo);

  //Access newly created company
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await companiesPage.goToCompany(companyName);

  // Go to Account page and update company
  await companySideNavPage.goToAccount();
  await companyTimelinePage.ensureEmailPermissionPromptIsRedirecting();
  // await page.goBack();
  // await companyTimelinePage.ensureLogEntryCanBeAdded();
  // await companyTimelinePage.ensureReminderCanBeAdded();
});

test('Create Contracts in an Company', async ({ page }, testInfo) => {
  const logoPage = new LogoPage(page);
  const loginPage = new LoginPage(page);
  const settingsPage = new SettingsPage(page);
  const settingsProductsPage = new SettingsProductsPage(page);
  const companiesPage = new CompaniesPage(page);
  const companyAccountPage = new CompanyAccountPage(page);
  const companySideNavPage = new CompanySideNavPage(page);

  // Login
  await loginPage.login();
  await logoPage.goToSettings();
  await settingsProductsPage.goToSettingsProductsPage();
  await settingsProductsPage.addProduct(SkuType.Subscription, 's1', '123.01');
  await settingsProductsPage.addProduct(SkuType.OneTime, 'o1', '111.99');

  await settingsPage.goBack();
  // Wait for redirect and load Companies page
  await companiesPage.goToCompanies();

  // Add company and check new entry
  const companyName = await companiesPage.addNonInitialCompany(testInfo);

  //Access newly created company
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await companiesPage.goToCompany(companyName);

  // Go to Account page and update company
  await companySideNavPage.goToAccount();
  await companyAccountPage.updateCompanyToCustomer();

  // Add the first contract to company and check new entry
  await companyAccountPage.addContractEmpty();
  await companyAccountPage.addBillingAddress(0);
  await companyAccountPage.checkContractsCount(1);
  await companyAccountPage.addSLIsToContract(0);
  await companyAccountPage.checkSLIsInAccountPanel();

  // Add the second first contract to company
  await companyAccountPage.addContractNonEmpty();
  await companyAccountPage.checkContractsCount(2);

  // Delete a contract
  await companyAccountPage.deleteContract(1);
  await companyAccountPage.checkContractsCount(1);
});

test('Create Note in an Company', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const companiesPage = new CompaniesPage(page);
  const companyAccountPage = new CompanyAccountPage(page);
  const companySideNavPage = new CompanySideNavPage(page);

  // Login
  await loginPage.login();
  // Wait for redirect and load Companies page
  await companiesPage.goToCompanies();

  // Add company and check new entry
  const companyName = await companiesPage.addNonInitialCompany(testInfo);

  //Access newly created company
  await new Promise((resolve) => setTimeout(resolve, 1500));
  await companiesPage.goToCompany(companyName);

  // Go to Account page and update company
  await companySideNavPage.goToAccount();
  await companyAccountPage.addNoteToCompany();
});

test('CmdK global menu', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const companiesPage = new CompaniesPage(page);
  const companiesCmdKPage = new CompanyCmdKPage(page);

  await loginPage.login();
  await companiesPage.goToCompanies();

  await companiesCmdKPage.accessCmdK();
  await companiesCmdKPage.verifyFinder();
  await companiesCmdKPage.verifyCompanyCreation(page, testInfo);
  await companiesCmdKPage.verifyNavigationToOpportunities(page);
  await companiesCmdKPage.verifyNavigationToCompanies(page);
  await companiesCmdKPage.verifyNavigationToContacts(page);
  await companiesCmdKPage.verifyNavigationToInvoices(page);
  await companiesCmdKPage.verifyNavigationToContracts(page);
  // await companiesCmdKPage.verifyNavigationToFlows(page);
  await companiesCmdKPage.verifyNavigationToSettings(page);
});

// test('Assign contact to flow', async ({ page }, testInfo) => {
//   const loginPage = new LoginPage(page);
//   const flowsPage = new FlowsPage(page);
//   const flowPage = new FlowPage(page);
//   const companiesPage = new companiesPage(page);
//   const companyPeoplePage = new companyPeoplePage(page);
//   const companySideNavPage = new companySideNavPage(page);
//   const contactsPage = new ContactsPage(page);
//
//   //
//   await loginPage.login();
//   await flowsPage.goToFlows();
//
//   const flowName = await flowsPage.addFlow();
//
//   await flowPage.checkNewFlowEntry(flowName);
//   // await flowPage.goToFlows();
//   await flowsPage.checkNewFlowEntry(flowName, flow.create);
//
//   await companiesPage.goToAllCompanies();
//
//   // Add company and check new entry
//   const companyName = await companiesPage.addNonInitialCompany(
//     testInfo,
//   );
//
//   //Access newly created company
//   await new Promise((resolve) => setTimeout(resolve, 1500));
//   await companiesPage.goToCompany(companyName);
//
//   // Go to People page
//   await companySideNavPage.goToPeople();
//
//   const contactName = await companyPeoplePage.createContactFromEmpty();
//
//   await companySideNavPage.goBack();
//   await contactsPage.waitForPageLoad();
//   await contactsPage.updateContactFlow(contactName, flowName);
//   await flowsPage.goToFlows();
//   await flowsPage.checkNewFlowEntry(flowName, flow.update);
// });

test('Create opportunities', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const companiesPage = new CompaniesPage(page);
  const opportunitiesKanbanPage = new OpportunitiesKanbanPage(page);

  await loginPage.login();
  await companiesPage.goToCompanies();

  const companyName = await companiesPage.addNonInitialCompany(testInfo);

  await opportunitiesKanbanPage.goToOpportunitiesKanban();
  await opportunitiesKanbanPage.checkOpportunitiesKanbanHeaderValues(
    0,
    0,
    0,
    0,
  );
  await opportunitiesKanbanPage.addOpportunity(companyName);

  const opportunityName = await opportunitiesKanbanPage.updateOpportunityName(
    companyName,
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

test('Purchase mailboxes', async ({ page }, testInfo) => {
  const loginPage = new LoginPage(page);
  const logoPage = new LogoPage(page);
  const mailboxesPage = new MailboxesPage(page);

  await loginPage.login();
  await logoPage.goToSettings();
  await mailboxesPage.goToMailboxes();
  await mailboxesPage.setupMailboxes();
  await mailboxesPage.searchForDomains();
  await mailboxesPage.addDomainsToCart();
  await mailboxesPage.setRedirectUrl();
  await mailboxesPage.setUsernames();
  await mailboxesPage.checkout();
  await mailboxesPage.fillInPaymentForm();
});
