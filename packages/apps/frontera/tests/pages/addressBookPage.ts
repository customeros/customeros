import { Page, expect } from '@playwright/test';

import {
  retryOperation,
  assertWithRetry,
  clickLocatorsThatAreVisible,
} from '../helper';

export class AddressBookPage {
  private page: Page;

  private sideNavItemAllOrgs = 'button[data-test="side-nav-item-all-orgs"]';
  private allOrgsAddOrg = 'button[data-test="all-orgs-add-org"]';
  private organizationNameInAllOrgsTable =
    'p[data-test="organization-name-in-all-orgs-table"]';
  private organizationWebsiteInAllOrgsTable =
    'p[data-test="organization-website-in-all-orgs-table"]';
  private organizationRelationshipInAllOrgsTable =
    'p[data-test="organization-relationship-in-all-orgs-table"]';
  private organizationHealthInAllOrgsTable =
    'span[data-test="organization-health-in-all-orgs-table"]';
  private organizationNextRenewalInAllOrgsTable =
    'span[data-test="organization-next-renewal-in-all-orgs-table"]';
  private organizationOnboardingInAllOrgsTable =
    'p[data-test="organization-onboarding-in-all-orgs-table"]';
  private organizationArrForecastInAllOrgsTable =
    'span[data-test="organization-arr-forecast-in-all-orgs-table"]';
  private organizationOwnerInAllOrgsTable =
    'p[data-test="organization-owner-in-all-orgs-table"]';
  private organizationContactsInAllOrgsTable =
    'div[data-test="organization-contacts-in-all-orgs-table"]';
  private organizationStageInAllOrgsTable =
    'p[data-test="organization-stage-in-all-orgs-table"]';
  private organizationLastTouchpointInAllOrgsTable =
    '[data-test="organization-last-touchpoint-in-all-orgs-table"]';
  private sideNavItemCustomers = 'button[data-test="side-nav-item-Customers"]';
  private organizationRelationshipButtonInAllOrgsTable =
    'button[data-test="organization-relationship-button-in-all-orgs-table"]';
  private relationshipCustomer = 'div[data-test="relationship-CUSTOMER"]';
  private allOrgsSelectAllOrgs = 'button[data-test="all-orgs-select-all-orgs"]';
  private orgActionsArchive = 'button[data-test="org-actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';
  private addressBookEmptyCreateOrg =
    'button[data-test="address-book-empty-create-org"]';

  constructor(page: Page) {
    this.page = page;
  }

  async waitForPageLoad() {
    const allOrgsMenubutton = this.page.locator(this.sideNavItemAllOrgs);

    await allOrgsMenubutton.waitFor({ state: 'visible' });
    await expect(allOrgsMenubutton).toBeVisible();
    await allOrgsMenubutton.click();
  }

