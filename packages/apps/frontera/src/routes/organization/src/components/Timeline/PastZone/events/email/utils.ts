import { EmailParticipant, InteractionEventParticipant } from '@graphql/types';

export const getEmailParticipantsByType = (
  sentTo: InteractionEventParticipant[],
): {
  cc: EmailParticipant[];
  to: EmailParticipant[];
  bcc: EmailParticipant[];
} => {
  const cc = (sentTo || []).filter(
    (e: InteractionEventParticipant): e is EmailParticipant =>
      e.__typename === 'EmailParticipant' && e.type === 'CC',
  );
  const bcc = (sentTo || []).filter(
    (e: InteractionEventParticipant): e is EmailParticipant =>
      e.__typename === 'EmailParticipant' && e.type === 'BCC',
  );
  const to = (sentTo || []).filter(
    (e: InteractionEventParticipant): e is EmailParticipant => e.type === 'TO',
  );

  return {
    cc,
    bcc,
    to,
  };
};

export const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
