import { QueryKey, InfiniteData, useQueryClient } from '@tanstack/react-query';

import { TimelineEvent } from '@organization/src/components/Timeline/types';
import { GetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';

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
                event.__typename !== 'Analysis' &&
                event.__typename !== 'InteractionSession' &&
                event.__typename !== 'Note' &&
                event.__typename !== 'PageView' &&
                event.id === updatedEvent?.id
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
