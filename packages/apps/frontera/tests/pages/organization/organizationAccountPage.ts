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

  private contractMenuDots = 'button[data-test="contract-menu-dots"]';
  private contractBillingDetailsAddresss =
    'button[data-test="contract-billing-details-address"]';
  private contractMenuEditContract =
    'button[data-test="contract-menu-edit-contract"]';
  private contractMenuDeleteContract =
    'div[data-test="contract-menu-delete-contract"]';
  private contractCardConfirmContractDeletion =
    'button[data-test="contract-card-confirm-contract-deletion"]';
  private contractCardAddSli = 'button[data-test="contract-card-add-sli"]';

  async addContractEmpty() {
    await this.page.click(this.orgAccountEmptyAddContract);
  }

  async addContractNonEmpty() {
    await this.page.click(this.orgAccountNonEmptyAddContract);
  }

  async openContractDotsMenu(contractIndex: number) {
    // await this.page.click(this.contractMenuDots);
    const firstContract = this.page
      .locator(this.contractMenuDots)
      .nth(contractIndex);

    await firstContract.click();
  }

  async deleteContract(contractIndex: number) {
    await this.openContractDotsMenu(contractIndex);
    await this.page.click(this.contractMenuDeleteContract);
    await this.page.click(this.contractCardConfirmContractDeletion);
  }

  async addServices() {
    await this.page.click(this.orgAccountAddServices);
  }

  async checkContractsCount(expectedNumberOfContracts: number) {
    const maxAttempts = 3;
    const retryInterval = 20000;

    await retryOperation(
      this.page,
      async () => {
        await assertWithRetry(async () => {
          const elements = this.page.locator(this.orgAccountAddServices);
          const actualNumberOfContracts = await elements.count();

          expect(
            actualNumberOfContracts,
            `Expected to have ${expectedNumberOfContracts} contract(s) and found ${actualNumberOfContracts}`,
          ).toBe(expectedNumberOfContracts);
        });
      },
      maxAttempts,
      retryInterval,
    );
  }
}
