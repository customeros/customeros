import { Page } from '@playwright/test';

import { ensureLocatorIsVisible } from '../../helper';

export class SettingsAccountsPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;
  private static settingsAccountsHeader =
    'h1[data-test="settings-accounts-header"]';
  settingsAccountsSelected =
    'button[data-test="settings-accounts"] div[aria-selected="true"]';

  static async ensureSettingsAccountsHeaderIsVisible(page: Page) {
    await ensureLocatorIsVisible(page, this.settingsAccountsHeader);
  }
}
