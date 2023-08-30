import {
  ContactParticipant,
  JobRoleParticipant,
  Meeting,
  OrganizationParticipant,
  UserParticipant,
} from '@graphql/types';
import { getName } from '@spaces/utils/getParticipantsName';


export const getParticipantName = (
  participant:
    | ContactParticipant
    | UserParticipant
    | JobRoleParticipant
    | OrganizationParticipant,
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
  if ((participant as OrganizationParticipant)?.organizationParticipant) {
    participantT = (participant as OrganizationParticipant)
      .organizationParticipant;
  }
  if (!participantT) return '';

  return getName(participantT);
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
