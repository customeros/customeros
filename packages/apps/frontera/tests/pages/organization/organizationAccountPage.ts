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
  private contractCardHeader = 'article[data-test="contract-card-header"]';
  private contractMenuDots = 'button[data-test="contract-menu-dots"]';
  private contractBillingDetailsAddress =
    'button[data-test="contract-billing-details-address"]';
  private contractBillingDetailsAddressCountry =
    'div[data-test="contract-billing-details-address-country"]';
  private contractMenuEditContract =
    'div[data-test="contract-menu-edit-contract"]';
  private contractMenuDeleteContract =
    'div[data-test="contract-menu-delete-contract"]';
  private contractCardConfirmContractDeletion =
    'button[data-test="contract-card-confirm-contract-deletion"]';
  private contractCardAddSli = 'button[data-test="contract-card-add-sli"]';
  private addNewServiceMenuSubscription =
    'div[data-test="add-new-service-menu-subscription"]';
  private addNewServiceMenuOneTime =
    'div[data-test="add-new-service-menu-one-time"]';
  private contractDetailsSaveDraft =
    'button[data-test="contract-details-save-draft"]';
  private subscriptionsInAccountPanel =
    'h1[data-test="account-panel-contract-subscription"]';
  private oneTimeInAccountPanel =
    'h1[data-test="account-panel-contract-one-time"]';
  private billingAddressSave = 'button[data-test="billing-address-save"]';

  async addContractEmpty() {
    await this.page.click(this.orgAccountEmptyAddContract);
  }

  async addContractNonEmpty() {
    await this.page.click(this.orgAccountNonEmptyAddContract);
  }

  async addBillingAddress(contractIndex: number) {
    await this.page.waitForResponse('**/customer-os-api');
    await this.openContractDotsMenu(contractIndex);
    await this.page.click(this.contractMenuEditContract);
    await this.page.click(this.contractBillingDetailsAddress);
    await this.page.click(this.contractBillingDetailsAddressCountry);

    const countryInput = this.page.locator(
      this.contractBillingDetailsAddressCountry,
    );

    await countryInput.pressSequentially(
      'South Georgia and the South Sandwich Islands',
    );
    await countryInput.press('Enter');
    await this.page.click(this.billingAddressSave);
    await this.page.waitForResponse('**/customer-os-api');
    await this.page.click(this.contractDetailsSaveDraft);
    await this.page.waitForResponse('**/customer-os-api');
  }

  async fillInBillingAddress() {
    await this.page.click(this.contractBillingDetailsAddressCountry);

    const countryInput = this.page.locator(
      this.contractBillingDetailsAddressCountry,
    );

    await countryInput.fill('South Georgia and the South Sandwich Islands');
    await countryInput.press('Enter');
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
          const elements = this.page.locator(this.contractCardHeader);
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

  async addSLIsToContract(contractIndex: number) {
    await this.openContractDotsMenu(contractIndex);
    await this.page.click(this.contractMenuEditContract);
    await this.page.click(this.contractCardAddSli);
    await this.page.click(this.addNewServiceMenuSubscription);
    await this.page.click(this.contractCardAddSli);
    await this.page.click(this.addNewServiceMenuOneTime);
    await this.page.click(this.contractDetailsSaveDraft);
  }

  async checkSLIsInAccountPanel() {
    const subscriptionSection = this.page.locator(
      this.subscriptionsInAccountPanel,
    );

    await expect(subscriptionSection).toBeVisible();

    const actualsubscriptionUnnamed = subscriptionSection
      .locator('..')
      .locator('p:has-text("Unnamed")');

    await expect(
      actualsubscriptionUnnamed,
      `Expected to have 1 Subscription SLI(s) and found ${actualsubscriptionUnnamed}`,
    ).toHaveCount(1);

    // Verify the text "Unnamed" under One-time
    const oneTimeSection = this.page.locator(this.oneTimeInAccountPanel);

    await expect(oneTimeSection).toBeVisible();

    const oneTimeUnnamed = oneTimeSection
      .locator('..')
      .locator('p:has-text("Unnamed")');

    await expect(
      oneTimeUnnamed,
      `Expected to have 1 One-time SLI(s) and found ${actualsubscriptionUnnamed}`,
    ).toHaveCount(1);
  }
}
