import { GraphQLClient } from 'graphql-request';
import {
  InfiniteData,
  keepPreviousData,
  useInfiniteQuery,
  UseInfiniteQueryOptions,
} from '@tanstack/react-query';

import {
  GetInvoicesQuery,
  GetInvoicesDocument,
  GetInvoicesQueryVariables,
} from '@shared/graphql/getInvoices.generated';

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

export const useGetInvoicesInfiniteQuery = <
  TData = InfiniteData<GetInvoicesQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetInvoicesQueryVariables,
  options?: Omit<
    UseInfiniteQueryOptions<GetInvoicesQuery, TError, TData>,
    'queryKey' | 'getNextPageParam' | 'initialPageParam'
  >,
) => {
  return useInfiniteQuery<GetInvoicesQuery, TError, TData>({
    queryKey: ['getInvoices.infinite', variables],
    queryFn: ({ pageParam = 1 }) =>
      fetcher<GetInvoicesQuery, GetInvoicesQueryVariables>(
        client,
        GetInvoicesDocument,
        {
          ...variables,
          pagination: { ...variables.pagination, page: pageParam as number },
        },
      )(),
    initialPageParam: 0,
    getNextPageParam: (lastPage, pages) => {
      const content = pages.flatMap((page) => page.invoices?.content ?? []);
      const totalElements = lastPage.invoices?.totalElements ?? 0;

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
