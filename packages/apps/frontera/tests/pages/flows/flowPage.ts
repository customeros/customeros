import { Page, expect } from '@playwright/test';

export class FlowPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private flowName = 'span[data-test="flows-flow-name"]';
  private goToFlows = 'span[data-test="go-to-flows"]';
  private flowContacts = 'button[data-test="flow-contacts"]';
  private startFlow = 'button[data-test="start-flow"]';

  async checkNewFlowEntry(expectedFlowName: string) {
    await expect(this.page.locator(this.flowName)).toHaveText(expectedFlowName);
  }
}
