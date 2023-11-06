import { QueryKey, InfiniteData, useQueryClient } from '@tanstack/react-query';

import { GetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';

// TODO: replace this type with something explicit or remove it.
// eslint-disable-next-line @typescript-eslint/no-explicit-any
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
            (event) => {
              // this code should exhaustively match agains all possible types of Events
              // __typename: "Analysis" does not have "id" which will cause this check to fail
              // @ts-expect-error TODO: match(event).with(***) for all cases.
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
