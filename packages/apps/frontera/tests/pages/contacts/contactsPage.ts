import { Page } from '@playwright/test';

import { clickLocatorsThatAreVisible } from '../../helper';

export class ContactsPage {
  private page: Page;

  private sideNavItemAllContacts =
    'button[data-test="side-nav-item-all-contacts"]';
  sideNavItemAllContactsSelected =
    'button[data-test="side-nav-item-all-contacts"] div[aria-selected="true"]';
  private allOrgsSelectAllOrgs = 'button[data-test="all-orgs-select-all-orgs"]';
  private contactsActionsArchive = 'button[data-test="actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';

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

        return true;
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
}
