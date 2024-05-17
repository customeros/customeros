import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { SessionStore } from '@store/Session/Session.store';

import { DataSource } from '@graphql/types';

type SendMailPayload = {
  to: string[];
  cc: string[];
  bcc: string[];
  content: string;
  channel: string;
  username: string;
  subject?: string;
  replyTo?: string;
  direction: string;
};

export class MailStore {
  isLoading = false;
  error: null | string = null;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
  }

  async send(
    payload: Omit<SendMailPayload, 'channel' | 'direction' | 'username'>,
    options?: {
      onError?: () => void;
      onSuccess?: (timelineEvent: unknown) => void;
    },
  ) {
    const decoratedPayload: SendMailPayload = {
      channel: 'EMAIL',
      username: this.root.session.value.profile.email,
      direction: 'OUTBOUND',
      ...payload,
    };

    try {
      this.isLoading = true;
      await this.transport.http.post(
        `/comms-api/mail/send/`,
        decoratedPayload,
        {
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      runInAction(() => {
        this.root.ui.toastSuccess(
          'Email successfully sent',
          'send-email-success',
        );
        // ideally should be removed when timeline is refactored to use mobx
      });
      const timelineEvent = generateTimelineEvent(
        decoratedPayload,
        this.root.session.value,
      );

      // this should be removed when timeline is refactored to use mobx
      try {
        options?.onSuccess?.(timelineEvent);
      } catch (_) {
        // this approach it's an ugly solution needed to prevent an unexplained typeerror
        // probably caused by passing outside callbacks within `options` object.
      }
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
        this.root.ui.toastError(
          'We were unable to send this email',
          'send-email-error',
        );
        // ideally should be removed when timeline is refactored to use mobx
        options?.onError?.();
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}

// temporary - must be removed when Timeline is migrated to mobx
const generateTimelineEvent = (
  request: SendMailPayload,
  user?: SessionStore['value'],
) => {
  const generateEmailParticipant = (type: string, email: string) => ({
    __typename: 'EmailParticipant',
    type,
    emailParticipant: {
      email,
      id: Math.random().toString(),
      contacts: [],
      users: [],
      organizations: [],
    },
  });

  const generateEmailParticipants = (type: string, emails: string[]) =>
    emails.map((email) => generateEmailParticipant(type, email));

  return {
    __typename: 'InteractionEvent',
    id: Math.random().toString(),
    date: new Date().toISOString(),
    channel: 'EMAIL',
    content: request.content,
    contentType: 'text/html',
    includes: [],
    issue: null,
    externalLinks: [
      {
        externalUrl: null,
        type: '',
      },
    ],
    repliesTo: request.replyTo,
    summary: null,
    meeting: null,
    actionItems: null,
    sentBy: [
      {
        emailParticipant: {
          email: user?.profile?.email,
          id: Math.random().toString(),
          contacts: [],
          users: [
            {
              __typename: 'User',
              id: Math.random().toString(),
              firstName: user?.profile?.email,
              lastName: '',
            },
          ],
          organizations: [],
        },
      },
    ],
    sentTo: [
      ...generateEmailParticipants('TO', request.to),
      ...generateEmailParticipants('CC', request.cc),
      ...generateEmailParticipants('BCC', request.bcc),
    ],
    interactionSession: {
      name: request.subject,
      events: [],
    },
    source: DataSource.Na,
  };
};
