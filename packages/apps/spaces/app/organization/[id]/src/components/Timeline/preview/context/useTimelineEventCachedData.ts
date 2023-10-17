import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import {
  GetTimelineEventsQuery,
  useGetTimelineEventsQuery,
} from '@organization/src/graphql/getTimelineEvents.generated';
import { QueryKey, useQueryClient } from '@tanstack/react-query';
import { TimelineEvent } from '@graphql/types';

export const useTimelineEventCachedData = () => {
  const [timelineMeta, _] = useTimelineMeta();
  const queryClient = useQueryClient();

  const getKeys = (timelineEventId: string) => {
    const singleEventQueryKey = useGetTimelineEventsQuery.getKey({
      ids: [timelineEventId],
    });

    const timelineQueryKey = useInfiniteGetTimelineQuery.getKey({
      ...timelineMeta.getTimelineVariables,
    });

    return [singleEventQueryKey, timelineQueryKey];
  };

  const getQueryData = (
    singleEventQueryKey: QueryKey,
    timelineQueryKey: QueryKey,
  ) => {
    const timelineEventsQueryCachedData = (
      queryClient.getQueryData(singleEventQueryKey) as GetTimelineEventsQuery
    )?.timelineEvents?.[0];

    const timelineInfiniteQueryCachedData =
      queryClient.getQueryData(timelineQueryKey);

    return [timelineEventsQueryCachedData, timelineInfiniteQueryCachedData];
  };

  const findTimelineEventByIdInPages = (
    pages: Array<{ organization: { timelineEvents: Array<TimelineEvent> } }>,
    eventId: string,
  ) => {
    if (!pages?.length || !eventId) {
      return null;
    }

    const eventMap = new Map<string, TimelineEvent>();

    pages.forEach((page) => {
      const timelineEvents = page?.organization?.timelineEvents;
      timelineEvents.forEach((event: TimelineEvent) =>
        eventMap.set(event.id, event),
      );
    });

    return eventMap.get(eventId) || null;
  };

  const handleFindTimelineEventInCache = (timelineEventId: string) => {
    const [singleEventQueryKey, timelineQueryKey] = getKeys(timelineEventId);

    const [timelineEventsQueryCachedData, timelineInfiniteQueryCachedData] =
      getQueryData(singleEventQueryKey, timelineQueryKey);

    return (
      timelineEventsQueryCachedData ||
      findTimelineEventByIdInPages(
        (timelineInfiniteQueryCachedData as unknown as any)?.pages,
        timelineEventId,
      )
    );
  };

  return { handleFindTimelineEventInCache };
};
