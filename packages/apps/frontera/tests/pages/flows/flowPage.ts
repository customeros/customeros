import { Page, expect } from '@playwright/test';

import { clickLocatorThatIsVisible } from '../../helper';

export class FlowPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  // private flowName = 'div[data-test="flow-name-in-flows-table"]';
  private flowName = 'input[data-test="flows-flow-name"]';
  private navigateToFlows = 'span[data-test="navigate-to-flows"]';
  private flowContacts = 'button[data-test="flow-contacts"]';
  // private saveFlow = 'button[data-test="save-flow"]';
  private startFlow = 'button[data-test="start-flow"]';
  private flowToggleSettings = 'button[data-test="flow-toggle-settings"]';
  private flowTriggerBlock = 'div[data-test="flow-trigger-block"]';
  private flowTriggerBlockOptions =
    'button[data-test="flow-trigger-block-options"]';
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
  private emailSettingsPanelDone =
    'button[data-test="email-settings-panel-done"]';
  private flowActionWait = 'div[data-test="flow-action-wait"]';
  private flowSendLinkedinMessage =
    'span[data-test="flow-send-linkedin-message"]';
  private flowSendLinkedinConnectionRequest =
    'span[data-test="flow-send-linkedin-connection-request"]';
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
        .toHaveValue(expectedFlowName),
      expect.soft(this.page.locator(this.navigateToFlows)).toHaveText('Flows'),
      expect.soft(this.page.locator(this.flowContacts)).toHaveText('0'),
      expect.soft(this.page.locator(this.startFlow).first()).toBeEnabled(),
      expect
        .soft(this.page.locator(this.startFlow).first())
        .toHaveText('Start flow'),
      expect
        .soft(this.page.locator(this.flowToggleSettings).first())
        .toBeEnabled(),
      expect.soft(this.page.locator(this.flowTriggerBlock)).toBeEnabled(),
      expect
        .soft(this.page.locator(this.flowTriggerBlock))
        .toHaveText('Flow triggers when'),
      expect.soft(this.page.locator(this.flowAddStepOrTrigger)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowEndFlow)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowEndFlow)).toHaveText('End Flow'),
      expect.soft(this.page.locator(this.flowZoomIn)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowZoomOut)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowFitToView)).toBeEnabled(),
      expect.soft(this.page.locator(this.flowTidyUp)).toBeEnabled(),
    ]);
    // await clickLocatorThatIsVisible(this.page, this.flowTriggerBlockOptions);
    //
    // await Promise.all([
    //   expect
    //     .soft(this.page.locator(this.triggersHubInput))
    //     .toHaveAttribute('placeholder', 'Search a trigger'),
    //   expect
    //     .soft(this.page.locator(this.flowTriggerRecordAddedManually))
    //     .toHaveText('Record is added to this flow'),
    //   expect
    //     .soft(this.page.locator(this.flowTriggerRecordCreated))
    //     .toHaveText('Record is created'),
    //   expect
    //     .soft(this.page.locator(this.flowTriggerRecordUpdated))
    //     .toHaveText('Record is updated'),
    //   expect
    //     .soft(this.page.locator(this.flowTriggerRecordMatchesCondition))
    //     .toHaveText('Record matches condition'),
    //   expect
    //     .soft(this.page.locator(this.flowTriggerWebhook))
    //     .toHaveText('Webhook is called'),
    // ]);
    //
    // await this.page.keyboard.press('Escape');
    //
    // await clickLocatorThatIsVisible(this.page, this.flowAddStepOrTrigger);
    // await Promise.all([
    //   expect(this.page.locator(this.stepsHubInput)).toHaveAttribute(
    //     'placeholder',
    //     'Search a step',
    //   ),
    //   expect(this.page.locator(this.flowActionSendEmail)).toHaveText(
    //     'Send email',
    //   ),
    //   expect(this.page.locator(this.flowActionWait)).toHaveText('Wait'),
    //   expect(this.page.locator(this.flowSendLinkedinMessage)).toHaveText(
    //     'Send LinkedIn message',
    //   ),
    //   expect(
    //     this.page.locator(this.flowSendLinkedinConnectionRequest),
    //   ).toHaveText('Send connection request'),
    //   expect(this.page.locator(this.flowCreateRecord)).toHaveText(
    //     'Create record',
    //   ),
    //   expect(this.page.locator(this.flowUpdateRecord)).toHaveText(
    //     'Update record',
    //   ),
    //   expect(this.page.locator(this.flowEnrichRecord)).toHaveText(
    //     'Enrich record',
    //   ),
    //   expect(this.page.locator(this.flowVerifyRecordProperty)).toHaveText(
    //     'Verify record property',
    //   ),
    //   expect(this.page.locator(this.flowConditions)).toHaveText('Conditions'),
    //   expect(this.page.locator(this.flowCreateToDo)).toHaveText('Create to-do'),
    // ]);
    // await clickLocatorThatIsVisible(this.page, this.flowActionSendEmail);
    // await clickLocatorThatIsVisible(this.page, this.emailSettingsPanelDone);

    await clickLocatorThatIsVisible(this.page, this.flowTidyUp);
    // await Promise.all([
    //   await expect(this.page.locator(this.saveFlow)).toHaveText(
    //     'Publish changes',
    //   ),
    //   await expect(this.page.locator(this.saveFlow)).toBeEnabled(),
    // ]);

    // await clickLocatorThatIsVisible(this.page, this.flowToggleSettings);
    // await expect(this.page.locator(this.flowAddSenders)).toBeEnabled();
    //
    // const requestPromise = createRequestPromise(
    //   this.page,
    //   'name',
    //   expectedFlowName,
    // );
    //
    // const responsePromise = createResponsePromise(
    //   this.page,
    //   'flow_Merge?.metadata?.id',
    //   undefined,
    // );

    // await clickLocatorThatIsVisible(this.page, this.saveFlow);
    await this.goToFlows();
    // await Promise.all([requestPromise, responsePromise]);
  }

  async goToFlows() {
    await clickLocatorThatIsVisible(this.page, this.navigateToFlows);
  }
}
