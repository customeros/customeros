import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import {
  GetTimelineQuery,
  useInfiniteGetTimelineQuery,
} from '@organization/src/graphql/getTimeline.generated';
import {
  GetTimelineEventsQuery,
  useGetTimelineEventsQuery,
} from '@organization/src/graphql/getTimelineEvents.generated';
import { InfiniteData, QueryKey, useQueryClient } from '@tanstack/react-query';

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

  const findTimelineEventByIdInPages = (pages: Array<any>, eventId: string) => {
    if (!pages?.length || !eventId) {
      return null;
    }
    // Providing correct type to pages and renaming function
    for (let i = 0; i < pages.length; i++) {
      const timelineEvents = pages?.[i]?.organization.timelineEvents;
      for (let j = 0; j < timelineEvents.length; j++) {
        if (timelineEvents[j]?.id === eventId) {
          return timelineEvents[j];
        }
      }
    }
    return null;
  };

  const handleFindTimelineEventInCache = (timelineEventId: string) => {
    const [singleEventQueryKey, timelineQueryKey] = getKeys(timelineEventId);

    const [timelineEventsQueryCachedData, timelineInfiniteQueryCachedData] =
      getQueryData(singleEventQueryKey, timelineQueryKey);

    return (
      timelineEventsQueryCachedData ||
      findTimelineEventByIdInPages(
        (timelineInfiniteQueryCachedData as InfiniteData<GetTimelineQuery>)
          ?.pages,
        timelineEventId,
      )
    );
  };

  return { handleFindTimelineEventInCache };
};
