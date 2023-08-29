'use client';
import React, { FC, useRef } from 'react';
import { DateTimeUtils } from '@spaces/utils/date';
import { Virtuoso, VirtuosoHandle } from 'react-virtuoso';
import { EmailStub, TimelineItem } from './events';
import { useInfiniteGetTimelineQuery } from '../../graphql/getTimeline.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useParams } from 'next/navigation';
import { InteractionEvent, Meeting } from '@graphql/types';
import { TimelineEventPreviewContextContextProvider } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { Button } from '@ui/form/Button';
import { Flex } from '@ui/layout/Flex';
import { EmptyTimeline } from '@organization/components/Timeline/EmptyTimeline';
import { TimelineItemSkeleton } from '@organization/components/Timeline/events/TimelineItem/TimelineItemSkeleton';
import { TimelineActions } from '@organization/components/Timeline/TimelineActions/TimelineActions';
import { useQueryClient } from '@tanstack/react-query';
import { SlackStub } from '@organization/components/Timeline/events/slack/SlackStub';
import { MeetingStub } from './events/meeting/MeetingStub';
import { TimelineEventPreviewModal } from '@organization/components/Timeline/preview/TimelineEventPreviewModal';

type InteractionEventWithDate = InteractionEvent & { date: string };
type Event = InteractionEventWithDate | Meeting;

const Header: FC<{ context?: any }> = ({ context: { loadMore, loading } }) => {
  return (
    <Button
      variant='outline'
      colorScheme='primary'
      loadingText='Loading'
      isLoading={loading}
      mt={4}
      size='sm'
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
  const virtuoso = useRef<VirtuosoHandle>(null);
  const queryClient = useQueryClient();

  const client = getGraphQLClient();
  const { data, isInitialLoading, isFetchingNextPage, fetchNextPage } =
    useInfiniteGetTimelineQuery(
      'from',
      client,
      {
        organizationId: id,
        from: NEW_DATE,
        size: 50,
      },
      {
        getNextPageParam: (lastPage) => {
          const lastEvent = lastPage?.organization?.timelineEvents?.slice(
            -1,
          )?.[0] as InteractionEventWithDate;
          return {
            from: lastEvent ? lastEvent.date : new Date(),
          };
        },
      },
    );
  const invalidateQuery = () =>
    queryClient.invalidateQueries(
      useInfiniteGetTimelineQuery.getKey({
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

  const flattenData = data?.pages.flatMap(
    (page) => page?.organization?.timelineEvents,
  ) as unknown as Event[];

  const loadedDataCount = data?.pages.flatMap(
    (page) => page?.organization?.timelineEvents,
  )?.length;

  const timelineEmailEvents = flattenData
    ?.filter((d) => {
      if (!d) return false;
      switch (d.__typename) {
        case 'InteractionEvent':
          return !!d?.id && ['EMAIL', 'SLACK'].includes(d.channel ?? '');
        case 'Meeting':
          return !!d.id;
        default:
          return false;
      }
    })
    .sort((a, b) => {
      const getDate = (a: Event) => {
        if (!a) return null;
        switch (a.__typename) {
          case 'InteractionEvent':
            return a.date;
          case 'Meeting':
            return a.createdAt;
          default:
            return null;
        }
      };
      const aDate = getDate(a);
      const bDate = getDate(b);

      return Date.parse(aDate) - Date.parse(bDate);
    });

  if (!timelineEmailEvents?.length) {
    return <EmptyTimeline invalidateQuery={invalidateQuery} />;
  }

  return (
    <>
      {isFetchingNextPage && (
        <Flex direction='column' mt={4} pl={6}>
          <TimelineItemSkeleton />
          <TimelineItemSkeleton />
        </Flex>
      )}
      <TimelineEventPreviewContextContextProvider
        data={timelineEmailEvents || []}
        id={id}
      >
        <Virtuoso<Event>
          ref={virtuoso}
          style={{ height: '100%', width: '100%', background: '#F9F9FB' }}
          initialItemCount={timelineEmailEvents?.length}
          initialTopMostItemIndex={timelineEmailEvents.length - 1}
          data={timelineEmailEvents}
          increaseViewportBy={300}
          atTopThreshold={100}
          context={{
            loadMore: () => fetchNextPage(),
            loading: isFetchingNextPage,
          }}
          itemContent={(index, timelineEvent) => {
            if (!timelineEvent) return null;
            switch (timelineEvent.__typename) {
              case 'InteractionEvent': {
                const showDate =
                  index === 0
                    ? true
                    : !DateTimeUtils.isSameDay(
                        (
                          timelineEmailEvents?.[
                            index - 1
                          ] as InteractionEventWithDate
                        )?.date,
                        timelineEvent.date,
                      );

                return (
                  <TimelineItem date={timelineEvent?.date} showDate={showDate}>
                    {timelineEvent.channel === 'EMAIL' && (
                      <EmailStub
                        email={timelineEvent as unknown as InteractionEvent}
                      />
                    )}
                    {timelineEvent.channel === 'SLACK' && (
                      <SlackStub
                        slackEvent={
                          timelineEvent as unknown as InteractionEvent
                        }
                      />
                    )}
                  </TimelineItem>
                );
              }
              case 'Meeting': {
                return (
                  <TimelineItem date={timelineEvent?.createdAt} showDate>
                    <MeetingStub data={timelineEvent} />
                  </TimelineItem>
                );
              }
              default:
                return null;
            }
          }}
          components={{
            Header: (rest) => (
              <Flex bg='gray.25' p={5}>
                {loadedDataCount &&
                !isFetchingNextPage &&
                data?.pages?.[0]?.organization?.timelineEventsTotalCount >
                  loadedDataCount ? (
                  <Header {...rest} />
                ) : null}
              </Flex>
            ),
            Footer: () => (
              <TimelineActions
                invalidateQuery={invalidateQuery}
                onScrollBottom={() => virtuoso?.current?.scrollBy({ top: 300 })}
              />
            ),
          }}
        />
        <TimelineEventPreviewModal invalidateQuery={invalidateQuery} />
      </TimelineEventPreviewContextContextProvider>
    </>
  );
};
