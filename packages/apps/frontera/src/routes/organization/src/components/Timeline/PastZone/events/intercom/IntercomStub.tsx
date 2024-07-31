import { FC } from 'react';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { User02 } from '@ui/media/icons/User02';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { getName } from '@utils/getParticipantsName';
import { InteractionEventWithDate } from '@organization/components/Timeline/types';
import {
  UserParticipant,
  ContactParticipant,
  JobRoleParticipant,
} from '@graphql/types';
import { IntercomMessageCard } from '@organization/components/Timeline/PastZone/events/intercom/IntercomMessageCard';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

// TODO unify with slack
export const IntercomStub: FC<{ intercomEvent: InteractionEventWithDate }> = ({
  intercomEvent,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();

  const intercomSender =
    (intercomEvent?.sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (intercomEvent?.sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant
      ?.contact ||
    (intercomEvent?.sentBy?.[0] as UserParticipant)?.userParticipant;
  const isSentByTenantUser =
    intercomEvent?.sentBy?.[0]?.__typename === 'UserParticipant';

  if (!intercomSender) {
    return null;
  }

  const intercomEventReplies = intercomEvent.interactionSession?.events?.filter(
    (e) => e?.id !== intercomEvent?.id,
  );
  const uniqThreadParticipants = intercomEventReplies
    ?.map((e) => {
      const threadSender =
        (e?.sentBy?.[0] as ContactParticipant)?.contactParticipant ||
        (e?.sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant?.contact ||
        (e?.sentBy?.[0] as UserParticipant)?.userParticipant;

      return threadSender;
    })
    ?.filter((v, i, a) => a.findIndex((t) => !!t && t?.id === v?.id) === i);

  return (
    <IntercomMessageCard
      showDateOnHover
      name={getName(intercomSender)}
      content={intercomEvent?.content || ''}
      onClick={() => openModal(intercomEvent.id)}
      profilePhotoUrl={intercomSender?.profilePhotoUrl}
      date={DateTimeUtils.formatTime(intercomEvent?.date)}
      className={cn(isSentByTenantUser ? 'ml-6' : 'ml-0')}
      sourceUrl={intercomEvent?.externalLinks?.[0]?.externalUrl}
    >
      {!!intercomEventReplies?.length && (
        <div className='flex mt-1'>
          <div className='flex gap-1 mr=1'>
            {uniqThreadParticipants?.map(
              ({ id, name, firstName, lastName, profilePhotoUrl }) => {
                const displayName =
                  name ?? [firstName, lastName].filter(Boolean).join(' ');

                return (
                  <Avatar
                    size='xs'
                    name={displayName}
                    variant='roundedSquareSmall'
                    icon={<User02 color='primary.700' />}
                    src={profilePhotoUrl ? profilePhotoUrl : undefined}
                    key={`uniq-intercom-thread-participant-${intercomEvent.id}-${id}`}
                  />
                );
              },
            )}
          </div>
          <Button size='sm' variant='link' className='text-sm'>
            {intercomEventReplies.length}{' '}
            {intercomEventReplies.length === 1 ? 'reply' : 'replies'}
          </Button>
        </div>
      )}
    </IntercomMessageCard>
  );
};
