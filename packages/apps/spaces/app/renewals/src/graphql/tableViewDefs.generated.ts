// @ts-nocheck remove this when typscript-react-query plugin is fixed
import type { InfiniteData } from '@tanstack/react-query';

useTableViewDefsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: TableViewDefsQueryVariables) =>
  (mutator: (cacheEntry: TableViewDefsQuery) => TableViewDefsQuery) => {
    const cacheKey = useTableViewDefsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TableViewDefsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TableViewDefsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteTableViewDefsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: TableViewDefsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<TableViewDefsQuery>,
    ) => InfiniteData<TableViewDefsQuery>,
  ) => {
    const cacheKey = useInfiniteTableViewDefsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TableViewDefsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TableViewDefsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
