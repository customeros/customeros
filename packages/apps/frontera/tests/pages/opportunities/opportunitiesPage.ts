import { Page } from '@playwright/test';

export class OpportunitiesPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  sideNavItemOpportunities = 'button[data-test="side-nav-item-Opportunities"]';
  sideNavItemOpportunitiesSelected =
    'button[data-test="side-nav-item-Opportunities"] div[aria-selected="true"]';
}
