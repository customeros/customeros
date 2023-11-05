import axios from 'axios';
import { DefaultSession } from 'next-auth/core/types';

import { DataSource } from '@graphql/types';
import { toastError, toastSuccess } from '@ui/presentation/Toast';

export type SendMailRequest = {
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

const generateTimelineEvent = (
  request: SendMailRequest,
  user?: DefaultSession['user'],
) => {
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
          email: user?.email,
          id: Math.random().toString(),
          contacts: [],
          users: [
            {
              __typename: 'User',
              id: Math.random().toString(),
              firstName: user?.name,
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

export const handleSendEmail = (
  textEmailContent: string,
  to: Array<string> = [],
  cc: Array<string> = [],
  bcc: Array<string> = [],
  replyTo: null | string,
  subject: null | string,
  onSuccess: (res: any) => void,
  onError: () => void,
  user?: DefaultSession['user'],
) => {
  const request: SendMailRequest = {
    channel: 'EMAIL',
    username: user?.email || '',
    content: textEmailContent || '',
    direction: 'OUTBOUND',
    to: to.filter((e) => e),
    cc: cc.filter((e) => e),
    bcc: bcc.filter((e) => e),
  };
  if (replyTo) {
    request.replyTo = replyTo;
  }
  if (subject) {
    request.subject = subject;
  }

  return axios
    .post(`/comms-api/mail/send/`, request, {
      headers: {
        'Content-Type': 'application/json',
      },
    })
    .then((res) => {
      const timelineEvent = generateTimelineEvent(request, user);
      onSuccess(timelineEvent);

      if (res.data) {
        toastSuccess(
          'Email successfully sent',
          `send-email-success-${subject}`,
        );
      }
    })
    .catch((reason) => {
      toastError(
        'We were unable to send this email',
        `send-email-error-${reason}-${subject}`,
      );
      onError();
    });
};
