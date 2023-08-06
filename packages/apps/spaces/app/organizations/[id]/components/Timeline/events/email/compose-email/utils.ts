import { SendMailRequest } from '@spaces/molecules/conversation-timeline-item/types';
import axios from 'axios';
import { toastError, toastSuccess } from '@ui/presentation/Toast';

export const handleSendEmail = (
  textEmailContent: string,
  destination: Array<string> = [],
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
    destination: destination,
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
