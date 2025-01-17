import { Page } from '@playwright/test';

export class SettingsPage {
  constructor(private page: Page) {
    this.page = page;
  }

  settingsGoBack = 'button[data-test="settings-go-back"]';
}
