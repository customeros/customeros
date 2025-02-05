import { Page } from '@playwright/test';

import { clickLocatorThatIsVisible } from '../../helper';

export class SettingsPage {
  constructor(private page: Page) {
    this.page = page;
  }

  settingsGoBack = 'button[data-test="settings-go-back"]';
  settingsHeader = 'p[data-test="settings-header"]';

  async goBack() {
    await clickLocatorThatIsVisible(this.page, this.settingsGoBack);
  }
}
