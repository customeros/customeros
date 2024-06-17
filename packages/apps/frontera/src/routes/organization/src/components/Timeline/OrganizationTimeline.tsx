import { Virtuoso } from 'react-virtuoso';
import { useParams } from 'react-router-dom';
import { FC, useMemo, useEffect, useCallback } from 'react';

import { observer } from 'mobx-react-lite';
import { useQueryClient } from '@tanstack/react-query';
import { setHours, setSeconds, setMinutes, setMilliseconds } from 'date-fns';

import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Meeting, ExternalSystemType } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { EmptyTimeline } from '@organization/components/Timeline/EmptyTimeline';
import { SlackStub } from '@organization/components/Timeline/PastZone/events/slack/SlackStub';
import { IssueStub } from '@organization/components/Timeline/PastZone/events/issue/IssueStub';
import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';
import { IntercomStub } from '@organization/components/Timeline/PastZone/events/intercom/IntercomStub';
import { LogEntryStub } from '@organization/components/Timeline/PastZone/events/logEntry/LogEntryStub';
import { UserActionStub } from '@organization/components/Timeline/PastZone/events/action/UserActionStub';
import { TimelineActions } from '@organization/components/Timeline/FutureZone/TimelineActions/TimelineActions';
import { TimelineItemSkeleton } from '@organization/components/Timeline/PastZone/events/TimelineItem/TimelineItemSkeleton';
import { TimelineEventPreviewModal } from '@organization/components/Timeline/shared/TimelineEventPreview/TimelineEventPreviewModal';

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
      className='mt-4'
      variant='outline'
      colorScheme='primary'
      loadingText='Loading'
      isLoading={loading}
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

export const OrganizationTimeline = observer(() => {
  const store = useStore();

  const styles = useMemo(
    () => ({ height: '100%', width: '100%', background: '#F9F9FB' }),
    [],
  );
  const id = useParams()?.id as string;
  const queryClient = useQueryClient();
  const { virtuosoRef } = useTimelineRefContext();
  const [timelineMeta, setTimelineMeta] = useTimelineMeta();
  const client = getGraphQLClient();

  const timeline =
    store.timelineEvents.getByOrganizationId(id)?.map((t) => t.value) ?? [];

  const { data, isFetchingNextPage, fetchNextPage, isPending } =
    useInfiniteGetTimelineQuery(
      client,
      {
        organizationId: id,
        from: NEW_DATE.toISOString(),
        size: store.demoMode ? 100 : 50,
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
  const timelineItemsLength = Object.values(data?.pages ?? []).reduce(
    (acc, curr) => curr.organization?.timelineEventsTotalCount + acc,
    0,
  );

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

  useEffect(() => {
    setTimelineMeta((prev) => ({
      ...prev,
      itemCount: timelineItemsLength,
    }));
  }, [timelineItemsLength]);

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

  const loadedDataCount = store.demoMode ? timeline.length : flattenData.length;

  const timelineEmailEvents = (store.demoMode ? timeline : flattenData)
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
      const aDate = getDate(a as TimelineEvent);
      const bDate = getDate(b as TimelineEvent);

      return Date.parse(aDate) - Date.parse(bDate);
    });

  useEffect(() => {
    store.timelineEvents.bootstrapTimeline(id);
  }, [id]);

  if (!isPending && !timelineEmailEvents?.length) {
    return (
      <>
        <EmptyTimeline invalidateQuery={invalidateQuery} />
        <TimelineEventPreviewModal invalidateQuery={invalidateQuery} />
      </>
    );
  }

  if (isPending && !isFetchingNextPage) {
    return (
      <div className='flex flex-col mt-4 pl-6 w-full'>
        <TimelineItemSkeleton />
        <TimelineItemSkeleton />
      </div>
    );
  }

  return (
    <>
      {isFetchingNextPage && (
        <div className='flex flex-col mt-4 pl-6 w-full'>
          <TimelineItemSkeleton />
          <TimelineItemSkeleton />
        </div>
      )}

      <Virtuoso<TimelineEvent>
        ref={virtuosoRef}
        style={styles}
        initialTopMostItemIndex={{
          align: 'start',
          index: timelineEmailEvents?.length - 1,
          behavior: 'auto',
        }}
        data={(timelineEmailEvents as TimelineEvent[]) ?? []}
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
            <div className='flex bg-gray-25 p-5'>
              {loadedDataCount &&
              !isFetchingNextPage &&
              data?.pages?.[0]?.organization?.timelineEventsTotalCount >
                loadedDataCount ? (
                <Header {...rest} />
              ) : null}
            </div>
          ),
          Footer: Footer,
        }}
      />
      <TimelineEventPreviewModal invalidateQuery={invalidateQuery} />
    </>
  );
});
