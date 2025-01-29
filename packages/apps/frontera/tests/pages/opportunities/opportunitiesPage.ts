import { Page } from '@playwright/test';

import { clickLocatorsThatAreVisible } from '../../helper';

export class OpportunitiesPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private sideNavItemOpportunities =
    'div[data-test="side-nav-item-opportunities"]';
  finderTableOpportunities = 'div[data-test="finder-table-OPPORTUNITIES"]';
  private allOrgsSelectAllOrgs = 'div[data-test="all-orgs-select-all-orgs"]';
  private opportunitiesActionsArchive = 'button[data-test="actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';
  private prospectsListButton = 'button[data-test="prospects-list-button"]';

  async goToOpportunitiesList() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.sideNavItemOpportunities,
      this.prospectsListButton,
    );
  }

  async selectAllOpportunities() {
    const opportunitiesSelectAllOpportunities = this.page.locator(
      this.allOrgsSelectAllOrgs,
    );

    await opportunitiesSelectAllOpportunities.click();
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
