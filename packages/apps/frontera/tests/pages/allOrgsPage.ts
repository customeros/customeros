import { Page, expect } from '@playwright/test';

export class AllOrgsPage {
  constructor(private page: Page) {}

  async waitForPageLoad() {
    const allOrgsMenubutton = this.page.locator(
      'button[data-test="side-nav-item-all-orgs"]',
    );
    await allOrgsMenubutton.waitFor({ state: 'visible' });
    await expect(allOrgsMenubutton).toBeVisible();
    await allOrgsMenubutton.click();
  }

  async addOrganization() {
    const addOrganizationButton = this.page.locator(
      'button[data-test="all-orgs-add-org"]',
    );

    await addOrganizationButton.click();
    await this.page.waitForSelector('[data-index="0"]', { timeout: 30000 });
  }

  async assertWithRetry(assertionFunc, maxRetries = 5, retryInterval = 3000) {
    let lastError;
    for (let i = 0; i < maxRetries; i++) {
      try {
        await assertionFunc();

        return;
      } catch (error) {
        lastError = error;
        if (i < maxRetries - 1) {
          console.warn(`Assertion failed, retrying in ${retryInterval}ms...`);
          await new Promise((resolve) => setTimeout(resolve, retryInterval));
        }
      }
    }
    throw lastError;
  }

  async checkNewEntry() {
    const newEntry = await this.page.locator('[data-index="0"]');
    await this.page.waitForTimeout(2000);
    await this.page.reload();

    await this.assertWithRetry(async () => {
      const organization = await newEntry
        .locator('a[data-test="organization-name-in-all-orgs-table"]')
        .first()
        .innerText();
      expect(organization).toBe('Unnamed');
    });

    await this.assertWithRetry(async () => {
      const website = await newEntry
        .locator('p[data-test="organization-website-in-all-orgs-table"]')
        .first()
        .innerText();
      expect(website).toBe('Unknown');
    });

    await this.assertWithRetry(async () => {
      const relationship = await newEntry
        .locator('p[data-test="organization-relationship-in-all-orgs-table"]')
        .innerText();
      expect(relationship).toBe('Prospect');
    });

    await this.assertWithRetry(async () => {
      const health = await newEntry
        .locator('span[data-test="organization-health-in-all-orgs-table"]')
        .innerText();
      expect(health).toBe('No set');
    });

    await this.assertWithRetry(async () => {
      const nextRenewal = await newEntry
        .locator(
          'span[data-test="organization-next-renewal-in-all-orgs-table"]',
        )
        .innerText();
      expect(nextRenewal).toBe('No contract');
    });

    await this.assertWithRetry(async () => {
      const onboarding = await newEntry
        .locator('p[data-test="organization-onboarding-in-all-orgs-table"]')
        .innerText();
      expect(onboarding).toBe('Not applicable');
    });

    await this.assertWithRetry(async () => {
      const arrForecast = await newEntry
        .locator(
          'span[data-test="organization-arr-forecast-in-all-orgs-table"]',
        )
        .innerText();
      expect(arrForecast).toBe('No contract');
    });

    const maxAttempts = 3;
    const evaluationTimeout = 5000; // 5 seconds, adjust as needed

    for (let attempts = 0; attempts < maxAttempts; attempts++) {
      try {
        await this.assertWithRetry(async () => {
          const owner = await newEntry
            .locator('p[data-test="organization-owner-in-all-orgs-table"]')
            .innerText();
          expect(owner).toBe('Silviu Basu');
        });

        await this.page.waitForSelector(
          '[data-test="organization-last-touchpoint-in-all-orgs-table"]',
          { state: 'attached', timeout: 10000 },
        );

        // Scroll the element into view
        await this.page.evaluate(() => {
          const element = document.querySelector(
            '[data-test="organization-last-touchpoint-in-all-orgs-table"]',
          );
          if (element) {
            element.scrollIntoView({
              behavior: 'auto',
              block: 'center',
              inline: 'center',
            });
          } else {
            console.warn('Last touchpoint element not found');
          }
        });

        // Use the new data-test attribute to locate and assert the element
        await Promise.race([
          this.assertWithRetry(async () => {
            const lastTouchpoint = this.page.locator(
              '[data-test="organization-last-touchpoint-in-all-orgs-table"]',
            );
            await expect(lastTouchpoint).toBeVisible();
            await expect(lastTouchpoint).toHaveText('Created');
          }),
          new Promise((_, reject) =>
            setTimeout(
              () => reject(new Error('Evaluation timed out')),
              evaluationTimeout,
            ),
          ),
        ]);

        break; // Success, exit the loop
      } catch (error) {
        console.error(`Attempt ${attempts + 1} failed:`, error);
        if (attempts === maxAttempts - 1) {
          throw error;
        }
        // Reload the page before the next attempt
        await this.page.waitForTimeout(10000);
        await this.page.reload();
      }
    }
  }

  async goToCustomersPage() {
    await this.page.click('button[data-test="side-nav-item-Customers"]');
  }

  async goToAllOrgsPage() {
    const allOrgsMenubutton = this.page.locator(
      'button[data-test="side-nav-item-all-orgs"]',
    );
    await allOrgsMenubutton.click();
  }

  async updateOrgToCustomer() {
    await this.page.click(
      'button[data-test="organization-relationship-button-in-all-orgs-table"]',
    );
    await this.page.click('div[data-test="relationship-CUSTOMER"]');
  }
}
