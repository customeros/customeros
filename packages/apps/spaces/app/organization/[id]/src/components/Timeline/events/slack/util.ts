import { match } from 'ts-pattern';

import { InteractionEventParticipant } from '@graphql/types';

export const getDisplayNameAndAvatar = (
  eventParticipant?: InteractionEventParticipant,
) => {
  return match(eventParticipant)
    .with(
      { __typename: 'ContactParticipant' },
      ({
        contactParticipant: { id, name, firstName, lastName, profilePhotoUrl },
      }) => ({
        id,
        displayName: name ?? [firstName, lastName].filter(Boolean).join(' '),
        photoUrl: profilePhotoUrl,
      }),
    )
    .with(
      { __typename: 'JobRoleParticipant' },
      ({ jobRoleParticipant: { contact } }) => ({
        id: contact?.id,
        displayName:
          contact?.name ??
          [contact?.firstName, contact?.lastName].filter(Boolean).join(' '),
        photoUrl: contact?.profilePhotoUrl ?? '',
      }),
    )
    .with(
      { __typename: 'UserParticipant' },
      ({
        userParticipant: { id, name, firstName, lastName, profilePhotoUrl },
      }) => ({
        id,
        displayName: name ?? [firstName, lastName].filter(Boolean).join(' '),
        photoUrl: profilePhotoUrl ?? '',
      }),
    )
    .with(
      { __typename: 'EmailParticipant' },
      ({ emailParticipant: { users, contacts } }) =>
        users.length
          ? {
              id: users?.[0]?.id,
              displayName:
                users?.[0]?.name ??
                [users?.[0]?.firstName, users?.[0]?.lastName]
                  .filter(Boolean)
                  .join(' '),
              photoUrl: users?.[0]?.profilePhotoUrl ?? '',
            }
          : {
              id: contacts?.[0]?.id,
              displayName: contacts?.[0]?.name ?? '',
              photoUrl: contacts?.[0]?.profilePhotoUrl ?? '',
            },
    )
    .otherwise(() => ({
      id: null,
      displayName: '',
      photoUrl: '',
    }));
};
