import { Page, expect } from '@playwright/test';

import {
  retryOperation,
  assertWithRetry,
  clickLocatorsThatAreVisible,
  clickLocatorThatIsVisibleWithIndex,
  clickLocatorThatIsVisibleAndHasText,
  clickLocatorWithSublocatorThatIsVisible,
} from '../../helper';

export class OrganizationAccountPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgAccountEmptyAddContract =
    'button[data-test="org-account-empty-new-contract"]';
  private orgAccountNonEmptyAddContract =
    'button[data-test="org-account-nonempty-new-contract"]';
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
  private organizationAccountNotesEditor =
    'div[data-test="organization-account-notes-editor"]';
  private organizationAccountRelationship =
    'div[data-test="organization-account-relationship"]';
  private relationshipCustomer = 'div[data-test="relationship-CUSTOMER"]';

  async addContractEmpty() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.orgAccountEmptyAddContract,
    );

    const response = await this.page.waitForResponse(
      (response) =>
        response.url().includes('customer-os-api') &&
        response
          .json()
          .then((body) => body.data && body.data.contract_Update !== undefined)
          .catch(() => false),
    );

    // Optional: You can log or assert the response data if needed
    await response.json();
    // await this.page.waitForResponse('**/customer-os-api');
  }

  async addContractNonEmpty() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.orgAccountNonEmptyAddContract,
    );
  }

  async addBillingAddress(contractIndex: number) {
    await this.page.waitForResponse('**/customer-os-api');
    await this.openContractDotsMenu(contractIndex);
    await clickLocatorsThatAreVisible(
      this.page,
      this.contractMenuEditContract,
      this.contractBillingDetailsAddress,
      this.contractBillingDetailsAddressCountry,
    );

    const countryInput = this.page.locator(
      this.contractBillingDetailsAddressCountry,
    );

    await countryInput.pressSequentially(
      'South Georgia and the South Sandwich Islands',
    );
    await countryInput.press('Enter');
    await clickLocatorsThatAreVisible(this.page, this.billingAddressSave);
    await this.page.waitForResponse('**/customer-os-api');
    await clickLocatorsThatAreVisible(this.page, this.contractDetailsSaveDraft);
    await this.page.waitForResponse('**/customer-os-api');
  }

  async openContractDotsMenu(contractIndex: number) {
    await clickLocatorThatIsVisibleWithIndex(
      this.page,
      this.contractMenuDots,
      contractIndex,
    );
  }

  async deleteContract(contractIndex: number) {
    await this.openContractDotsMenu(contractIndex);
    await clickLocatorsThatAreVisible(
      this.page,
      this.contractMenuDeleteContract,
      this.contractCardConfirmContractDeletion,
    );
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
    await clickLocatorsThatAreVisible(
      this.page,
      this.contractMenuEditContract,
      this.contractCardAddSli,
      this.addNewServiceMenuSubscription,
      this.contractCardAddSli,
      this.addNewServiceMenuOneTime,
      this.contractDetailsSaveDraft,
    );
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

  async updateOrgToCustomer() {
    await clickLocatorThatIsVisibleAndHasText(
      this.page,
      this.organizationAccountRelationship,
      'Prospect',
    );

    await clickLocatorThatIsVisibleAndHasText(
      this.page,
      this.relationshipCustomer,
      'Customer',
    );
  }

  async addNoteToOrg() {
    const editor = await clickLocatorWithSublocatorThatIsVisible(
      this.page,
      this.organizationAccountNotesEditor,
      '.ProseMirror',
    );

    // Set up the request listener before typing
    const requestPromise = this.page.waitForRequest((request) => {
      if (
        request.method() === 'POST' &&
        request.url().includes('customer-os-api')
      ) {
        const postData = request.postData();

        if (postData) {
          const parsedData = JSON.parse(postData);

          return parsedData.variables.input.notes === '<p>Test Note!</p>';
        }
      }

      return false;
    });

    // Set up the response listener
    const responsePromise = this.page.waitForResponse(async (response) => {
      if (
        response.request().method() === 'POST' &&
        response.url().includes('customer-os-api')
      ) {
        const responseBody = await response.json();

        // Adjust this condition based on your API response structure
        return (
          responseBody.data?.organization_Update?.metadata?.id !== undefined
        );
      }

      return false;
    });

    // Type the note
    await editor.pressSequentially('Test Note!', { delay: 500 });

    // Wait for both the request to be sent and the response to be received
    await Promise.all([requestPromise, responsePromise]);

    // Wait a bit to ensure the UI has updated
    await this.page.waitForTimeout(1000);

    await this.page.reload();

    // Wait for the editor to be visible after reload
    await expect(editor).toBeVisible({ timeout: 20000 });

    // Use a retry mechanism for checking the text
    await expect(async () => {
      const text = await editor.innerText();

      expect(text).toBe('Test Note!');
    }).toPass({ timeout: 10000 });
  }
}
