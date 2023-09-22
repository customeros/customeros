'use client';
import React, { FC, useCallback, useRef } from 'react';
import { DateTimeUtils } from '@spaces/utils/date';
import { Virtuoso, VirtuosoHandle } from 'react-virtuoso';
import { EmailStub, TimelineItem } from './events';
import { useInfiniteGetTimelineQuery } from '../../graphql/getTimeline.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useParams } from 'next/navigation';
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
import { InteractionEventWithDate, TimelineEvent } from './types';
import { UserActionStub } from '@organization/components/Timeline/events/action/UserActionStub';
import { IntercomStub } from '@organization/components/Timeline/events/intercom/IntercomStub';
import { ExternalSystemType } from '@spaces/graphql';
import { LogEntryStub } from '@organization/components/Timeline/events/logEntry/LogEntryStub';

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

export const NEW_DATE = new Date(new Date().setDate(new Date().getDate() + 1));

function getEventDate(event?: TimelineEvent) {
  return (event as InteractionEventWithDate)?.date || event?.createdAt;
}

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
  const invalidateQuery = useCallback(() => {
    queryClient.invalidateQueries(['GetTimeline.infinite']);
  }, []);

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
  ) as unknown as TimelineEvent[];

  const loadedDataCount = data?.pages.flatMap(
    (page) => page?.organization?.timelineEvents,
  )?.length;

  const timelineEmailEvents = flattenData
    ?.filter((d) => {
      if (!d) return false;
      switch (d.__typename) {
        case 'InteractionEvent':
          return (
            !!d?.id && ['EMAIL', 'SLACK', 'CHAT'].includes(d.channel ?? '')
          );
        case 'Meeting':
        case 'LogEntry':
        case 'Action':
          return !!d.id;
        default:
          return false;
      }
    })
    .sort((a, b) => {
      const getDate = (a: TimelineEvent) => {
        if (!a) return null;
        switch (a.__typename) {
          case 'InteractionEvent':
            return a.date;
          case 'Meeting':
          case 'Action':
            return a.createdAt;
          case 'LogEntry':
            return a.logEntryStartedAt;

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
        <Virtuoso<TimelineEvent>
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
            const showDate =
              index === 0
                ? true
                : !DateTimeUtils.isSameDay(
                    getEventDate(
                      timelineEmailEvents?.[
                        index - 1
                      ] as InteractionEventWithDate,
                    ),
                    getEventDate(timelineEvent as TimelineEvent),
                  );
            switch (timelineEvent.__typename) {
              case 'InteractionEvent': {
                return (
                  <TimelineItem date={timelineEvent?.date} showDate={showDate}>
                    {timelineEvent.channel === 'EMAIL' && (
                      <EmailStub email={timelineEvent} />
                    )}
                    {timelineEvent.channel === 'CHAT' && (
                      <>
                        {timelineEvent.externalLinks?.[0]?.type ===
                          ExternalSystemType.Slack && (
                          <SlackStub slackEvent={timelineEvent} />
                        )}
                        {timelineEvent.externalLinks?.[0]?.type ===
                          ExternalSystemType.Intercom && (
                          <IntercomStub intercomEvent={timelineEvent} />
                        )}
                      </>
                    )}
                  </TimelineItem>
                );
              }
              case 'Meeting': {
                return (
                  <TimelineItem
                    date={timelineEvent?.createdAt}
                    showDate={showDate}
                  >
                    <MeetingStub data={timelineEvent} />
                  </TimelineItem>
                );
              }
              case 'Action': {
                return (
                  <TimelineItem
                    date={timelineEvent?.createdAt}
                    showDate={showDate}
                  >
                    <UserActionStub data={timelineEvent} />
                  </TimelineItem>
                );
              }
              case 'LogEntry': {
                return (
                  <TimelineItem
                    date={timelineEvent?.logEntryStartedAt}
                    showDate={showDate}
                  >
                    <LogEntryStub data={timelineEvent} />
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
            Footer: () => {
              const memoizedScrollBy = useCallback(() => {
                virtuoso?.current?.scrollBy({ top: 300 });
              }, [virtuoso]);
              return (
                <TimelineActions
                  invalidateQuery={invalidateQuery}
                  onScrollBottom={memoizedScrollBy}
                />
              );
            },
          }}
        />
        <TimelineEventPreviewModal invalidateQuery={invalidateQuery} />
      </TimelineEventPreviewContextContextProvider>
    </>
  );
};
