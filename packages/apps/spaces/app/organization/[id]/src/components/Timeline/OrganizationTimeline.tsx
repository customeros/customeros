'use client';
import { Virtuoso } from 'react-virtuoso';
import { useParams } from 'next/navigation';
import React, { FC, useMemo, useEffect, useCallback } from 'react';

import { useIsRestoring, useQueryClient } from '@tanstack/react-query';
import { setHours, setSeconds, setMinutes, setMilliseconds } from 'date-fns';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { Meeting, ExternalSystemType } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { EmptyTimeline } from '@organization/src/components/Timeline/EmptyTimeline';
import { SlackStub } from '@organization/src/components/Timeline/PastZone/events/slack/SlackStub';
import { IssueStub } from '@organization/src/components/Timeline/PastZone/events/issue/IssueStub';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';
import { IntercomStub } from '@organization/src/components/Timeline/PastZone/events/intercom/IntercomStub';
import { LogEntryStub } from '@organization/src/components/Timeline/PastZone/events/logEntry/LogEntryStub';
import { UserActionStub } from '@organization/src/components/Timeline/PastZone/events/action/UserActionStub';
import { TimelineActions } from '@organization/src/components/Timeline/FutureZone/TimelineActions/TimelineActions';
import { TimelineItemSkeleton } from '@organization/src/components/Timeline/PastZone/events/TimelineItem/TimelineItemSkeleton';
import { TimelineEventPreviewModal } from '@organization/src/components/Timeline/shared/TimelineEventPreview/TimelineEventPreviewModal';

import { useTimelineMeta } from './state';
import { FutureZone } from './FutureZone/FutureZone';
import { EmailStub, TimelineItem } from './PastZone/events';
import { MeetingStub } from './PastZone/events/meeting/MeetingStub';
import { useInfiniteGetTimelineQuery } from '../../graphql/getTimeline.generated';
import {
  TimelineEvent,
  LogEntryWithAliases,
  InteractionEventWithDate,
} from './types';

// TODO: type this context accordingly
// eslint-disable-next-line @typescript-eslint/no-explicit-any
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
    (event as Meeting)?.createdAt
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
  const { data, isFetchingNextPage, fetchNextPage, isPending } =
    useInfiniteGetTimelineQuery(
      client,
      {
        organizationId: id,
        from: NEW_DATE.toISOString(),
        size: 50,
      },
      {
        initialPageParam: 0,
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
    queryClient.invalidateQueries({ queryKey: ['GetTimeline.infinite'] });
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
    return (
      <>
        <TimelineActions invalidateQuery={invalidateQuery} />
        <FutureZone />
      </>
    );
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
        case 'Issue':
          return !!d.id;
        case 'Action':
          return !!d.id && d.actionType !== 'CREATED';
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

  if (!isRestoring && !isPending && !timelineEmailEvents?.length) {
    return (
      <>
        <EmptyTimeline invalidateQuery={invalidateQuery} />
        <TimelineEventPreviewModal invalidateQuery={invalidateQuery} />
      </>
    );
  }

  return (
    <>
      {(isFetchingNextPage || isPending) && (
        <Flex direction='column' mt={4} pl={6}>
          <TimelineItemSkeleton />
          <TimelineItemSkeleton />
        </Flex>
      )}

      <Virtuoso<TimelineEvent>
        ref={virtuosoRef}
        style={styles}
        initialTopMostItemIndex={timelineEmailEvents?.length - 1}
        data={timelineEmailEvents ?? []}
        increaseViewportBy={300}
        atTopThreshold={100}
        context={virtuosoContext}
        itemContent={(index, timelineEvent) => {
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
              return <div>not supported</div>;
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
