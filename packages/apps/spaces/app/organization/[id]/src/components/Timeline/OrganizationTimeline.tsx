'use client';
import React, { FC, useCallback, useEffect, useMemo } from 'react';
import { setHours, setSeconds, setMinutes, setMilliseconds } from 'date-fns';
import { useIsRestoring } from '@tanstack/react-query';
import { DateTimeUtils } from '@spaces/utils/date';
import { Virtuoso } from 'react-virtuoso';
import { EmailStub, TimelineItem } from './events';
import { useInfiniteGetTimelineQuery } from '../../graphql/getTimeline.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useParams } from 'next/navigation';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';
import { Button } from '@ui/form/Button';
import { Flex } from '@ui/layout/Flex';
import { EmptyTimeline } from '@organization/src/components/Timeline/EmptyTimeline';
import { TimelineItemSkeleton } from '@organization/src/components/Timeline/events/TimelineItem/TimelineItemSkeleton';
import { TimelineActions } from '@organization/src/components/Timeline/TimelineActions/TimelineActions';
import { useQueryClient } from '@tanstack/react-query';
import { SlackStub } from '@organization/src/components/Timeline/events/slack/SlackStub';
import { MeetingStub } from './events/meeting/MeetingStub';
import { TimelineEventPreviewModal } from '@organization/src/components/Timeline/preview/TimelineEventPreviewModal';
import {
  InteractionEventWithDate,
  LogEntryWithAliases,
  TimelineEvent,
} from './types';
import { UserActionStub } from '@organization/src/components/Timeline/events/action/UserActionStub';
import { IntercomStub } from '@organization/src/components/Timeline/events/intercom/IntercomStub';
import { ExternalSystemType } from '@graphql/types';
import { LogEntryStub } from '@organization/src/components/Timeline/events/logEntry/LogEntryStub';

import { useTimelineMeta } from './shared/state';
import { IssueStub } from '@organization/src/components/Timeline/events/issue/IssueStub';

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

export const NEW_DATE = setSeconds(
  setMilliseconds(
    setMinutes(
      setHours(
        new Date(new Date(new Date().setDate(new Date().getDate() + 1))),
        0,
      ),
      0,
    ),
    0,
  ),
  0,
);

function getEventDate(event?: TimelineEvent) {
  return (
    (event as InteractionEventWithDate)?.date ||
    (event as LogEntryWithAliases)?.logEntryStartedAt ||
    event?.createdAt
  );
}
export const OrganizationTimeline: FC = () => {
  const styles = useMemo(
    () => ({ height: '100%', width: '100%', background: '#F9F9FB' }),
    [],
  );
  const id = useParams()?.id as string;
  const queryClient = useQueryClient();
  const { virtuosoRef } = useTimelineRefContext();
  const [timelineMeta, setTimelineMeta] = useTimelineMeta();
  const isRestoring = useIsRestoring();
  const client = getGraphQLClient();
  const { data, isFetchingNextPage, fetchNextPage } =
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
          const lastEventDate = getEventDate(lastEvent as TimelineEvent);
          return {
            from: lastEvent ? lastEventDate : new Date(),
          };
        },
      },
    );
  const invalidateQuery = useCallback(() => {
    queryClient.invalidateQueries(['GetTimeline.infinite']);
  }, []);

  useEffect(() => {
    setTimelineMeta({
      ...timelineMeta,
      getTimelineVariables: {
        organizationId: id,
        from: NEW_DATE.toISOString(),
        size: 50,
      },
    });
  }, [NEW_DATE, id]);

  const virtuosoContext = useMemo(
    () => ({
      loadMore: () => fetchNextPage(),
      loading: isFetchingNextPage,
    }),
    [fetchNextPage, isFetchingNextPage],
  );
  const Footer = useCallback(() => {
    return <TimelineActions invalidateQuery={invalidateQuery} />;
  }, [invalidateQuery]);

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
        case 'Issue':
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
          case 'Issue':
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

  if (!isRestoring && !timelineEmailEvents?.length) {
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

      <Virtuoso<TimelineEvent>
        ref={virtuosoRef}
        style={styles}
        initialItemCount={timelineEmailEvents?.length}
        initialTopMostItemIndex={timelineEmailEvents?.length - 1}
        data={timelineEmailEvents}
        increaseViewportBy={300}
        atTopThreshold={100}
        context={virtuosoContext}
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
            case 'Issue': {
              return (
                <TimelineItem
                  date={timelineEvent?.createdAt}
                  showDate={showDate}
                >
                  <IssueStub data={timelineEvent} />
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
          Footer: Footer,
        }}
      />
      <TimelineEventPreviewModal invalidateQuery={invalidateQuery} />
    </>
  );
};
