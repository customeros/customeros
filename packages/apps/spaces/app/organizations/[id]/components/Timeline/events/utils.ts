import {
  ContactParticipant,
  JobRoleParticipant,
  Meeting,
  OrganizationParticipant,
  UserParticipant,
} from '@graphql/types';
import { getName } from '@spaces/utils/getParticipantsName';

type Participant =
  | ContactParticipant
  | UserParticipant
  | JobRoleParticipant
  | OrganizationParticipant;

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

  return contact.emails?.[0]?.email ?? getName(contact);
};

export const getParticipants = (
  data: Meeting | undefined,
): (string | string[] | number)[] => {
  if (data?.attendedBy?.length) {
    const fullArray = data?.attendedBy
      ?.map((participant) => getParticipantName(participant))
      .filter(Boolean);

    if (!data?.note?.[0]?.html || !data?.agenda)
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
