'use client';
import React, { FC } from 'react';
import { DateTimeUtils } from '@spaces/utils/date';
import { Virtuoso } from 'react-virtuoso';
import { EmailTimelineItem, TimelineItem } from './events';
import { useGetTimelineQuery } from '../../graphql/getTimeline.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useParams } from 'next/navigation';
import { InteractionEvent } from '@graphql/types';
import { EmailPreviewModal } from './events/email/EmailPreviewModal';
import { TimelineEventPreviewContextContextProvider } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';

const NEW_DATE = new Date();
export const OrganizationTimeline: FC = () => {
  const id = useParams()?.id as string;

  const client = getGraphQLClient();
  const { data, isInitialLoading} = useGetTimelineQuery(client, {
    organizationId: id,
    from: NEW_DATE,
    size: 500,
  });

  if (isInitialLoading) {
    // Todo
    return <div></div>;
  }

  const timelineEmailEvents = (
    data?.organization?.timelineEvents as unknown as InteractionEvent[]
  )
    ?.filter((d: InteractionEvent) => !!d?.id && d.channel === 'EMAIL')
    ?.reverse();

  if (!timelineEmailEvents?.length) {
    return null;
  }

  return (
    <TimelineEventPreviewContextContextProvider data={timelineEmailEvents || []}>
      <Virtuoso
        style={{ height: '100%', width: '100%' }}
        initialItemCount={timelineEmailEvents?.length}
        initialTopMostItemIndex={timelineEmailEvents.length - 1}
        data={timelineEmailEvents}
        increaseViewportBy={300}
        overscan={10}
        atTopThreshold={100}
        itemContent={(index, timelineEvent: InteractionEvent) => {
          if (timelineEvent.__typename !== 'InteractionEvent') return null;

          const showDate =
            index === 0
              ? true
              : !DateTimeUtils.isSameDay(
                  // @ts-expect-error this is correct, generated types did not picked up alias correctly
                  timelineEmailEvents?.[index - 1]?.date,
                  // @ts-expect-error this is correct, generated types did not picked up alias correctly
                  timelineEvent.date,
                );
          return (
            // @ts-expect-error this is correct, generated types did not picked up alias correctly
            <TimelineItem date={timelineEvent?.date} showDate={showDate}>
              <EmailTimelineItem
                email={timelineEvent as unknown as InteractionEvent}
              />
            </TimelineItem>
          );
        }}
      />

      <EmailPreviewModal />
    </TimelineEventPreviewContextContextProvider>
  );
};
