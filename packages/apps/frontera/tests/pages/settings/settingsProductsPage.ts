import { Page } from '@playwright/test';

export class SettingsProductsPage {
  constructor(page) {
    this.page = page;
  }

  private page: Page;
  settingsProducts = 'button[data-test="sideNav-settings-products"]';
}
