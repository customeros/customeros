import { Page, expect } from '@playwright/test';

export class AllOrgsPage {
  constructor(private page: Page) {}

  async waitForPageLoad() {
    const allOrgsMenubutton = this.page.locator('button:has-text("All orgs")');
    await allOrgsMenubutton.waitFor({ state: 'visible' });
    await expect(allOrgsMenubutton).toBeVisible();
    await allOrgsMenubutton.click();
  }

  async addOrganization() {
    const addOrganizationsButton = await this.page.waitForSelector(
      'button.inline-flex.items-center.justify-center.font-semibold.text-gray-700.border.border-solid.border-gray-300.hover\\:bg-gray-50.hover\\:text-gray-700.focus\\:bg-gray-50.rounded-lg.text-sm',
    );
    await addOrganizationsButton.click();
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
          // console.log(`Assertion failed, retrying in ${retryInterval}ms...`);
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
        .locator('text=Unnamed')
        .first()
        .innerText();
      expect(organization).toBe('Unnamed');
    });

    await this.assertWithRetry(async () => {
      const website = await newEntry
        .locator('text=Unknown')
        .first()
        .innerText();
      expect(website).toBe('Unknown');
    });

    await this.assertWithRetry(async () => {
      const relationship = await newEntry.locator('text=Prospect').innerText();
      expect(relationship).toBe('Prospect');
    });

    await this.assertWithRetry(async () => {
      const health = await newEntry
        .locator('text=Unknown >> nth=1')
        .innerText();
      expect(health).toBe('Unknown');
    });

    await this.assertWithRetry(async () => {
      const nextRenewal = await newEntry
        .locator('text=Unknown >> nth=2')
        .innerText();
      expect(nextRenewal).toBe('Unknown');
    });

    await this.assertWithRetry(async () => {
      const onboarding = await newEntry
        .locator('text=Not applicable')
        .innerText();
      expect(onboarding).toBe('Not applicable');
    });

    await this.assertWithRetry(async () => {
      const arrForecast = await newEntry
        .locator('text=Unknown >> nth=3')
        .innerText();
      expect(arrForecast).toBe('Unknown');
    });

    await this.assertWithRetry(async () => {
      const owner = await newEntry.locator('text=Owner').innerText();
      expect(owner).toBe('Owner');
    });

    const maxAttempts = 3;
    const evaluationTimeout = 5000; // 5 seconds, adjust as needed

    for (let attempts = 0; attempts < maxAttempts; attempts++) {
      try {
        // First, wait for the element to be present
        await this.page.waitForSelector('.flex.flex-1.relative.w-full', {
          timeout: 5000,
        });

        // Then try to scroll
        await this.page.evaluate(() => {
          const element = document.querySelector(
            '.flex.flex-1.relative.w-full',
          );
          if (element) {
            element.scrollTo(10000, 0);
          } else {
            console.warn('Scroll element not found');
          }
        });

        await Promise.race([
          this.assertWithRetry(async () => {
            const lastTouchpoint = await newEntry
              .locator('text=Created')
              .first();
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

        break;
      } catch (error) {
        console.error(`Attempt ${attempts + 1} failed:`, error);

        if (attempts === maxAttempts - 1) {
          throw error;
        }

        // Reload the page before the next attempt
        await this.page.reload();
      }
    }
  }

  async goToCustomersPage() {
    await this.page.click('button:has-text("Customers"):has(svg)');
  }

  async goToAllOrgsPage() {
    const allOrgsMenubutton = this.page.locator('button:has-text("All orgs")');
    await allOrgsMenubutton.click();
  }

  async updateOrgToCustomer() {
    await this.page.click('#edit-button');
    await this.page.click('div.text-gray-700:has-text("Customer")');
  }
}
