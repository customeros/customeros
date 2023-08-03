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

const Header: FC<any> = ({ context: { loadMore, loading } }) => {
  return (
    <Button
      variant='outline'
      loadingText='Loading'
      isLoading={loading}
      borderColor='gray.200'
      color='gray.500'
      mt={4}
      onClick={loadMore}
      isDisabled={loading}
      _hover={{
        background: 'primary.50',
        color: 'primary.700',
        borderColor: 'primary.200',
      }}
      _focus={{
        background: 'primary.50',
        color: 'primary.700',
        borderColor: 'primary.200',
      }}
      _focusVisible={{
        background: 'primary.50',
        color: 'primary.700',
        borderColor: 'primary.200',
        boxShadow: '0 0 0 4px var(--chakra-colors-primary-100)',
      }}
    >
      Load more
    </Button>
  );
};

export const OrganizationTimeline: FC = () => {
  const id = useParams()?.id as string;
  const newDate = useMemo(() => new Date(), [id]);
  const virtuoso = useRef(null);

  const client = getGraphQLClient();
  const { data, isInitialLoading, isLoading } = useGetTimelineQuery(client, {
    organizationId: id,
    from: newDate,
    size: 100,
  });

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
    ?.filter((d: InteractionEvent) => !!d?.id && d.channel === 'EMAIL')
    ?.reverse();

  if (!timelineEmailEvents?.length) {
    return (
      <Flex direction='column' height='100%'>
        <EmptyTimeline />
        <Flex bg='#F9F9FB' direction='column' flex={1} pl={6}>
          <div>
            <TimelineActions
              // @ts-expect-error shouldn't cause error
              onScrollBottom={() => virtuoso?.current?.scrollBy({ top: 300 })}
            />
          </div>

          <Flex flex={1} height='100%' bg='#F9F9FB' />
        </Flex>
      </Flex>
    );
  }
  console.log('🏷️ ----- timelineEmailEvents: ', timelineEmailEvents);
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
              <EmailStub email={timelineEvent as unknown as InteractionEvent} />
            </TimelineItem>
          );
        }}
        components={{
          Footer: () => (
            <TimelineActions
              // @ts-expect-error shouldn't cause error
              onScrollBottom={() => virtuoso?.current?.scrollBy({ top: 300 })}
            />
          ),
        }}
      />

      <EmailPreviewModal />
    </TimelineEventPreviewContextContextProvider>
  );
};
