// @ts-nocheck remove this when typscript-react-query plugin is fixed
import type { InfiniteData } from '@tanstack/react-query';

useGetUsersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetUsersQueryVariables) =>
  (mutator: (cacheEntry: GetUsersQuery) => GetUsersQuery) => {
    const cacheKey = useGetUsersQuery.getKey(variables);
    const previousEntries = queryClient.getQueryData<GetUsersQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetUsersQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetUsersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetUsersQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetUsersQuery>,
    ) => InfiniteData<GetUsersQuery>,
  ) => {
    const cacheKey = useInfiniteGetUsersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetUsersQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetUsersQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  };
