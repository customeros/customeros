import { InfiniteData, QueryKey, useQueryClient } from '@tanstack/react-query';
import { GetTimelineQuery } from '@organization/graphql/getTimeline.generated';

export function useUpdateCacheWithNewEvent<T>() {
  const queryClient = useQueryClient();

  return async (newTimelineEvent: T, queryKey: QueryKey) => {
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
  };
}
