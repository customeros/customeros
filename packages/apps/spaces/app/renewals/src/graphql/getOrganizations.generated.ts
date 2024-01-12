// @ts-nocheck remove this when typscript-react-query plugin is fixed
import type { InfiniteData } from '@tanstack/react-query';

useGetOrganizationsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetOrganizationsQueryVariables) =>
  (mutator: (cacheEntry: GetOrganizationsQuery) => GetOrganizationsQuery) => {
    const cacheKey = useGetOrganizationsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetOrganizationsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetOrganizationsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetOrganizationsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetOrganizationsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetOrganizationsQuery>,
    ) => InfiniteData<GetOrganizationsQuery>,
  ) => {
    const cacheKey = useInfiniteGetOrganizationsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
