import { Page } from '@playwright/test';

export class TargetsPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  sideNavItemTargets = 'button[data-test="side-nav-item-Targets"]';
  sideNavItemTargetsSelected =
    'button[data-test="side-nav-item-Targets"] div[aria-selected="true"]';
}
