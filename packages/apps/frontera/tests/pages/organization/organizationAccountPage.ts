import { Page, expect } from '@playwright/test';

import { retryOperation, assertWithRetry } from '../../helper';

export class OrganizationAccountPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgAccountEmptyAddContract =
    'button[data-test="org-account-empty-new-contract"]';
  private orgAccountNonEmptyAddContract =
    'button[data-test="org-account-nonempty-new-contract"]';
  private orgAccountAddServices =
    'button[data-test="org-account-add-services"]';

  async addContractEmpty() {
    await this.page.click(this.orgAccountEmptyAddContract);
  }

  async addContractNonEmpty() {
    await this.page.click(this.orgAccountNonEmptyAddContract);
  }

  async addServices() {
    await this.page.click(this.orgAccountAddServices);
  }

  async checkContractsCount(expectedNumberOfContracts: number) {
    const maxAttempts = 3;
    const retryInterval = 20000;

    await retryOperation(
      async () => {
        await assertWithRetry(async () => {
          const elements = this.page.locator(this.orgAccountAddServices);
          const actualNumberOfContracts = await elements.count();
          expect(
            actualNumberOfContracts,
            `Expected to have ${expectedNumberOfContracts} customer(s) and found ${actualNumberOfContracts}`,
          ).toBe(expectedNumberOfContracts);
        });
      },
      maxAttempts,
      retryInterval,
    );
  }
}
