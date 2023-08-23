import axios from 'axios';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
export type SendMailRequest = {
  username: string;
  content: string;
  channel: string;
  direction: string;
  to: string[];
  cc: string[];
  bcc: string[];
  subject?: string;
  replyTo?: string;
};

export const handleSendEmail = (
  textEmailContent: string,
  to: Array<string> = [],
  cc: Array<string> = [],
  bcc: Array<string> = [],
  replyTo: null | string,
  subject: null | string,
  onSuccess: () => void,
  onError: () => void,
  userEmail?: string | null,
) => {
  const request: SendMailRequest = {
    channel: 'EMAIL',
    username: userEmail || '',
    content: textEmailContent || '',
    direction: 'OUTBOUND',
    to: to,
    cc: cc,
    bcc: bcc,
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
      if (res.data) {
        toastSuccess(
          'Email successfully sent',
          `send-email-success-${subject}`,
        );
        onSuccess();
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
