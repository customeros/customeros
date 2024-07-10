import { Page } from '@playwright/test';

export class LoginPage {
  constructor(private page: Page) {}

  async login(email: string, password: string) {
    await this.page.goto('https://app.customeros.ai/');

    const googleLoginButtonSelector = 'text=Sign in with Microsoft';
    await this.page.waitForSelector(googleLoginButtonSelector, {
      state: 'visible',
    });
    await this.page.locator(googleLoginButtonSelector).click();

    await this.page.fill(
      'input[aria-label="Enter your email or phone"]',
      email,
    );
    await this.page.click('#idSIButton9');

    await this.page.fill(
      'input[aria-label="Enter the password for silviu@openline.dev"]',
      password,
    );
    await this.page.click('#idSIButton9');

    //   // Complete additional steps
    await this.page.click('#idSubmit_ProofUp_Redirect');
    await this.page.click('a.L0g5CbGcDigQv3yxT1b_:has-text("Skip setup")', {
      timeout: 60000,
    });
    await this.page.click('#upgradeConsentCheckbox');
    await this.page.click('#idSIButton9');
  }
}
