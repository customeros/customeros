import {
  ContactParticipant,
  JobRoleParticipant,
  Meeting,
  OrganizationParticipant,
  EmailParticipant,
  UserParticipant,
  Contact,
} from '@graphql/types';
import { getName } from '@spaces/utils/getParticipantsName';

type Participant =
  | ContactParticipant
  | UserParticipant
  | JobRoleParticipant
  | OrganizationParticipant
  | EmailParticipant;

export const getParticipant = (participant: Participant) => {
  switch (participant?.__typename) {
    case 'ContactParticipant':
      return participant.contactParticipant;
    case 'UserParticipant':
      return participant.userParticipant;
    case 'JobRoleParticipant':
      return participant.jobRoleParticipant?.contact;
    case 'OrganizationParticipant':
      return participant.organizationParticipant;
    case 'EmailParticipant':
      return participant.emailParticipant;
    default:
      return null;
  }
};

export const getParticipantName = (participant: Participant): string => {
  const contact = getParticipant(participant);

  if (!contact) return '';

  return getName(contact);
};

export const getParticipantEmail = (participant: Participant): string => {
  const contact = getParticipant(participant);

  if (!contact) return '';

  return getName(contact);
};

export const getParticipants = (
  data: Meeting | undefined,
): (string | string[] | number)[] => {
  const owner = data?.createdBy?.length
    ? getParticipantName(data?.createdBy?.[0])
    : null;

  if (data?.attendedBy?.length) {
    const fullArray = data?.attendedBy
      ?.map((participant) => getParticipantName(participant))
      .filter((p) => p !== owner)
      .filter(Boolean);

    if (!data?.note?.[0]?.content || !data?.agenda)
      return [fullArray.join(data.attendedBy.length > 2 ? ', ' : ' and '), ''];

    return fullArray
      .filter((_, i) => {
        return i < 2;
      })
      .join(data.attendedBy.length > 2 ? ', ' : ' and ')
      .concat(
        data.attendedBy.length > 2 ? ' + ' + (data.attendedBy.length - 2) : '',
      )
      .split(' + ');
  }
  return [];
};

export const getMentionOptionLabel = (
  contact: Pick<Contact, 'name' | 'firstName' | 'lastName' | 'emails'>,
): string | undefined | null => {
  return (
    contact?.name ||
    [contact?.firstName, contact?.lastName].filter(Boolean).join(' ') ||
    contact?.emails?.[0]?.email
  );
};
