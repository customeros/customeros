import { VirtuosoHandle } from 'react-virtuoso';

import { QueryKey, InfiniteData, useQueryClient } from '@tanstack/react-query';

import { GetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';

export function useUpdateCacheWithNewEvent<T>(
  virtuosoRef?: React.RefObject<VirtuosoHandle> | null,
) {
  const queryClient = useQueryClient();

  return async (newTimelineEvent: T, queryKey: QueryKey) => {
    const previousTimelineEntries =
      queryClient.getQueryData<InfiniteData<GetTimelineQuery>>(queryKey);

    const timelineEntries =
      previousTimelineEntries?.pages?.[0]?.organization?.timelineEvents;

    await queryClient.cancelQueries({ queryKey });
    queryClient.setQueryData<InfiniteData<GetTimelineQuery>>(
      queryKey,
      (currentCache): InfiniteData<GetTimelineQuery> => {
        return {
          ...currentCache,
          pages: currentCache?.pages?.map((p, idx) => {
            if (idx !== 0) return p;

            return {
              ...p,
              organization: {
                ...p?.organization,
                timelineEvents: [
                  newTimelineEvent,
                  ...(p?.organization?.timelineEvents ?? []),
                ],
                timelineEventsTotalCount:
                  p?.organization?.timelineEventsTotalCount + 1,
              },
            };
          }),
        } as InfiniteData<GetTimelineQuery>;
      },
    );
    virtuosoRef?.current?.scrollToIndex({
      index: (timelineEntries?.length ?? 0) + 1,
    });

    return { previousTimelineEntries };
  };
}
