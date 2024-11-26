import { QueryKey, InfiniteData, useQueryClient } from '@tanstack/react-query';

import { LogEntry } from '@graphql/types';
import { TimelineEvent } from '@organization/components/Timeline/types';
import { GetTimelineQuery } from '@organization/graphql/getTimeline.generated';

export function useUpdateCacheWithExistingEvent() {
  const queryClient = useQueryClient();

  return async (updatedEvent: TimelineEvent, queryKey: QueryKey) => {
    await queryClient.cancelQueries({ queryKey });

    const previousTimelineEntries =
      queryClient.getQueryData<InfiniteData<GetTimelineQuery>>(queryKey);

    queryClient.setQueryData<InfiniteData<GetTimelineQuery>>(
      queryKey,
      (currentCache): InfiniteData<GetTimelineQuery> => {
        const updatedPages = currentCache?.pages?.map((page) => {
          const updatedEvents = page?.organization?.timelineEvents?.map(
            (event) => {
              if (
                // those entities need fragments on the query string in order to have `id` field present
                // we filter them out on this condition to avoid TS saying there's no `id` present when checking.
                event.__typename === 'LogEntry' &&
                event.id === (updatedEvent as LogEntry)?.id
              ) {
                return { ...event, ...updatedEvent };
              }

              return event;
            },
          );

          return {
            ...page,
            organization: {
              ...page?.organization,
              timelineEvents: updatedEvents,
            },
          };
        });

        return {
          ...currentCache,
          pages: updatedPages,
        } as InfiniteData<GetTimelineQuery>;
      },
    );

    return { previousTimelineEntries };
  };
}
