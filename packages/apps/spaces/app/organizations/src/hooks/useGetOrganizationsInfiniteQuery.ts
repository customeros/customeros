import { GraphQLClient } from 'graphql-request';
import {
  InfiniteData,
  keepPreviousData,
  useInfiniteQuery,
  UseInfiniteQueryOptions,
} from '@tanstack/react-query';

import {
  GetOrganizationsQuery,
  GetOrganizationsDocument,
  GetOrganizationsQueryVariables,
} from '../graphql/getOrganizations.generated';

function fetcher<TData, TVariables extends { [key: string]: unknown }>(
  client: GraphQLClient,
  query: string,
  variables?: TVariables,
  requestHeaders?: RequestInit['headers'],
) {
  return async (): Promise<TData> =>
    client.request({
      document: query,
      variables,
      requestHeaders,
    });
}

export const useGetOrganizationsInfiniteQuery = <
  TData = InfiniteData<GetOrganizationsQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationsQueryVariables,
  options?: Omit<
    UseInfiniteQueryOptions<GetOrganizationsQuery, TError, TData>,
    'queryKey' | 'getNextPageParam' | 'initialPageParam'
  >,
) => {
  return useInfiniteQuery<GetOrganizationsQuery, TError, TData>({
    queryKey: ['getOrganizations.infinite', variables],
    queryFn: ({ pageParam = 1 }) =>
      fetcher<GetOrganizationsQuery, GetOrganizationsQueryVariables>(
        client,
        GetOrganizationsDocument,
        {
          ...variables,
          pagination: { ...variables.pagination, page: pageParam as number },
        },
      )(),
    initialPageParam: 1,
    getNextPageParam: (lastPage, pages) => {
      const content = pages.flatMap(
        (page) => page.dashboardView_Organizations?.content ?? [],
      );
      const totalElements =
        lastPage.dashboardView_Organizations?.totalElements ?? 0;

      if (content.length >= totalElements) {
        return undefined;
      }

      return pages.length + 1;
    },
    refetchOnWindowFocus: false,
    placeholderData: keepPreviousData,
    ...options,
  });
};
