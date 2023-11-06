import { QueryKey, useQueryClient } from '@tanstack/react-query';

import { TimelineEvent } from '@graphql/types';
import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import {
  GetTimelineEventsQuery,
  useGetTimelineEventsQuery,
} from '@organization/src/graphql/getTimelineEvents.generated';

// TODO: revisit if this abstraction is useful for us or not or if we need to re-implement in a better way
// Typing is done just so the compiler does not complain.

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
    const timelineEventsQueryCachedData =
      (queryClient.getQueryData(singleEventQueryKey) as GetTimelineEventsQuery)
        ?.timelineEvents?.[0] ?? null;

    const timelineInfiniteQueryCachedData =
      queryClient.getQueryData(timelineQueryKey);

    return [timelineEventsQueryCachedData, timelineInfiniteQueryCachedData];
  };

  const findTimelineEventByIdInPages = (
    pages: Array<{ organization: { timelineEvents: Array<TimelineEvent> } }>,
    eventId: string,
  ): TimelineEvent | null => {
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

    return (eventMap.get(eventId) as TimelineEvent) || null;
  };

  const handleFindTimelineEventInCache = (
    timelineEventId: string,
  ): TimelineEvent | undefined | null => {
    const [singleEventQueryKey, timelineQueryKey] = getKeys(timelineEventId);

    const [timelineEventsQueryCachedData, timelineInfiniteQueryCachedData] =
      getQueryData(singleEventQueryKey, timelineQueryKey);

    return (
      (timelineEventsQueryCachedData as TimelineEvent) ||
      findTimelineEventByIdInPages(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (timelineInfiniteQueryCachedData as unknown as any)?.pages,
        timelineEventId,
      )
    );
  };

  return { handleFindTimelineEventInCache };
};
