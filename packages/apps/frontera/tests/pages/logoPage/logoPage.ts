import { Page } from '@playwright/test';

import { SettingsPage } from '../settings/settingsPage';
import {
  ensureLocatorIsVisible,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class LogoPage {
  constructor(private page: Page) {
    this.page = page;
  }

  private settingsButton = 'button[data-test="settings-button"]';
  // private sideNavItemLogoSettings = 'div[data-test="logo-settings"]';

  async goToSettings() {
    const settingsPage = new SettingsPage(this.page);

    await clickLocatorsThatAreVisible(this.page, this.settingsButton);
    // await clickLocatorsThatAreVisible(this.page, this.sideNavItemLogoSettings);
    await ensureLocatorIsVisible(this.page, settingsPage.settingsHeader);
  }
}
