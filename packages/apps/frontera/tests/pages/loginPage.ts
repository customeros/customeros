import { Page, expect } from '@playwright/test';

// import config from '../config';

export class LoginPage {
  constructor(private page: Page) {}

  async login() {
    // const loginUrl = config.LOCAL_LOGIN_URL;
    const loginUrl = process.env.PROD_FE_TEST_USER_URL;

    // Listen to all responses from the page
    this.page.on('response', async (response) => {
      const url = response.url();
      const status = response.status();

      if (
        url.includes('https://mid.customeros.ai/customer-os-api') &&
        status === 401
      ) {
        expect(
          status,
          'Test failed: Expected not to receive a 401 Unauthorized response from the API endpoint',
        ).not.toBe(401);
      }
    });

    await this.page.goto(loginUrl);
  }
}
