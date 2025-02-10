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

export class CompaniesPage {
  private page: Page;

  private sideNavItemAllOrgs = sideNavSelectors.sideNavItemCompanies;
  private sideNavItemCustomers = sideNavSelectors.sideNavItemCustomers;
  private sideNavItemCustomersSelected = sideNavSelectors.sideNavItemCustomers;
  private finderTableOrganizations = '[data-test="finder-table-ORGANIZATIONS"]';
  allOrgsAddOrg = 'button[data-test="all-orgs-add-org"]';
  private createOrganizationFromTable =
    'button[data-test="create-organization-from-table"]';
  private organizationsCreateNewOrgOrgName =
    'input[data-test="organizations-create-new-org-org-name"]';
  private addOrgModalAddOrg = 'div[data-test="add-org-modal-add-org"]';
  private organizationNameInAllOrgsTable =
    '[data-test="organization-name-in-all-orgs-table"]';
  private organizationWebsiteInAllOrgsTable =
    'span[data-test="organization-website-in-all-orgs-table"]';
  private organizationRelationshipInAllOrgsTable =
    'p[data-test="organization-relationship-in-all-orgs-table"]';
  private organizationHealthInAllOrgsTable =
    'span[data-test="organization-health-in-all-orgs-table"]';
  private organizationNextRenewalInAllOrgsTable =
    'span[data-test="organization-next-renewal-in-all-orgs-table"]';
  private organizationOnboardingInAllOrgsTable =
    'span[data-test="organization-onboarding-in-all-orgs-table"]';
  private organizationArrForecastInAllOrgsTable =
    'span[data-test="organization-arr-forecast-in-all-orgs-table"]';
  private organizationOwnerInAllOrgsTable =
    'div[data-test="organization-owner-in-all-orgs-table"]';
  private organizationContactsInAllOrgsTable =
    'div[data-test="organization-contacts-in-all-orgs-table"]';
  private organizationStageInAllOrgsTable =
    'div[data-test="organization-stage-in-all-orgs-table"]';
  // private organizationLastTouchpointInAllOrgsTable =
  //   '[data-test="organization-last-touchpoint-in-all-orgs-table"]';
  private organizationLastTouchpointDateInAllOrgsTable =
    '[data-test="organization-last-touchpoint-date-in-all-orgs-table"]';
  private organizationRelationshipButtonInAllOrgsTable =
    'button[data-test="organization-relationship-button-in-all-orgs-table"]';
  private relationshipCustomer =
    'div[data-test="org-dashboard-relationship-CUSTOMER"]';
  private allOrgsSelectAllOrgs = 'div[data-test="all-orgs-select-all-orgs"]';
  private orgActionsArchive = 'button[data-test="org-actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';

  constructor(page: Page) {
    this.page = page;
  }

  async goToCompanies() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllOrgs);
  }

  async addInitialOrganization() {
    const initialOrg = true;

    return await this.addOrganization(this.allOrgsAddOrg, initialOrg);
  }

  async addNonInitialCompany(testInfo: TestInfo) {
    const initialOrg = false;

    return await this.addOrganization(
      this.createOrganizationFromTable,
      initialOrg,
      testInfo,
    );
  }

  async addOrganization(
    organizationCreatorLocator: string,
    initialOrg: boolean,
    testInfo?: TestInfo,
  ) {
    await clickLocatorThatIsVisible(this.page, organizationCreatorLocator);

    const organizationName = randomUUID();

    await this.page
      .locator(this.organizationsCreateNewOrgOrgName)
      .fill(organizationName);

    await this.page.waitForTimeout(1000);

    const requestPromise = createRequestPromise(
      this.page,
      'input?.name',
      organizationName,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await clickLocatorThatIsVisible(this.page, this.addOrgModalAddOrg);
    await Promise.all([requestPromise, responsePromise]);

    initialOrg && (await this.page.reload());
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

  async checkNewCompanyEntry(organizationName: string) {
    const maxAttempts = 3;
    const retryInterval = 20000;

    const newEntry = this.page.locator(
      `${this.finderTableOrganizations} div:has(${this.organizationNameInAllOrgsTable}:text("${organizationName}"))`,
    );

    await this.page.waitForTimeout(2000);
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

          expect(stage).toBe('Lead');
        });
      },
      maxAttempts,
      retryInterval,
    );

    await retryOperation(
      this.page,
      async () => {
        await this.page.waitForSelector(
          this.organizationLastTouchpointDateInAllOrgsTable,
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
        }, this.organizationLastTouchpointDateInAllOrgsTable);

        await assertWithRetry(async () => {
          const lastTouchpointDate = await newEntry
            .locator(this.organizationLastTouchpointDateInAllOrgsTable)
            .innerText();

          expect(lastTouchpointDate).toBe('today');
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

  async goToCompaniesPage() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllOrgs);
  }

  async updateCompanyToCustomer(organizationName: string) {
    const rowLocator = this.page
      .locator(`${this.finderTableOrganizations} div[data-index]`)
      .filter({
        has: this.page.locator(
          `${this.organizationNameInAllOrgsTable}:text("${organizationName}")`,
        ),
      });

    await rowLocator
      .locator('[data-test="organization-relationship-in-all-orgs-table"]')
      .click();

    await clickLocatorThatIsVisible(this.page, this.relationshipCustomer);
    await this.page.waitForTimeout(5000);
  }

  async goToCompany(organizationName: string) {
    await this.page
      .locator(
        `${this.finderTableOrganizations} ${this.organizationNameInAllOrgsTable}:has-text("${organizationName}")`,
      )
      .click();
  }

  async selectAllOrgs() {
    const allOrgsSelectAllOrgs = this.page.locator(this.allOrgsSelectAllOrgs);

    await allOrgsSelectAllOrgs.click();
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
