import { QueryKey, InfiniteData, useQueryClient } from '@tanstack/react-query';

import { GetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
type EventWithId = { id: string; [key: string]: any };

export function useUpdateCacheWithExistingEvent() {
  const queryClient = useQueryClient();

  return async (updatedEvent: EventWithId, queryKey: QueryKey) => {
    await queryClient.cancelQueries({ queryKey });
    const previousTimelineEntries =
      queryClient.getQueryData<InfiniteData<GetTimelineQuery>>(queryKey);

    queryClient.setQueryData<InfiniteData<GetTimelineQuery>>(
      queryKey,
      (currentCache): InfiniteData<GetTimelineQuery> => {
        const updatedPages = currentCache?.pages?.map((page) => {
          const updatedEvents = page?.organization?.timelineEvents?.map(
            (event: Record<string, any>) => {
              if (event.id === updatedEvent?.id) {
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
