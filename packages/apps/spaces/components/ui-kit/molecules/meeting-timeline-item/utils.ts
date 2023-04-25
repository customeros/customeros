import {
  Contact,
  ContactParticipant,
  MeetingParticipant,
  User,
  UserParticipant,
} from '../../../../graphQL/__generated__/generated';

export const getAttendeeDataFromParticipant = (
  participant: MeetingParticipant,
): Contact | User => {
  if (
    participant.__typename !== 'ContactParticipant' &&
    participant.__typename !== 'UserParticipant'
  ) {
    throw new Error(
      'Meeting participant type error. Participant is neither contact nor user',
    );
  }

  return participant.__typename === 'ContactParticipant'
    ? participant.contactParticipant
    : (participant as UserParticipant).userParticipant;
};
