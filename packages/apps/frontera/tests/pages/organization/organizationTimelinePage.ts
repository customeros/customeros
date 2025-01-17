import { Page, expect, Locator } from '@playwright/test';

import { SettingsAccountsPage } from '../settings/settingsAccounts';
import {
  createRequestPromise,
  createResponsePromise,
  ensureLocatorIsVisible,
  clickLocatorThatIsVisible,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class OrganizationTimelinePage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private timelinetimelineEmailButton =
    'button[data-test="timeline-email-button"]';
  private timelineEmailPermissionPrompt =
    'button[data-test="timeline-email-permission-prompt"]';

  private timelineLogButton = 'button[data-test="timeline-log-button"]';
  private timelineLogEditor = 'div[data-test="timeline-log-editor"]';
  private timelineLogConfirmationButton =
    'button[data-test="timeline-log-confirmation-button"]';
  private timelineLogEntryText = 'div[data-test="timeline-log-entry-text"]';
  private timelineLogEntry = 'div[data-test="timeline-log-entry"]';
  private timelinePoppedUpLogEntryText =
    'div[data-test="timeline-popped-up-log-entry-text"]';
  private timelinePoppedUpLogEntryClose =
    'button[data-test="timeline-popped-up-log-entry-close"]';

  private timelineReminderButton =
    'button[data-test="timeline-reminder-button"]';
  private timelineReminderEditor =
    'textarea[data-test="timeline-reminder-editor"]';
  private timelineReminderList = 'div[data-test="timeline-reminder-list"]';
  private timelineReminderDismiss =
    'button[data-test="timeline-reminder-dismiss"]';

  async ensureEmailPermissionPromptIsRedirecting() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.timelinetimelineEmailButton,
    );
    await clickLocatorsThatAreVisible(
      this.page,
      this.timelineEmailPermissionPrompt,
    );
    await SettingsAccountsPage.ensureSettingsAccountsHeaderIsVisible(this.page);
  }

  async addLogEntry() {
    await clickLocatorsThatAreVisible(this.page, this.timelineLogButton);

    const [logEntryEditor] = await Promise.all([
      this.page
        .locator(this.timelineLogEditor)
        .locator('..')
        .locator('..')
        .locator('span'),
    ]);

    await logEntryEditor.click();

    await logEntryEditor.pressSequentially('Test Log Entry!', {
      delay: 500,
    });

    const responsePromise = createResponsePromise(
      this.page,
      'organization?.timelineEventsTotalCount',
      undefined,
    );

    await clickLocatorsThatAreVisible(
      this.page,
      this.timelineLogConfirmationButton,
    );

    await Promise.all([responsePromise]);
  }

  async ensureLogEntryCanBeAdded() {
    await this.addLogEntry();

    await this.ensureLogEntryAppearedInTheTimeline();
    await this.ensureLogEntryPopupAppears();
  }

  private async ensureLogEntryAppearedInTheTimeline() {
    const editor = await ensureLocatorIsVisible(
      this.page,
      this.timelineLogEntryText,
    );

    // Use a retry mechanism for checking the text
    await expect(async () => {
      const text = await editor.innerText();

      expect(text).toBe('Test Log Entry!');
    }).toPass({ timeout: 10000 });
  }

  private async ensureLogEntryPopupAppears() {
    await clickLocatorsThatAreVisible(this.page, this.timelineLogEntry);

    const divLocator: Locator = this.page.locator(
      this.timelinePoppedUpLogEntryText,
    );

    // Find the span within the div and get its text content
    const spanText: string = await divLocator.locator('span').innerText();

    // Check if the text content matches the expected text
    expect(spanText).toBe('Test Log Entry!');

    await clickLocatorsThatAreVisible(
      this.page,
      this.timelinePoppedUpLogEntryClose,
    );
  }

  async ensureReminderCanBeAdded() {
    await this.addReminder();
    await this.ensureReminderWasCreated();
    await clickLocatorsThatAreVisible(this.page, this.timelineReminderDismiss);
  }

  private async addReminder() {
    await clickLocatorsThatAreVisible(this.page, this.timelineReminderButton);

    const timelineReminderEditor = await clickLocatorThatIsVisible(
      this.page,
      this.timelineReminderEditor,
    );

    await this.page.waitForResponse('**/customer-os-api');

    const requestPromise = createRequestPromise(
      this.page,
      'content',
      'Test Reminder!',
    );

    // Type the note
    await timelineReminderEditor.pressSequentially('Test Reminder!', {
      delay: 500,
    });

    // Wait for both the request to be sent and the response to be received
    await Promise.all([requestPromise]);
  }

  private async ensureReminderWasCreated() {
    const parentDiv = this.page.locator(this.timelineReminderList);

    // Ensure that the parent div has only one child
    const remindersCount = await parentDiv.evaluate(
      (node) => node.children.length,
    );

    expect(
      remindersCount,
      `Expected to find exactly 1 Reminder but found ${remindersCount}`,
    ).toBe(1);

    if (remindersCount === 1) {
      // Locate the textarea within the child
      const textarea = parentDiv.locator(this.timelineReminderEditor);

      // Ensure the textarea contains the text "Test Reminder!"
      await expect(textarea).toHaveValue('Test Reminder!');
    } else {
      expect(
        remindersCount,
        `Expected to find exactly 0 Reminder but found ${remindersCount}`,
      ).toBe(0);
    }
  }
}
