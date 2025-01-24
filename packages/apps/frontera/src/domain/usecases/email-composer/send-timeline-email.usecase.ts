import { RootStore } from '@store/root';
import { action, computed, reaction, observable } from 'mobx';
import { EmailSenderSelectUsecase } from '@domain/usecases/email-composer/email-sender-select.usecase.ts';
import { EmailParticipantSelectUsecase } from '@domain/usecases/email-composer/email-participant-select.usecase.ts';

import { prepareEmailContent } from '@shared/components/EmailTemplate/EmailTemplate';

type EmailAddress = {
  value: string;
  label: string;
};

type EmailState = {
  from: string;
  subject: string;
  content: string;
  to: EmailAddress[];
  cc: EmailAddress[];
  bcc: EmailAddress[];
  fromProvider: string;
};

export class TimelineEmailUsecase {
  @observable public accessor isSending = false;
  @observable private accessor defaultEmailAddress: EmailAddress = {
    value: '',
    label: '',
  };
  @observable public accessor emailState: EmailState = {
    from: '',
    fromProvider: '',
    to: [],
    cc: [],
    bcc: [],
    subject: '',
    content: '',
  };
  public toSelector: EmailParticipantSelectUsecase;
  public ccSelector: EmailParticipantSelectUsecase;
  public bccSelector: EmailParticipantSelectUsecase;
  public fromSelector: EmailSenderSelectUsecase;
  public emailContent: string = '';
  public subject: string;

  private root = RootStore.getInstance();
  private queryKey: string;
  private updateTimelineCache: (response: unknown, queryKey: string) => void;
  private organizationId: string;

  constructor(
    organizationId: string,
    attendees: string[],
    currentUserId?: string,
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
    reaction(() => this.defaultEmailAddress, this.handleDefaultEmailChange);
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
  private handleDefaultEmailChange = () => {
    if (this.defaultEmailAddress.label) {
      this.emailState = {
        ...this.emailState,
        to: [this.defaultEmailAddress],
      };
    }
  };

  @action
  public setDefaultEmailAddress(email: string) {
    this.defaultEmailAddress = { value: email, label: email };
  }

  @action
  public resetEditor() {
    this.emailContent = '';
    this.subject = '';
    this.showConfirmationDialog = false;
    this.defaultEmailAddress = { value: '', label: '' };
    this.toSelector.reset();
    this.ccSelector.reset();
    this.bccSelector.reset();
  }

  @computed
  public get canExitSafely() {
    const isFormEmpty =
      !this.emailContent.length ||
      this.emailContent === `<p class="my-3"><br></p>`;
    const areFieldsEmpty =
      !this.fromSelector.selectedEmail?.value ||
      !this.toSelector.selectedEmails.length === 0;

    const showEmailEditorConfirmationDialog = !isFormEmpty || !areFieldsEmpty;

    return !showEmailEditorConfirmationDialog;
  }

  @action
  public handleExitEditor() {
    this.resetEditor();
    this.defaultEmailAddress = { value: '', label: '' };
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
            from: this.fromSelector.selectedEmail
              ? this.fromSelector.selectedEmail?.value
              : this.defaultEmailAddress.value,
            to: this.toSelector.selectedEmails?.map(({ value }) => value),
            cc: this.ccSelector.selectedEmails?.map(({ value }) => value),
            bcc: this.bccSelector.selectedEmails?.map(({ value }) => value),
            replyTo: replyToId,
            content: emailContent,
            subject: this.subject,
          },
          {
            onSuccess: (response) => {
              this.updateTimelineCache(response, this.queryKey);
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
