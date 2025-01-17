import { Page } from '@playwright/test';

export class LogoPage {
  constructor(private page: Page) {
    this.page = page;
  }

  private sideNavItemLogo = 'button[data-test="logo-button"]';
}
