'use client';
import React, { FC, useMemo, useRef } from 'react';
import { DateTimeUtils } from '@spaces/utils/date';
import { Virtuoso } from 'react-virtuoso';
import { EmailStub, TimelineItem } from './events';
import { useGetTimelineQuery } from '../../graphql/getTimeline.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useParams } from 'next/navigation';
import { InteractionEvent } from '@graphql/types';
import { EmailPreviewModal } from './events/email/EmailPreviewModal';
import { TimelineEventPreviewContextContextProvider } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { Button } from '@ui/form/Button';
import { Flex } from '@ui/layout/Flex';
import { EmptyTimeline } from '@organization/components/Timeline/EmptyTimeline';
import { TimelineItemSkeleton } from '@organization/components/Timeline/events/TimelineItem/TimelineItemSkeleton';
import { TimelineActions } from '@organization/components/Timeline/TimelineActions/TimelineActions';
import { useQueryClient } from '@tanstack/react-query';
import { SlackStub } from '@organization/components/Timeline/events/slack/SlackStub';

const Header: FC<any> = ({ context: { loadMore, loading } }) => {
  return (
    <Button
      variant='outline'
      loadingText='Loading'
      isLoading={loading}
      mt={4}
      onClick={loadMore}
      isDisabled={loading}
    >
      Load more
    </Button>
  );
};

const NEW_DATE = new Date();

export const OrganizationTimeline: FC = () => {
  const id = useParams()?.id as string;
  const virtuoso = useRef(null);
  const queryClient = useQueryClient();

  const client = getGraphQLClient();
  const { data, isInitialLoading, isLoading } = useGetTimelineQuery(client, {
    organizationId: id,
    from: NEW_DATE,
    size: 100,
  });
  const invalidateQuery = () =>
    queryClient.invalidateQueries(
      useGetTimelineQuery.getKey({
        organizationId: id,
        from: NEW_DATE,
        size: 100,
      }),
    );

  if (isInitialLoading) {
    return (
      <Flex direction='column' mt={4} pl={6}>
        <TimelineItemSkeleton />
        <TimelineItemSkeleton />
        <TimelineItemSkeleton />
      </Flex>
    );
  }

  const timelineEmailEvents = (
    data?.organization?.timelineEvents as unknown as InteractionEvent[]
  )
    ?.filter(
      (d: InteractionEvent) =>
        !!d?.id && (d.channel === 'EMAIL' || d.channel === 'SLACK'),
    )
    ?.reverse();
  console.log('🏷️ ----- timelineEmailEvents: ', timelineEmailEvents);
  if (!timelineEmailEvents?.length) {
    return <EmptyTimeline invalidateQuery={invalidateQuery} />;
  }

  return (
    <TimelineEventPreviewContextContextProvider
      data={timelineEmailEvents || []}
    >
      <Virtuoso
        ref={virtuoso}
        style={{ height: '100%', width: '100%', background: '#F9F9FB' }}
        initialItemCount={timelineEmailEvents?.length}
        initialTopMostItemIndex={timelineEmailEvents.length - 1}
        data={timelineEmailEvents}
        increaseViewportBy={300}
        atTopThreshold={100}
        // context={{ loadMore: () => null, loading: isLoading }}
        // components={{ Header }}
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
              {timelineEvent.channel === 'EMAIL' && (
                <EmailStub
                  email={timelineEvent as unknown as InteractionEvent}
                />
              )}
              {timelineEvent.channel === 'SLACK' && (
                <SlackStub
                  slackEvent={timelineEvent as unknown as InteractionEvent}
                />
              )}
            </TimelineItem>
          );
        }}
        components={{
          Footer: () => (
            <TimelineActions
              invalidateQuery={invalidateQuery}
              // @ts-expect-error shouldn't cause error
              onScrollBottom={() => virtuoso?.current?.scrollBy({ top: 300 })}
            />
          ),
        }}
      />

      <EmailPreviewModal invalidateQuery={invalidateQuery} />
    </TimelineEventPreviewContextContextProvider>
  );
};
