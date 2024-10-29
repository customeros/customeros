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
  private saveFlow = 'button[data-test="save-flow"]';
  private startFlow = 'button[data-test="start-flow"]';
  private flowToggleSettings = 'button[data-test="flow-toggle-settings"]';
  private flowTriggerBlock = 'span[data-test="flow-trigger-block"]';
  private triggersHubInput = 'input[data-test="TriggersHub-input"]';
  private flowTriggerRecordAddedManually =
    'div[data-test="flow-trigger-record-added-manually"]';
  private flowTriggerRecordCreated =
    'span[data-test="flow-trigger-record-created"]';
  private flowTriggerRecordUpdated =
    'span[data-test="flow-trigger-record-updated"]';
  private flowTriggerRecordMatchesCondition =
    'span[data-test="flow-trigger-record-matches-condition"]';
  private flowTriggerWebhook = 'span[data-test="flow-trigger-webhook"]';
  private flowAddStepOrTrigger = 'button[data-test="flow-add-step-or-trigger"]';
  private stepsHubInput = 'input[data-test="StepsHub-input"]';
  private flowActionSendEmail = 'div[data-test="flow-action-send-email"]';
  private flowActionWait = 'div[data-test="flow-action-wait"]';
  private flowSendLinkedinMessage =
    'span[data-test="flow-send-linkedin-message"]';
  private flowCreateRecord = 'span[data-test="flow-create-record"]';
  private flowUpdateRecord = 'span[data-test="flow-update-record"]';
  private flowEnrichRecord = 'span[data-test="flow-enrich-record"]';
  private flowVerifyRecordProperty =
    'span[data-test="flow-verify-record-property"]';
  private flowConditions = 'span[data-test="flow-conditions"]';
  private flowCreateToDo = 'span[data-test="flow-create-to-do"]';
  private flowEndFlow = 'span[data-test="flow-end-flow"]';
  private flowZoomIn = 'button[data-test="flow-zoom-in"]';
  private flowZoomOut = 'button[data-test="flow-zoom-out"]';
  private flowFitToView = 'button[data-test="flow-fit-to-view"]';
  private flowTidyUp = 'button[data-test="flow-tidy-up"]';
  private flowAddSenders = 'button[data-test="flow-add-senders"]';

  async checkNewFlowEntry(expectedFlowName: string) {
    await Promise.all([
      expect
        .soft(this.page.locator(this.flowName))
        .toHaveText(expectedFlowName),
      expect.soft(this.page.locator(this.navigateToFlows)).toHaveText('Flows'),
      expect.soft(this.page.locator(this.flowContacts)).toHaveText('0'),
      expect.soft(this.page.locator(this.startFlow)).toBeEnabled(),
      expect.soft(this.page.locator(this.startFlow)).toHaveText('Start flow'),
      expect.soft(this.page.locator(this.flowToggleSettings)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowTriggerBlock)).toBeEnabled(),
      expect
        .soft(this.page.locator(this.flowTriggerBlock))
        .toHaveText('What should trigger this flow?'),
      expect.soft(this.page.locator(this.flowAddStepOrTrigger)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowEndFlow)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowEndFlow)).toHaveText('End Flow'),
      expect.soft(this.page.locator(this.flowZoomIn)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowZoomOut)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowFitToView)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowTidyUp)).toBeEnabled(),
    ]);
    await clickLocatorThatIsVisible(this.page, this.flowTriggerBlock);

    await Promise.all([
      expect(this.page.locator(this.triggersHubInput)).toHaveAttribute(
        'placeholder',
        'Search a trigger',
      ),
      expect(this.page.locator(this.flowTriggerRecordAddedManually)).toHaveText(
        'Record added manually...',
      ),
      expect(this.page.locator(this.flowTriggerRecordCreated)).toHaveText(
        'Record created',
      ),
      expect(this.page.locator(this.flowTriggerRecordUpdated)).toHaveText(
        'Record updated',
      ),
      expect(
        this.page.locator(this.flowTriggerRecordMatchesCondition),
      ).toHaveText('Record matches condition'),
      expect(this.page.locator(this.flowTriggerWebhook)).toHaveText('Webhook'),
    ]);

    await this.page.keyboard.press('Escape');

    await clickLocatorThatIsVisible(this.page, this.flowAddStepOrTrigger);
    await Promise.all([
      expect(this.page.locator(this.stepsHubInput)).toHaveAttribute(
        'placeholder',
        'Search a step',
      ),
      expect(this.page.locator(this.flowActionSendEmail)).toHaveText(
        'Send email',
      ),
      expect(this.page.locator(this.flowActionWait)).toHaveText('Wait'),
      expect(this.page.locator(this.flowSendLinkedinMessage)).toHaveText(
        'Send LinkedIn message',
      ),
      expect(this.page.locator(this.flowCreateRecord)).toHaveText(
        'Create record',
      ),
      expect(this.page.locator(this.flowUpdateRecord)).toHaveText(
        'Update record',
      ),
      expect(this.page.locator(this.flowEnrichRecord)).toHaveText(
        'Enrich record',
      ),
      expect(this.page.locator(this.flowVerifyRecordProperty)).toHaveText(
        'Verify record property',
      ),
      expect(this.page.locator(this.flowConditions)).toHaveText('Conditions'),
      expect(this.page.locator(this.flowCreateToDo)).toHaveText('Create to-do'),
    ]);
    await clickLocatorThatIsVisible(this.page, this.flowActionSendEmail);

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

    await clickLocatorThatIsVisible(this.page, this.flowToggleSettings);
    await expect(this.page.locator(this.flowAddSenders)).toBeEnabled();
  }

  async goToFlows() {
    await clickLocatorThatIsVisible(this.page, this.navigateToFlows);
  }
}
