import { Page, expect } from '@playwright/test';

import {
  createRequestPromise,
  createResponsePromise,
  clickLocatorThatIsVisible,
} from '../../helper';

export class FlowPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private flowName = 'span[data-test="flows-flow-name"]';
  private navigateToFlows = 'span[data-test="navigate-to-flows"]';
  private flowContacts = 'button[data-test="flow-contacts"]';
  private startFlow = 'button[data-test="start-flow"]';
  private flowTidyUp = 'button[data-test="flow-tidy-up"]';
  private saveFlow = 'button[data-test="save-flow"]';

  async checkNewFlowEntry(expectedFlowName: string) {
    await Promise.all([
      expect
        .soft(this.page.locator(this.flowName))
        .toHaveText(expectedFlowName),
      expect.soft(this.page.locator(this.navigateToFlows)).toHaveText('Flows'),
      expect.soft(this.page.locator(this.flowContacts)).toHaveText('0'),
      expect.soft(this.page.locator(this.startFlow)).toHaveText('Start flow'),
      expect.soft(this.page.locator(this.startFlow)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowTidyUp)).toBeEnabled(),
    ]);
    await clickLocatorThatIsVisible(this.page, this.flowTidyUp);
    await Promise.all([
      await expect(this.page.locator(this.saveFlow)).toHaveText('Save'),
      await expect(this.page.locator(this.saveFlow)).toBeEnabled(),
    ]);

    const requestPromise = createRequestPromise(
      this.page,
      'name',
      expectedFlowName,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'flow_Merge?.metadata?.id',
      undefined,
    );

    await clickLocatorThatIsVisible(this.page, this.saveFlow);
    await Promise.all([requestPromise, responsePromise]);
  }

  async goToFlows() {
    await clickLocatorThatIsVisible(this.page, this.navigateToFlows);
  }
}
