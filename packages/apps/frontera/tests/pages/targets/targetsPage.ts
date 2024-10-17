import { Page } from '@playwright/test';

export class TargetsPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  sideNavItemTargets = 'div[data-test="side-nav-item-Targets"]';
  sideNavItemTargetsSelected =
    'div[data-test="side-nav-item-Targets"] div[aria-selected="true"]';
}
