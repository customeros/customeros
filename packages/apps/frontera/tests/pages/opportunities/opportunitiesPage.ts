import { Page } from '@playwright/test';

import { clickLocatorsThatAreVisible } from '../../helper';

export class OpportunitiesPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private sideNavItemOpportunities =
    'button[data-test="side-nav-item-opportunities"]';
  private allOrgsSelectAllOrgs = 'button[data-test="all-orgs-select-all-orgs"]';
  private opportunitiesActionsArchive = 'button[data-test="actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';

  async goToOpportunities() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemOpportunities);
  }

  async selectAllOpportunities() {
    const opportunitiesSelectAllOpportunities = this.page.locator(
      this.allOrgsSelectAllOrgs,
    );

    try {
      await opportunitiesSelectAllOpportunities.waitFor({
        state: 'visible',
        timeout: 2000,
      });

      const isVisible = await opportunitiesSelectAllOpportunities.isVisible();

      if (isVisible) {
        await opportunitiesSelectAllOpportunities.click();

        // Wait for a short time to allow for any asynchronous updates
        await this.page.waitForTimeout(100);

        // Check if the button is checked after clicking
        return (
          (await opportunitiesSelectAllOpportunities.getAttribute(
            'aria-checked',
          )) === 'true' ||
          (await opportunitiesSelectAllOpportunities.getAttribute(
            'data-state',
          )) === 'checked'
        );
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
    await clickLocatorsThatAreVisible(
      this.page,
      this.opportunitiesActionsArchive,
    );
  }

  async confirmArchiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.orgActionsConfirmArchive);
  }
}
