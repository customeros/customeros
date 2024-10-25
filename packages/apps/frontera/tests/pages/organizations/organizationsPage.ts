import { randomUUID } from 'crypto';
import { Page, expect, TestInfo } from '@playwright/test';

import { sideNavSelectors } from '../sideNavSelectors';
import {
  retryOperation,
  assertWithRetry,
  createRequestPromise,
  createResponsePromise,
  ensureLocatorIsVisible,
  clickLocatorThatIsVisible,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class OrganizationsPage {
  private page: Page;

  private sideNavItemAllOrgs = sideNavSelectors.sideNavItemAllOrgs;
  private sideNavItemCustomers = sideNavSelectors.sideNavItemCustomers;
  private sideNavItemCustomersSelected =
    sideNavSelectors.sideNavItemCustomersSelected;
  private finderTableOrganizations =
    'div[data-test="finder-table-ORGANIZATIONS"]';
  private allOrgsAddOrg = 'button[data-test="all-orgs-add-org"]';
  private createOrganizationFromTable =
    'button[data-test="create-organization-from-table"]';
  private organizationsCreateNewOrgOrgName =
    'input[data-test="organizations-create-new-org-org-name"]';
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
  private organizationRelationshipButtonInAllOrgsTable =
    'button[data-test="organization-relationship-button-in-all-orgs-table"]';
  private relationshipCustomer = 'div[data-test="relationship-CUSTOMER"]';
  private allOrgsSelectAllOrgs = 'button[data-test="all-orgs-select-all-orgs"]';
  private orgActionsArchive = 'button[data-test="org-actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';

  constructor(page: Page) {
    this.page = page;
  }

  async goToAllOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllOrgs);
  }

  async addInitialOrganization() {
    return await this.addOrganization(this.createOrganizationFromTable);
  }

  async addNonInitialOrganization(testInfo: TestInfo) {
    return await this.addOrganization(
      this.createOrganizationFromTable,
      testInfo,
    );
  }

  async addOrganization(
    organizationCreatorLocator: string,
    testInfo?: TestInfo,
  ) {
    await clickLocatorsThatAreVisible(
      this.page,
      organizationCreatorLocator,
      this.organizationsCreateNewOrgOrgName,
    );

    const organizationName = randomUUID();

    const requestPromise = createRequestPromise(
      this.page,
      'name',
      organizationName,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await this.page.keyboard.type(organizationName);
    await this.page.waitForTimeout(300);
    await this.page.keyboard.press('Enter');

    await Promise.all([requestPromise, responsePromise]);
    await this.page.waitForSelector(
      `${this.finderTableOrganizations} ${this.organizationNameInAllOrgsTable}:has-text("${organizationName}")`,
      { timeout: 30000 },
    );

    if (testInfo) {
      process.stdout.write(
        '\nOrganization ' +
          organizationName +
          ' was created for the test: ' +
          testInfo.title,
      );
    } else {
      process.stdout.write(
        '\nInitial Organization ' + organizationName + ' was created',
      );
    }

    return organizationName;
  }

  async checkNewOrganizationEntry(organizationName: string) {
    const maxAttempts = 3;
    const retryInterval = 20000;

    const newEntry = this.page
      .locator(
        `${this.finderTableOrganizations} ${this.organizationNameInAllOrgsTable}:has-text("${organizationName}")`,
      )
      .locator('..')
      .locator('..')
      .locator('..');

    await this.page.waitForTimeout(2000);
    await this.page.reload();
    await this.page.waitForSelector('[data-index="0"]', { timeout: 30000 });

    await assertWithRetry(async () => {
      const organization = await newEntry
        .locator(this.organizationNameInAllOrgsTable)
        .innerText();

      expect(organization).toBe(organizationName);
    });

    await assertWithRetry(async () => {
      const website = await newEntry
        .locator(this.organizationWebsiteInAllOrgsTable)
        .innerText();

      expect(website).toBe('Not set');
    });

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

          expect(owner).toBe('No owner');
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

    await this.page.waitForTimeout(2000);
  }

  async goToCustomersPage() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemCustomers);

    await ensureLocatorIsVisible(this.page, this.sideNavItemCustomersSelected);
  }

  async goToAllOrgsPage() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllOrgs);
  }

  async updateOrgToCustomer(organizationName: string) {
    const newEntry = this.page
      .locator(
        `${this.finderTableOrganizations} ${this.organizationNameInAllOrgsTable}:has-text("${organizationName}")`,
      )
      .locator('..')
      .locator('..')
      .locator('..');

    await newEntry
      .locator(this.organizationRelationshipButtonInAllOrgsTable)
      .click();
    await clickLocatorThatIsVisible(this.page, this.relationshipCustomer);
  }

  async goToOrganization(organizationName: string) {
    await this.page
      .locator(
        `${this.finderTableOrganizations} ${this.organizationNameInAllOrgsTable}:has-text("${organizationName}")`,
      )
      .click();
  }

  async selectAllOrgs(): Promise<boolean> {
    const allOrgsSelectAllOrgs = this.page.locator(this.allOrgsSelectAllOrgs);

    try {
      await allOrgsSelectAllOrgs.waitFor({ state: 'visible', timeout: 10000 });

      const isVisible = await allOrgsSelectAllOrgs.isVisible();

      if (isVisible) {
        await allOrgsSelectAllOrgs.click();

        const isAllOrgsSelectAllOrgs = await allOrgsSelectAllOrgs.getAttribute(
          'data-state',
        );

        return isAllOrgsSelectAllOrgs === 'checked';
      }
    } catch (error) {
      if (error.name === 'TimeoutError') {
        // Silently return false if the element is not found
        return false;
      }
      // Re-throw any other errors
      throw error;
    }

    return false;
  }

  async archiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.orgActionsArchive);
  }

  async confirmArchiveOrgs() {
    const responsePromise = this.page.waitForResponse(async (response) => {
      if (
        response.request().method() === 'POST' &&
        response.url().includes('customer-os-api')
      ) {
        const responseBody = await response.json();

        return responseBody.data?.organization_HideAll?.result !== undefined;
      }

      return false;
    });

    await clickLocatorsThatAreVisible(this.page, this.orgActionsConfirmArchive);

    await Promise.all([responsePromise]);
  }
}
