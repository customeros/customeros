import { randomUUID } from 'crypto';
import { Page, expect } from '@playwright/test';

import {
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
  flowNameInAllOrgsTable = 'div[data-test="flow-name-in-flows-table"]';
  flowEndedEarlyInFlowTable =
    'div[data-test="flow-ended-early-in-flows-table"]';
  flowNotStartedInFlowsTable =
    'div[data-test="flow-not-started-in-flows-table"]';
  flowStatusInFlowsTable = 'p[data-test="flow-status-in-flows-table"]';
  flowInProgressInFlowsTable =
    'div[data-test="flow-in-progress-in-flows-table"]';
  flowCompletedInFlowsTable = 'div[data-test="flow-completed-in-flows-table"]';
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

  async checkNewFlowEntry(expectedFlowName: string) {
    const flowNameInAllOrgsTable = this.page
      .locator(
        `${this.finderTableFlows} ${this.flowNameInAllOrgsTable}:has-text("${expectedFlowName}")`,
      )
      .locator('..')
      .locator('..')
      .locator('..')
      .locator('..')
      .locator('..');

    await this.page.waitForTimeout(2000);
    await this.page.reload();
    await this.page.waitForSelector('[data-index="0"]', { timeout: 30000 });

    const actualFlow = await flowNameInAllOrgsTable
      .locator(this.flowNameInAllOrgsTable)
      .innerText();

    const actualFlowEndedEarly = await flowNameInAllOrgsTable
      .locator(this.flowEndedEarlyInFlowTable)
      .innerText();

    const actualFlowNotStarted = await flowNameInAllOrgsTable
      .locator(this.flowNotStartedInFlowsTable)
      .innerText();

    const actualFlowStatusInAllOrgsTable = await flowNameInAllOrgsTable
      .locator(this.flowStatusInFlowsTable)
      .innerText();

    const actualflowInProgressInFlowsTable = await flowNameInAllOrgsTable
      .locator(this.flowInProgressInFlowsTable)
      .innerText();

    const actualFlowCompletedInFlowsTable = await flowNameInAllOrgsTable
      .locator(this.flowCompletedInFlowsTable)
      .innerText();

    await Promise.all([
      expect.soft(actualFlow).toBe(expectedFlowName),
      expect.soft(actualFlowEndedEarly).toBe('No data yet'),
      expect.soft(actualFlowNotStarted).toBe('No data yet'),
      expect.soft(actualFlowStatusInAllOrgsTable).toBe('Not Started'),
      expect.soft(actualflowInProgressInFlowsTable).toBe('No data yet'),
      expect.soft(actualFlowCompletedInFlowsTable).toBe('No data yet'),
    ]);
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
