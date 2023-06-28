import {
  ContactParticipant,
  EmailParticipant,
  PhoneNumberParticipant,
  UserParticipant,
} from '@spaces/graphql';

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
