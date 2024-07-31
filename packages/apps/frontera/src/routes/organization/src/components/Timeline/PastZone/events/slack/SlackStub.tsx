import { FC } from 'react';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { User02 } from '@ui/media/icons/User02';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { InteractionEventWithDate } from '@organization/components/Timeline/types';
import { SlackMessageCard } from '@organization/components/Timeline/PastZone/events/slack/SlackMessageCard';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { getDisplayNameAndAvatar } from './util';

export const SlackStub: FC<{ slackEvent: InteractionEventWithDate }> = ({
  slackEvent,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();

  const slackSender = getDisplayNameAndAvatar(slackEvent?.sentBy?.[0]);

  const isSentByTenantUser =
    slackEvent?.sentBy?.[0]?.__typename === 'UserParticipant';
  const slackEventReplies = slackEvent.interactionSession?.events?.filter(
    (e) => e?.id !== slackEvent?.id,
  );
  const uniqThreadParticipants = slackEventReplies
    ?.map((e) => {
      return getDisplayNameAndAvatar(e.sentBy?.[0]);
    })
    ?.filter((v, i, a) => a.findIndex((t) => !!t && t?.id === v?.id) === i);

  return (
    <SlackMessageCard
      showDateOnHover
      content={slackEvent?.content || ''}
      onClick={() => openModal(slackEvent.id)}
      name={slackSender.displayName || 'Unknown'}
      date={DateTimeUtils.formatTime(slackEvent?.date)}
      profilePhotoUrl={slackSender.photoUrl || undefined}
      className={cn(isSentByTenantUser ? 'ml-6' : 'ml-0')}
      sourceUrl={slackEvent?.externalLinks?.[0]?.externalUrl}
    >
      {!!slackEventReplies?.length && (
        <div className='flex mt-1'>
          <div className='flex gap-1 mr-1'>
            {uniqThreadParticipants?.map(({ id, displayName, photoUrl }) => {
              return (
                <Avatar
                  size='xs'
                  name={displayName}
                  src={photoUrl ?? undefined}
                  variant='roundedSquareSmall'
                  icon={<User02 className='text-primary-700' />}
                  key={`uniq-slack-thread-participant-${slackEvent.id}-${id}`}
                />
              );
            })}
          </div>
          <Button size='sm' variant='link' className='text-sm'>
            {slackEventReplies.length}{' '}
            {slackEventReplies.length === 1 ? 'reply' : 'replies'}
          </Button>
        </div>
      )}
    </SlackMessageCard>
  );
};
