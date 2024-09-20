import { randomUUID } from 'crypto';
import { Page, expect } from '@playwright/test';

import {
  assertWithRetry,
  createRequestPromise,
  createResponsePromise,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class FlowsPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private finderTableFlows = 'div[data-test="finder-table-FLOW"]';
  sideNavItemAllFlows = 'button[data-test="side-nav-item-all-flows"]';
  sideNavItemAllFlowsSelected =
    'button[data-test="side-nav-item-all-flows"] div[aria-selected="true"]';
  private allOrgsSelectAllOrgs = 'button[data-test="all-orgs-select-all-orgs"]';
  addNewFlow = 'button[data-test="add-new-flow"]';
  createNewFlowModalTitle = 'h1[data-test="create-new-flow-modal-title"]';
  createNewFlowName = 'input[data-test="create-new-flow-name"]';
  cancelCreateNewFlow = 'button[data-test="cancel-create-new-flow"]';
  confirmCreateNewFlow = 'button[data-test="confirm-create-new-flow"]';
  flowNameInAllOrgsTable = 'div[data-test="flow-name-in-all-orgs-table"]';
  flowEndedEarlyInAllOrgsTable =
    'div[data-test="flow-ended-early-in-all-orgs-table"]';
  private flowsActionsArchive = 'button[data-test="actions-archive"]';
  private orgActionsConfirmArchive =
    'button[data-test="org-actions-confirm-archive"]';

  async goToFlows() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllFlows);
  }

  async addFlow() {
    // testInfo?: TestInfo, // organizationCreatorLocator: string,
    await clickLocatorsThatAreVisible(this.page, this.addNewFlow);

    const createNewFlowModalTitleInput = this.page.locator(
      this.createNewFlowModalTitle,
    );

    await expect(createNewFlowModalTitleInput).toHaveText('Create new flow');

    await clickLocatorsThatAreVisible(this.page, this.createNewFlowName);

    const flowName = randomUUID();

    const requestPromise = createRequestPromise(this.page, 'name', flowName);

    const responsePromise = createResponsePromise(
      this.page,
      'flow_Merge?.metadata?.id',
      undefined,
    );

    await this.page.keyboard.type(flowName);
    await this.page.keyboard.press('Enter');

    await Promise.all([requestPromise, responsePromise]);
    await this.page.waitForSelector(
      `${this.finderTableFlows} ${this.flowNameInAllOrgsTable}:has-text("${flowName}")`,
      { timeout: 30000 },
    );

    return flowName;
  }

  async checkNewFlowEntry(flowName: string) {
    const flowNameInAllOrgsTable = this.page
      .locator(
        `${this.finderTableFlows} ${this.flowNameInAllOrgsTable}:has-text("${flowName}")`,
      )
      .locator('..')
      .locator('..')
      .locator('..')
      .locator('..')
      .locator('..');

    await this.page.waitForTimeout(2000);
    await this.page.reload();
    await this.page.waitForSelector('[data-index="0"]', { timeout: 30000 });

    await assertWithRetry(async () => {
      const flow = await flowNameInAllOrgsTable
        .locator(this.flowNameInAllOrgsTable)
        .innerText();

      expect(flow).toBe(flowName);
    });

    await assertWithRetry(async () => {
      const flow = await flowNameInAllOrgsTable
        .locator(this.flowEndedEarlyInAllOrgsTable)
        .innerText();

      expect(flow).toBe('No data yet');
    });
  }

  async waitForPageLoad() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavItemAllFlows);
  }

  async selectAllFlows() {
    const allFlowsSelectAllContacts = this.page.locator(
      this.allOrgsSelectAllOrgs,
    );

    try {
      await allFlowsSelectAllContacts.waitFor({
        state: 'visible',
        timeout: 2000,
      });

      const isVisible = await allFlowsSelectAllContacts.isVisible();

      if (isVisible) {
        await allFlowsSelectAllContacts.click();

        // Wait for a short time to allow for any asynchronous updates
        await this.page.waitForTimeout(100);

        // Check if the button is checked after clicking
        return (
          (await allFlowsSelectAllContacts.getAttribute('aria-checked')) ===
            'true' ||
          (await allFlowsSelectAllContacts.getAttribute('data-state')) ===
            'checked'
        );
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
    await clickLocatorsThatAreVisible(this.page, this.flowsActionsArchive);
  }

  async confirmArchiveOrgs() {
    await clickLocatorsThatAreVisible(this.page, this.orgActionsConfirmArchive);
  }
}
