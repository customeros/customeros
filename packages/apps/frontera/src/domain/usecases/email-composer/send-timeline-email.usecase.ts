import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';
import { EmailSenderSelectUsecase } from '@domain/usecases/email-composer/email-sender-select.usecase.ts';
import { EmailParticipantSelectUsecase } from '@domain/usecases/email-composer/email-participant-select.usecase.ts';

import { prepareEmailContent } from '@shared/components/EmailTemplate/EmailTemplate';

export class TimelineEmailUsecase {
  @observable public accessor isSending = false;
  @observable public accessor subject: string = '';

  public toSelector: EmailParticipantSelectUsecase;
  public ccSelector: EmailParticipantSelectUsecase;
  public bccSelector: EmailParticipantSelectUsecase;
  public fromSelector: EmailSenderSelectUsecase;
  public emailContent: string = '';

  private root = RootStore.getInstance();
  private organizationId: string;
  private onSuccess: (response: unknown) => void;

  constructor(
    organizationId: string,
    attendees: string[],
    currentUserId: string,
    onSuccess: (response: unknown) => void,
  ) {
    this.organizationId = organizationId;

    this.fromSelector = new EmailSenderSelectUsecase(
      organizationId,
      attendees,
      currentUserId,
    );
    this.toSelector = new EmailParticipantSelectUsecase(organizationId);
    this.ccSelector = new EmailParticipantSelectUsecase(organizationId);
    this.bccSelector = new EmailParticipantSelectUsecase(organizationId);
    this.updateEmailContent = this.updateEmailContent.bind(this);
    this.updateSubject = this.updateSubject.bind(this);
    this.createEmail = this.createEmail.bind(this);
    this.onSuccess = onSuccess;
  }

  @action
  public updateEmailContent(emailContent: string) {
    this.emailContent = emailContent;
  }

  @action
  public updateSubject(subject: string) {
    this.subject = subject;
  }

  @action
  public resetEditor() {
    this.emailContent = '';
    this.subject = '';
    this.toSelector.reset();
    this.ccSelector.reset();
    this.bccSelector.reset();
  }

  @computed
  public get canExitSafely() {
    const isFormEmpty =
      !this.emailContent.length ||
      this.emailContent === `<p class="my-3"><br></p>`;

    const showEmailEditorConfirmationDialog = !isFormEmpty;

    return !showEmailEditorConfirmationDialog;
  }

  @action
  public handleExitEditor() {
    this.resetEditor();
  }

  @action
  public async createEmail(replyToId?: string) {
    this.isSending = true;

    try {
      const emailContent = await prepareEmailContent(this.emailContent);

      if (emailContent) {
        await this.root.mail.send(
          {
            fromProvider: this.fromSelector.selectedEmail?.provider ?? '',
            from: this.fromSelector.selectedEmail?.value ?? '',
            to: this.toSelector.selectedEmails?.map(({ value }) => value),
            cc: this.ccSelector.selectedEmails?.map(({ value }) => value),
            bcc: this.bccSelector.selectedEmails?.map(({ value }) => value),
            replyTo: replyToId,
            content: emailContent,
            subject: this.subject,
          },
          {
            onSuccess: (response) => {
              this.onSuccess(response);
              this.isSending = false;
              this.resetEditor();
            },
            onError: () => {
              this.isSending = false;
            },
          },
        );
      }
    } catch (error) {
      console.error('Error saving email:', error);
      this.root.ui.toastError('Error saving email', 'email-save-error');
      this.isSending = false;
    }
  }
}
