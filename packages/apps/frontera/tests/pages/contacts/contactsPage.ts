import { Page, expect } from '@playwright/test';

import { assertWithRetry, clickLocatorsThatAreVisible } from '../../helper';

export class ContactsPage {
  private page: Page;

  private sideNavItemAllContacts =
    'div[data-test="side-nav-item-all-contacts"]';
  sideNavItemAllContactsSelected =
    'div[data-test="side-nav-item-all-contacts"] div[aria-selected="true"]';
  private allOrgsSelectAllOrgs = 'button[data-test="all-orgs-select-all-orgs"]';
  private contactsActionsArchive = 'button[data-test="actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';
  private finderTableContacts = 'div[data-test="finder-table-CONTACTS"]';
  private contactNameInContactsTable =
    'p[data-test="contact-name-in-contacts-table"]';
  private flowName = 'div[data-test="flow-name"]';

  private getContactFlowEditSelector(contactId: string) {
    return `button[data-test="contact-flow-edit-${contactId}"]`;
  }

  constructor(page: Page) {
    this.page = page;
  }

  async waitForPageLoad() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllContacts);
  }

  async selectAllContacts(): Promise<boolean> {
    const allContactsSelectAllContacts = this.page.locator(
      this.allOrgsSelectAllOrgs,
    );

    try {
      await allContactsSelectAllContacts.waitFor({
        state: 'visible',
        timeout: 2000,
      });

      const isVisible = await allContactsSelectAllContacts.isVisible();

      if (isVisible) {
        await allContactsSelectAllContacts.click();

        const isAllContactsSelectAllContacts =
          await allContactsSelectAllContacts.getAttribute('data-state');

        return isAllContactsSelectAllContacts === 'checked';
      }
    } catch (error) {
      if (error.name === 'TimeoutError') {
        // Silently return false if the element is not found
        return false;
      }
      // Re-throw any other errors
      throw error;
    }

    return false;
  }

  async archiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.contactsActionsArchive);
  }

  async confirmArchiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.orgActionsConfirmArchive);
  }

  async updateContactFlow(
    contact: { contactId: string; contactName: string },
    flowName: string,
  ) {
    await this.page.reload();
    await assertWithRetry(async () => {
      const organizationIsVisible = await this.page
        .locator(this.getContactFlowEditSelector(contact.contactId))
        .isVisible();

      expect(organizationIsVisible).toBe(true);
    });

    await clickLocatorsThatAreVisible(
      this.page,
      this.getContactFlowEditSelector(contact.contactId),
    );

    const contactFlowInContactsTable = this.page
      .locator(
        `${this.finderTableContacts} ${this.contactNameInContactsTable}:has-text("${contact.contactName}")`,
      )
      .locator('..')
      .locator('..')
      .locator('..')
      .locator(this.flowName);

    await contactFlowInContactsTable.pressSequentially(flowName);
    await contactFlowInContactsTable.press('Enter');

    await expect(
      contactFlowInContactsTable,
      `Expected to have flow ${flowName} allocated to contact ${contact.contactName}`,
    ).toHaveText(flowName);

    await this.page.reload();
  }
}
