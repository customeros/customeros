import {
  ContactParticipant,
  JobRoleParticipant,
  Meeting,
  OrganizationParticipant,
  UserParticipant,
} from '@graphql/types';
import { getName } from '@spaces/utils/getParticipantsName';

export const getParticipantEmailOrName = (
  participant:
    | ContactParticipant
    | UserParticipant
    | JobRoleParticipant
    | OrganizationParticipant,
): string => {
  let participantT;

  switch (participant?.__typename) {
    case 'ContactParticipant':
      participantT = participant.contactParticipant;
      break;
    case 'UserParticipant':
      participantT = participant.userParticipant;
      break;
    case 'JobRoleParticipant':
      participantT = participant.jobRoleParticipant?.contact;
      break;
    case 'OrganizationParticipant':
      participantT = participant.organizationParticipant;
      break;
    default:
      participantT = null;
  }

  if (!participantT) return '';

  return participantT.emails?.[0]?.email || getName(participantT);
};

export const getParticipants = (
  data: Meeting | undefined,
): (string | string[] | number)[] => {
  if (data?.attendedBy?.length) {
    const fullArray = data?.attendedBy
      ?.map((participant) => getParticipantEmailOrName(participant))
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
