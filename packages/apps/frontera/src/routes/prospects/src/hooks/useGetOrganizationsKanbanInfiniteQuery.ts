import { GraphQLClient } from 'graphql-request';
import {
  InfiniteData,
  keepPreviousData,
  useInfiniteQuery,
  UseInfiniteQueryOptions,
} from '@tanstack/react-query';

import {
  GetOrganizationsKanbanQuery,
  GetOrganizationsKanbanDocument,
  GetOrganizationsKanbanQueryVariables,
} from '../graphql/getOrganizationsKanban.generated';

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

export const useGetOrganizationsKanbanInfiniteQuery = <
  TData = InfiniteData<GetOrganizationsKanbanQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationsKanbanQueryVariables,
  options?: Omit<
    UseInfiniteQueryOptions<GetOrganizationsKanbanQuery, TError, TData>,
    'queryKey' | 'getNextPageParam' | 'initialPageParam'
  >,
) => {
  return useInfiniteQuery<GetOrganizationsKanbanQuery, TError, TData>({
    queryKey: ['getOrganizationKanban.infinite'],
    queryFn: ({ pageParam = 1 }) =>
      fetcher<
        GetOrganizationsKanbanQuery,
        GetOrganizationsKanbanQueryVariables
      >(client, GetOrganizationsKanbanDocument, {
        ...variables,
        pagination: { ...variables.pagination, page: pageParam as number },
      })(),
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