  async addOrganization() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.allOrgsAddOrg,
      this.addressBookEmptyCreateOrg,
    );

    const orgCreationResp = await this.page.waitForResponse(
      '**/customer-os-api',
    );
    const orgCreationJson = await orgCreationResp.json();

    expect(
      'error' in orgCreationJson,
      `The createOrganization mutation returned errors`,
    ).toBe(false);
    await this.page.waitForSelector('[data-index="0"]', { timeout: 30000 });
  }

  async checkNewEntry() {
    const maxAttempts = 3;
    const retryInterval = 20000;

    const newEntry = this.page.locator('[data-index="0"]');

    await this.page.waitForTimeout(2000);
    await this.page.reload();
    await this.page.waitForSelector('[data-index="0"]', { timeout: 30000 });

    await assertWithRetry(async () => {
      const organization = await newEntry
        .locator(this.organizationNameInAllOrgsTable)
        .first()
        .innerText();

      expect(organization).toBe('Unnamed');
    });

    await assertWithRetry(async () => {
      const website = await newEntry
        .locator(this.organizationWebsiteInAllOrgsTable)
        .first()
        .innerText();

      expect(website).toBe('Unknown');
    });

    // await this.page.waitForResponse('**/customer-os-api');

    await assertWithRetry(async () => {
      const relationship = await newEntry
        .locator(this.organizationRelationshipInAllOrgsTable)
        .innerText();

      expect(relationship).toBe('Prospect');
    });

    await assertWithRetry(async () => {
      const health = await newEntry
        .locator(this.organizationHealthInAllOrgsTable)
        .innerText();

      expect(health).toBe('No set');
    });

    await assertWithRetry(async () => {
      const nextRenewal = await newEntry
        .locator(this.organizationNextRenewalInAllOrgsTable)
        .innerText();

      expect(nextRenewal).toBe('No contract');
    });

    await assertWithRetry(async () => {
      const onboarding = await newEntry
        .locator(this.organizationOnboardingInAllOrgsTable)
        .innerText();

      expect(onboarding).toBe('Not applicable');
    });

    await assertWithRetry(async () => {
      const arrForecast = await newEntry
        .locator(this.organizationArrForecastInAllOrgsTable)
        .innerText();

      expect(arrForecast).toBe('No contract');
    });

    await retryOperation(
      this.page,
      async () => {
        await assertWithRetry(async () => {
          const owner = await newEntry
            .locator(this.organizationOwnerInAllOrgsTable)
            .innerText();

          expect(owner).toBe('Silviu Basu');
        });
      },
      maxAttempts,
      retryInterval,
    );

    await retryOperation(
      this.page,
      async () => {
        await this.page.waitForSelector(
          this.organizationContactsInAllOrgsTable,
          { state: 'attached', timeout: 10000 },
        );

        await this.page.evaluate((selector) => {
          const element = document.querySelector(selector);

          if (element) {
            element.scrollIntoView({
              behavior: 'auto',
              block: 'center',
              inline: 'center',
            });
          } else {
            console.warn('Contacts element not found');
          }
        }, this.organizationContactsInAllOrgsTable);

        await assertWithRetry(async () => {
          const contacts = await newEntry
            .locator(this.organizationContactsInAllOrgsTable)
            .innerText();

          expect(contacts).toBe('0');
        });
      },
      maxAttempts,
      retryInterval,
    );

    await retryOperation(
      this.page,
      async () => {
        await this.page.waitForSelector(this.organizationStageInAllOrgsTable, {
          state: 'attached',
          timeout: 10000,
        });

        await this.page.evaluate((selector) => {
          const element = document.querySelector(selector);

          if (element) {
            element.scrollIntoView({
              behavior: 'auto',
              block: 'center',
              inline: 'center',
            });
          } else {
            console.warn('Stage element not found');
          }
        }, this.organizationStageInAllOrgsTable);

        await assertWithRetry(async () => {
          const stage = await newEntry
            .locator(this.organizationStageInAllOrgsTable)
            .innerText();

          expect(stage).toBe('Target');
        });
      },
      maxAttempts,
      retryInterval,
    );

    await retryOperation(
      this.page,
      async () => {
        await this.page.waitForSelector(
          this.organizationLastTouchpointInAllOrgsTable,
          { state: 'attached', timeout: 10000 },
        );

        await this.page.evaluate((selector) => {
          const element = document.querySelector(selector);

          if (element) {
            element.scrollIntoView({
              behavior: 'auto',
              block: 'center',
              inline: 'center',
            });
          } else {
            console.warn('Last touchpoint element not found');
          }
        }, this.organizationLastTouchpointInAllOrgsTable);

        await assertWithRetry(async () => {
          const lastTouchpoint = await newEntry
            .locator(this.organizationLastTouchpointInAllOrgsTable)
            .innerText();

          expect(lastTouchpoint).toBe('Created');
        });
      },
      maxAttempts,
      retryInterval,
    );
  }

  async goToCustomersPage() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemCustomers);
  }

  async goToAllOrgsPage() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllOrgs);
  }

  async updateOrgToCustomer() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.organizationRelationshipButtonInAllOrgsTable,
      this.relationshipCustomer,
    );
  }

  async goToOrganization() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.organizationNameInAllOrgsTable,
    );
  }

  async selectAllOrgs() {
    const allOrgsSelectAllOrgs = this.page.locator(this.allOrgsSelectAllOrgs);

    await allOrgsSelectAllOrgs.waitFor({ state: 'visible' });

    const isVisible = await allOrgsSelectAllOrgs.isVisible();

    if (isVisible) {
      await allOrgsSelectAllOrgs.click();
    }
  }

  async archiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.orgActionsArchive);
  }

  async confirmArchiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.orgActionsConfirmArchive);
  }
}
