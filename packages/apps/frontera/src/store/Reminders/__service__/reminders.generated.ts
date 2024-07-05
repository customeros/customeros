// @ts-nocheck remove this when typscript-react-query plugin is fixed

useRemindersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: RemindersQueryVariables) =>
  (mutator: (cacheEntry: RemindersQuery) => RemindersQuery) => {
    const cacheKey = useRemindersQuery.getKey(variables);
    const previousEntries = queryClient.getQueryData<RemindersQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<RemindersQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteRemindersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: RemindersQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<RemindersQuery>,
    ) => InfiniteData<RemindersQuery>,
  ) => {
    const cacheKey = useInfiniteRemindersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<RemindersQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<RemindersQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  };
