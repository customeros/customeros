import { GraphQLClient } from 'graphql-request';
import {
  useInfiniteQuery,
  UseInfiniteQueryOptions,
} from '@tanstack/react-query';

import {
  GetRenewalsQuery,
  GetRenewalsDocument,
  GetRenewalsQueryVariables,
} from '../graphql/getRenewals.generated';

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

export const useGetRenewalsInfiniteQuery = <
  TData = GetRenewalsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetRenewalsQueryVariables,
  _options?: UseInfiniteQueryOptions<GetRenewalsQuery, TError, TData>,
) => {
  return useInfiniteQuery<GetRenewalsQuery, TError, TData>({
    queryKey: ['getRenewals.infinite', variables],
    queryFn: ({ pageParam = 1 }) =>
      fetcher<GetRenewalsQuery, GetRenewalsQueryVariables>(
        client,
        GetRenewalsDocument,
        {
          ...variables,
          pagination: { ...variables.pagination, page: pageParam },
        },
      )(),
    getNextPageParam: (lastPage, pages) => {
      const content = pages.flatMap(
        (page) => page.dashboardView_Renewals?.content ?? [],
      );
      const totalElements = lastPage.dashboardView_Renewals?.totalElements ?? 0;

      if (content.length >= totalElements) {
        return undefined;
      }

      return pages.length + 1;
    },
    refetchOnWindowFocus: false,
    keepPreviousData: true,
  });
};
