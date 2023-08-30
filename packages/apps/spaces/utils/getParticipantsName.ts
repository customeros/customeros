import {
  EmailParticipant,
  Organization,
} from '@graphql/types';
import { Contact, InteractionEventParticipant, User } from '@graphql/types';
// todo cleanup, move those to helper functions in timeline events utils folder
export const getName = (
  data: Contact | User | Organization,
  email?: string | null | undefined,
  rawEmail?: string | undefined | null,
): string => {
  if ((data as Contact | Organization)?.name) {
    return <string>(data as Contact | Organization).name;
  }
  const personData: Contact | User = (data as Contact | User)
  if (personData?.firstName || personData?.lastName) {
    return `${personData.firstName} ${personData.lastName}`;
  }
  return email || rawEmail || '';
};

export const getEmailParticipantName = (
  participant: EmailParticipant,
): string => {
  if (!participant?.emailParticipant) {
    return '';
  }
  const { emailParticipant } = participant;
  const { contacts, users, email, rawEmail } = emailParticipant;

  if (contacts.length) {
    return contacts.map((c) => getName(c, email, rawEmail)).join(' ');
  }
  if (users.length) {
    return users.map((c) => getName(c, email, rawEmail)).join(' ');
  }

  return email || rawEmail || '';
};

export const getEmailParticipantsName = (
  participants: EmailParticipant[],
): string => {
  return participants
    ?.map((participant) => getEmailParticipantName(participant))
    .join(', ');
};

export const getParticipantNameAndEmail = (
  participant: EmailParticipant,
  keyName = 'email',
): { [x: string]: string; label: string } => {
  const { emailParticipant } = participant;
  const { contacts, users, email, rawEmail } = emailParticipant;

  if (contacts.length) {
    const label = contacts.find((c) => c?.name || c?.firstName || c?.lastName);

    return {
      label: label
        ? label?.name || `${label?.firstName} ${label?.lastName}`.trim()
        : '',
      [keyName]: email || rawEmail || '',
    };
  }
  if (users.length) {
    const label = users.find((c) => c?.firstName || c?.lastName);

    return {
      label: label ? `${label?.firstName} ${label?.lastName}`.trim() : '',
      [keyName]: email || rawEmail || '',
    };
  }

  return {
    label: '',
    [keyName]: email || rawEmail || '',
  };
};

export const getEmailParticipantsNameAndEmail = (
  participants: InteractionEventParticipant[],
  label = 'email',
): Array<{ [x: string]: string; label: string }> => {
  return participants?.map((participant) =>
    getParticipantNameAndEmail(participant as EmailParticipant, label),
  );
};
