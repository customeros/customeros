'use client';
import React, { FC } from 'react';
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

const NEW_DATE = new Date();
//
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

  const client = getGraphQLClient();
  const { data, isInitialLoading, isLoading } = useGetTimelineQuery(client, {
    organizationId: id,
    from: NEW_DATE,
    size: 100,
  });

  if (isInitialLoading) {
    return (
      <Flex direction='column' mt={4}>
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
    return <EmptyTimeline />;
  }

  return (
    <TimelineEventPreviewContextContextProvider
      data={timelineEmailEvents || []}
    >
      <Virtuoso
        style={{ height: '100%', width: '100%' }}
        initialItemCount={timelineEmailEvents?.length}
        initialTopMostItemIndex={timelineEmailEvents.length - 1}
        data={timelineEmailEvents}
        increaseViewportBy={300}
        overscan={10}
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
      />

      <EmailPreviewModal />
    </TimelineEventPreviewContextContextProvider>
  );
};
