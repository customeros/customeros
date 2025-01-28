import { Page, expect } from '@playwright/test';

import {
  assertWithRetry,
  clickLocatorThatIsVisible,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class ContactsPage {
  private page: Page;

  sideNavItemAllContacts = 'div[data-test="side-nav-item-all-contacts"]';
  sideNavItemAllContactsSelected =
    'div[data-test="side-nav-item-all-contacts"] div[aria-selected="true"]';
  private allOrgsSelectAllOrgs = 'div[data-test="all-orgs-select-all-orgs"]';
  private contactsActionsArchive = 'button[data-test="actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';
  finderTableContacts = 'div[data-test="finder-table-CONTACTS"]';
  private addContactToFlowConfirmation =
    'button[data-test="add-contact-to-flow-confirmation"]';
  private contactNameInContactsTable =
    'p[data-test="contact-name-in-contacts-table"]';
  private contactCurrentFlowsInContactsTable =
    'div[data-test="contact-current-flows-in-contacts-table"]';
  private contactsFlowFoundEntry = 'div[data-test="contacts-flow-found-entry"]';

  private flowName = 'div[data-test="flow-name"]';

  private getContactFlowEditSelector(contactName: string) {
    return `[data-test="contact-name-in-contacts-table"]:has-text("${contactName}")`;
  }

  private getContactFlowSelector(contactName: string) {
    return `[data-index]:has(${this.contactNameInContactsTable}:text("${contactName}")) ${this.contactCurrentFlowsInContactsTable}, [data-index]:has(${this.contactNameInContactsTable}:text("${contactName}")) ${this.flowName}`;
  }

  constructor(page: Page) {
    this.page = page;
  }

  async waitForPageLoad() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllContacts);
  }

  async selectAllContacts() {
    const allContactsSelectAllContacts = this.page.locator(
      this.allOrgsSelectAllOrgs,
    );

    await allContactsSelectAllContacts.click();
  }

  async archiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.contactsActionsArchive);
  }

  async confirmArchiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.orgActionsConfirmArchive);
  }

  async updateContactFlow(contactName: string, flowName: string) {
    await this.page.reload();
    await assertWithRetry(async () => {
      const organizationIsVisible = await this.page
        .locator(this.getContactFlowEditSelector(contactName))
        .isVisible();

      expect(organizationIsVisible).toBe(true);
    });

    await this.page.waitForTimeout(5000);

    const contactFlowInContactsTable = await clickLocatorThatIsVisible(
      this.page,
      this.getContactFlowSelector(contactName),
    );

    await this.page.keyboard.press('Escape');

    await clickLocatorThatIsVisible(
      this.page,
      this.getContactFlowSelector(contactName),
    );

    const specificcCntactsFlowFoundEntry = `${this.contactsFlowFoundEntry}[data-value="${flowName}"]`;

    await clickLocatorThatIsVisible(this.page, specificcCntactsFlowFoundEntry);
    await clickLocatorThatIsVisible(
      this.page,
      this.addContactToFlowConfirmation,
    );

    await expect(
      contactFlowInContactsTable,
      `Expected to have flow ${flowName} allocated to contact ${contactName}`,
    ).toHaveText(flowName);

    await this.page.reload();
  }
}
