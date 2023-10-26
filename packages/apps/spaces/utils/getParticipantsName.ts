import { EmailParticipant, Organization } from '@graphql/types';
import {
  Contact,
  InteractionEventParticipant,
  User,
  Email,
} from '@graphql/types';
// todo cleanup, move those to helper functions in timeline events utils folder
export const getName = (
  data: Contact | User | Organization | Email,
  email?: string | null | undefined,
  rawEmail?: string | undefined | null,
): string => {
  if ((data as Contact | Organization)?.name) {
    return <string>(data as Contact | Organization).name;
  }
  const personData: Contact | User = data as Contact | User;
  if (personData?.firstName || personData?.lastName) {
    return `${personData.firstName} ${personData.lastName}`;
  }
  if ((data as Email)?.rawEmail) {
    return getEmailParticipantName(data as Email);
  }
  return email || rawEmail || 'Unknown';
};

export const getEmailParticipantName = (emailParticipant: Email): string => {
  const { contacts, users, email, rawEmail } = emailParticipant;

  if (contacts.length) {
    return getName(contacts[0], email, rawEmail);
  }
  if (users.length) {
    return getName(users[0], email, rawEmail);
  }

  return email || rawEmail || '';
};

export const getEmailParticipantsName = (
  participants: EmailParticipant[],
): string => {
  return participants
    ?.map((participant) =>
      getEmailParticipantName(participant.emailParticipant),
    )
    .join(', ');
};

export const getParticipantNameAndEmail = (
  participant: EmailParticipant,
  keyName = 'email',
): { [x: string]: string; label: string } => {
  const { emailParticipant } = participant;
  const { contacts, users, email, rawEmail } = emailParticipant;

  if (contacts.length) {
    const label = contacts.find(
      (c) =>
        (c?.name && c.name.toLowerCase() !== 'unknown') ||
        (c?.firstName && c.firstName.toLowerCase() !== 'unknown') ||
        (c?.lastName && c.lastName.toLowerCase() !== 'unknown'),
    );

    return {
      label: label
        ? label?.name || `${label?.firstName} ${label?.lastName}`.trim()
        : '',
      [keyName]: email || rawEmail || '',
    };
  }
  if (users.length) {
    const label = users.find(
      (c) =>
        (c?.name && c.name.toLowerCase() !== 'unknown') ||
        (c?.firstName && c.firstName.toLowerCase() !== 'unknown') ||
        (c?.lastName && c.lastName.toLowerCase() !== 'unknown'),
    );

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
