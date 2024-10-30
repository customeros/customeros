import { Page, expect, TestInfo } from '@playwright/test';

import { FlowsPage } from '../flows/flowsPage';
import { TargetsPage } from '../targets/targetsPage';
import { sideNavSelectors } from '../sideNavSelectors';
import { SettingsPage } from '../settings/settingsPage';
import { InvoicesPage } from '../invoices/invoicesPage';
import { ContactsPage } from '../contacts/contactsPage';
import { OrganizationsPage } from './organizationsPage';
import { ContractsPage } from '../contracts/contractsPage';
// import { CustomersPage } from '../customers/customersPage';
import { CustomerMapPage } from '../customer-map/customerMapPage';
import { SettingsAccountsPage } from '../settings/settingsAccounts';
import {
  ensureLocatorIsVisible,
  clickLocatorThatIsVisible,
} from '../../helper';
import { OpportunitiesKanbanPage } from '../opportunitiesKanban/opportunitiesKanbanPage';

export class OrganizationsCmdKPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private sideNavItemCustomersSelected =
    sideNavSelectors.sideNavItemCustomersSelected;
  private organizationsHub = 'div[data-test="organization-hub"]';
  private organizationsHubSpan = 'div[data-test="organization-hub"] span';
  private organizationHubInput = 'input[data-test="organization-hub-input"]';
  // private organizationHubNavigate =
  //   'div[data-test="organization-hub-navigate"]';
  private organizationHubAddNewOrgs =
    'div[data-test="organization-hub-add-new-orgs"]';
  private organizationHubGt = 'div[data-test="organization-hub-gt"]';
  private organizationHubGo = 'div[data-test="organization-hub-go"]';
  private organizationHubGc = 'div[data-test="organization-hub-gc"]';
  private organizationHubGz = 'div[data-test="organization-hub-gz"]';
  private organizationHubGn = 'div[data-test="organization-hub-gn"]';
  private organizationHubGi = 'div[data-test="organization-hub-gi"]';
  private organizationHubGr = 'div[data-test="organization-hub-gr"]';
  private organizationHubGf = 'div[data-test="organization-hub-gf"]';
  private organizationHubGs = 'div[data-test="organization-hub-gs"]';
  private organizationHubGd = 'div[data-test="organization-hub-gd"]';

  private async openCmdK() {
    await this.page.waitForSelector('div[data-test="search-orgs"]', {
      state: 'visible',
    });

    await this.page.keyboard.down('Meta');
    await this.page.keyboard.press('KeyK');
    await this.page.keyboard.up('Meta');

    await this.page.waitForSelector(this.organizationsHub, {
      state: 'visible',
    });
  }

  private async verifyNavigationWithClick(
    organizationHubNavigationDestination: string,
    sideNavItemSelected: string,
  ) {
    await clickLocatorThatIsVisible(
      this.page,
      organizationHubNavigationDestination,
    );

    const sideNavItemSelectedVisible = await ensureLocatorIsVisible(
      this.page,
      sideNavItemSelected,
    );

    const sideNavItemSelectedTextContent =
      await sideNavItemSelectedVisible.getAttribute('aria-selected');

    expect(sideNavItemSelectedTextContent).toBe('true');
  }

  private async verifyNavigationWithKeyboard(
    secondKey: string,
    sideNavItemSelected: string,
  ) {
    await this.page.keyboard.press('KeyG');
    await this.page.keyboard.press(secondKey);

    const sideNavItemSelectedVisible = await ensureLocatorIsVisible(
      this.page,
      sideNavItemSelected,
    );

    const sideNavItemSelectedTextContent =
      await sideNavItemSelectedVisible.getAttribute('aria-selected');

    expect(sideNavItemSelectedTextContent).toBe('true');
  }

  async accessCmdK() {
    await this.openCmdK();

    const organizationsHub = this.page
      .locator(this.organizationsHubSpan)
      .first();

    const organizationsHubText = await organizationsHub.textContent();

    const organizationHubInput = this.page
      .locator(this.organizationHubInput)
      .first();

    const organizationHubInputText = await organizationHubInput.getAttribute(
      'placeholder',
    );

    const organizationHubAddNewOrgs = this.page
      .locator(this.organizationHubAddNewOrgs)
      .first();

    const organizationHubAddNewOrgsText =
      await organizationHubAddNewOrgs.textContent();

    const navigationGroup = this.page
      .locator(this.organizationsHub)
      .locator('div[data-test="organization-hub-navigate"] div[role="group"]')
      .first();

    const navigationItems = await navigationGroup
      .locator('div[cmdk-item]')
      .all();

    const organizationHubGt = this.page.locator(this.organizationHubGt);
    const organizationHubGtText = await organizationHubGt.textContent();
    const navigationItemTextOne = await navigationItems[0].textContent();

    const organizationHubGo = this.page.locator(this.organizationHubGo);
    const organizationHubGoText = await organizationHubGo.textContent();
    const navigationItemTextTwo = await navigationItems[1].textContent();

    const organizationHubGc = this.page.locator(this.organizationHubGc);
    const organizationHubGcText = await organizationHubGc.textContent();
    const navigationItemTextThree = await navigationItems[2].textContent();

    const organizationHubGz = this.page.locator(this.organizationHubGz);
    const organizationHubGzText = await organizationHubGz.textContent();
    const navigationItemTextFour = await navigationItems[3].textContent();

    const organizationHubGn = this.page.locator(this.organizationHubGn);
    const organizationHubGnText = await organizationHubGn.textContent();
    const navigationItemTextFive = await navigationItems[4].textContent();

    const organizationHubGi = this.page.locator(this.organizationHubGi);
    const organizationHubGiText = await organizationHubGi.textContent();
    const navigationItemTextSix = await navigationItems[5].textContent();

    const organizationHubGr = this.page.locator(this.organizationHubGr);
    const organizationHubGrText = await organizationHubGr.textContent();
    const navigationItemTextSeven = await navigationItems[6].textContent();

    const organizationHubGq = this.page.locator(this.organizationHubGf);
    const organizationHubGqText = await organizationHubGq.textContent();
    const navigationItemTextEight = await navigationItems[7].textContent();

    const organizationHubGs = this.page.locator(this.organizationHubGs);
    const organizationHubGsText = await organizationHubGs.textContent();
    const navigationItemTextNine = await navigationItems[8].textContent();

    const organizationHubGd = this.page.locator(this.organizationHubGd);
    const organizationHubGdText = await organizationHubGd.textContent();
    const navigationItemTextTen = await navigationItems[9].textContent();

    await Promise.all([
      expect.soft(organizationsHubText.trim()).toBe('Organizations'),
      expect
        .soft(organizationHubInputText.trim())
        .toBe('Type a command or search'),
      expect
        .soft(organizationHubAddNewOrgsText.trim())
        .toBe('Add new organizations...'),
      expect.soft(navigationItems).toHaveLength(10),
      expect
        .soft(organizationHubGtText.trim())
        .toBe(navigationItemTextOne.trim()),
      expect
        .soft(organizationHubGoText.trim())
        .toBe(navigationItemTextTwo.trim()),
      expect
        .soft(organizationHubGcText.trim())
        .toBe(navigationItemTextThree.trim()),
      expect
        .soft(organizationHubGzText.trim())
        .toBe(navigationItemTextFour.trim()),
      expect
        .soft(organizationHubGnText.trim())
        .toBe(navigationItemTextFive.trim()),
      expect
        .soft(organizationHubGiText.trim())
        .toBe(navigationItemTextSix.trim()),
      expect
        .soft(organizationHubGrText.trim())
        .toBe(navigationItemTextSeven.trim()),
      expect
        .soft(organizationHubGqText.trim())
        .toBe(navigationItemTextEight.trim()),
      expect
        .soft(organizationHubGsText.trim())
        .toBe(navigationItemTextNine.trim()),
      expect
        .soft(organizationHubGdText.trim())
        .toBe(navigationItemTextTen.trim()),
    ]);

    await this.page.keyboard.press('Escape');

    const organizationsHubCount = await this.page
      .locator(this.organizationsHubSpan)
      .count();

    expect(organizationsHubCount).toBe(0);
  }

  async verifyFinder() {
    await this.openCmdK();

    await this.page
      .locator(this.organizationHubInput)
      .pressSequentially('go to customer');

    const navigationGroup = this.page
      .locator(this.organizationsHub)
      .locator('div[data-test="organization-hub-navigate"] div[role="group"]')
      .first();

    const navigationItems = await navigationGroup
      .locator('div[cmdk-item]')
      .all();

    const organizationHubGc = this.page.locator(this.organizationHubGc);
    const organizationHubGcText = await organizationHubGc.textContent();
    const navigationItemTextOne = await navigationItems[0].textContent();

    await Promise.all([
      expect.soft(navigationItems).toHaveLength(2),
      expect
        .soft(organizationHubGcText.trim())
        .toBe(navigationItemTextOne.trim()),
    ]);

    await this.page.keyboard.press('Escape');
  }

  async verifyOrganizationCreation(page: Page, testInfo: TestInfo) {
    const organizationsPage = new OrganizationsPage(page);

    await this.openCmdK();
    await organizationsPage.addOrganization(
      this.organizationHubAddNewOrgs,
      testInfo,
    );
    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToTargets(page: Page) {
    const targetsPage = new TargetsPage(page);
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyT',
      targetsPage.sideNavItemTargetsSelected,
    );

    await this.page.goBack();
    await this.page.waitForLoadState('load');

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGt,
      targetsPage.sideNavItemTargetsSelected,
    );

    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToOpportunities(page: Page) {
    const opportunitiesPage = new OpportunitiesKanbanPage(page);
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyO',
      opportunitiesPage.sideNavItemOpportunitiesSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGo,
      opportunitiesPage.sideNavItemOpportunitiesSelected,
    );

    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToCustomers(page: Page) {
    // const customersPage = new CustomersPage(page);
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyC',
      this.sideNavItemCustomersSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGc,
      this.sideNavItemCustomersSelected,
    );

    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToContacts(page: Page) {
    const contactsPage = new ContactsPage(page);
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyN',
      contactsPage.sideNavItemAllContactsSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGn,
      contactsPage.sideNavItemAllContactsSelected,
    );

    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToInvoices(page: Page) {
    const invoicesPage = new InvoicesPage();
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyI',
      invoicesPage.sideNavItemAllUpcomingSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGi,
      invoicesPage.sideNavItemAllUpcomingSelected,
    );

    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToContracts(page: Page) {
    const contractsPage = new ContractsPage();
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyR',
      contractsPage.sideNavItemAllContractsSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGr,
      contractsPage.sideNavItemAllContractsSelected,
    );

    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToFlows(page: Page) {
    const flowsPage = new FlowsPage(page);
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyF',
      flowsPage.sideNavItemAllFlowsSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGf,
      flowsPage.sideNavItemAllFlowsSelected,
    );

    await organizationsPage.goToAllOrgs();
  }

  async verifyNavigationToSettings(page: Page) {
    const settingsAccountsPage = new SettingsAccountsPage(page);
    const settingsPage = new SettingsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyS',
      settingsAccountsPage.settingsAccountsSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGs,
      settingsAccountsPage.settingsAccountsSelected,
    );

    await clickLocatorThatIsVisible(this.page, settingsPage.settingsGoBack);
  }

  async verifyNavigationToCustomerMap(page: Page) {
    const customerMapPage = new CustomerMapPage();
    const organizationsPage = new OrganizationsPage(page);

    await this.verifyNavigationWithKeyboard(
      'KeyD',
      customerMapPage.sideNavItemAllCustomerMapSelected,
    );

    await this.page.goBack();

    await this.openCmdK();
    await this.verifyNavigationWithClick(
      this.organizationHubGd,
      customerMapPage.sideNavItemAllCustomerMapSelected,
    );

    await organizationsPage.goToAllOrgs();
  }
}
