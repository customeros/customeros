import {
  ContactParticipant,
  EmailParticipant,
  JobRoleParticipant,
  PhoneNumberParticipant,
  UserParticipant,
} from '@graphql/types';
import { Contact, InteractionEventParticipant, User } from '@graphql/types';

type Participant =
  | EmailParticipant
  | PhoneNumberParticipant
  | ContactParticipant
  | UserParticipant;

export const getParticipantNames = (participants: Participant[]): string[] => {
  return participants.map((participant) => {
    if (participant.__typename === 'EmailParticipant') {
      const { emailParticipant } = participant;
      const { contacts, users } = emailParticipant;
      if (contacts.length) {
        return contacts
          .map((c) => (c?.name ? c.name : `${c.firstName} ${c.lastName}`))
          .join(' ');
      }
      if (users.length) {
        return users.map((c) => `${c.firstName} ${c.lastName}`).join(' ');
      }

      const participantName =
        contacts?.[0]?.name ||
        users?.[0]?.firstName + ' ' + users?.[0]?.lastName;
      return participantName || 'Unnamed';
    } else if (participant.__typename === 'PhoneNumberParticipant') {
      const { phoneNumberParticipant } = participant;
      const { contacts, users } = phoneNumberParticipant;
      const participantName =
        contacts?.[0]?.name ||
        users?.[0]?.firstName + ' ' + users?.[0]?.lastName;
      return participantName || 'name';
    } else if (participant.__typename === 'ContactParticipant') {
      const { contactParticipant } = participant;
      const { name, firstName, lastName } = contactParticipant;
      return firstName + ' ' + lastName || name || 'Unnamed';
    } else if (participant.__typename === 'UserParticipant') {
      const { userParticipant } = participant;
      const { firstName, lastName } = userParticipant;
      return firstName + ' ' + lastName || 'Unnamed';
    }
    return 'Unnamed';
  });
};

export const getName = (
  data: Contact | User,
  email?: string | null | undefined,
  rawEmail?: string | undefined | null,
): string => {
  if (data.__typename === 'Contact' && data?.name) {
    return data.name;
  }
  if (data?.firstName || data?.lastName) {
    return `${data.firstName} ${data.lastName}`;
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

export const getParticipantName = (
  participant: ContactParticipant | UserParticipant | JobRoleParticipant,
): string => {
  let participantT;
  if ((participant as ContactParticipant)?.contactParticipant) {
    participantT = (participant as ContactParticipant).contactParticipant;
  }
  if ((participant as UserParticipant)?.userParticipant) {
    participantT = (participant as UserParticipant).userParticipant;
  }
  if ((participant as JobRoleParticipant)?.jobRoleParticipant?.contact) {
    participantT = (participant as JobRoleParticipant).jobRoleParticipant
      ?.contact;
  }
  if (!participantT) return '';

  return getName(participantT);
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
      label: label?.name || `${label?.firstName} ${label?.lastName}`.trim(),
      [keyName]: email || rawEmail || '',
    };
  }
  if (users.length) {
    const label = users.find((c) => c?.firstName || c?.lastName);

    return {
      label: `${label?.firstName} ${label?.lastName}`.trim(),
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
