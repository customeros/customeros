import { Page } from '@playwright/test';

import { clickLocatorsThatAreVisible } from '../helper';

export class ContactsPage {
  private page: Page;

  private sideNavItemAllContacts =
    'button[data-test="side-nav-item-all-contacts"]';
  private allOrgsSelectAllOrgs = 'button[data-test="all-orgs-select-all-orgs"]';
  private contactsActionsArchive =
    'button[data-test="contacts-actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';

  constructor(page: Page) {
    this.page = page;
  }

  async waitForPageLoad() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllContacts);
  }

  async selectAllOrgs(): Promise<boolean> {
    const allOrgsSelectAllOrgs = this.page.locator(this.allOrgsSelectAllOrgs);

    await allOrgsSelectAllOrgs.waitFor({ state: 'visible', timeout: 2000 });

    const isVisible = await allOrgsSelectAllOrgs.isVisible();

    if (isVisible) {
      await allOrgsSelectAllOrgs.click();

      return true;
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
