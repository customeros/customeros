// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import type { InfiniteData } from '@tanstack/react-query';
import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
} from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(
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
export type NewCustomersQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;

export type NewCustomersQuery = {
  __typename?: 'Query';
  dashboard_NewCustomers?: {
    __typename?: 'DashboardNewCustomers';
    thisMonthCount: number;
    thisMonthIncreasePercentage: string;
    perMonth: Array<{
      __typename?: 'DashboardNewCustomersPerMonth';
      year: number;
      month: number;
      count: number;
    } | null>;
  } | null;
};

export const NewCustomersDocument = `
    query NewCustomers($period: DashboardPeriodInput) {
  dashboard_NewCustomers(period: $period) {
    thisMonthCount
    thisMonthIncreasePercentage
    perMonth {
      year
      month
      count
    }
  }
}
    `;
export const useNewCustomersQuery = <
  TData = NewCustomersQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: NewCustomersQueryVariables,
  options?: UseQueryOptions<NewCustomersQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<NewCustomersQuery, TError, TData>(
    variables === undefined ? ['NewCustomers'] : ['NewCustomers', variables],
    fetcher<NewCustomersQuery, NewCustomersQueryVariables>(
      client,
      NewCustomersDocument,
      variables,
      headers,
    ),
    options,
  );
useNewCustomersQuery.document = NewCustomersDocument;

useNewCustomersQuery.getKey = (variables?: NewCustomersQueryVariables) =>
  variables === undefined ? ['NewCustomers'] : ['NewCustomers', variables];
export const useInfiniteNewCustomersQuery = <
  TData = NewCustomersQuery,
  TError = unknown,
>(
  pageParamKey: keyof NewCustomersQueryVariables,
  client: GraphQLClient,
  variables?: NewCustomersQueryVariables,
  options?: UseInfiniteQueryOptions<NewCustomersQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<NewCustomersQuery, TError, TData>(
    variables === undefined
      ? ['NewCustomers.infinite']
      : ['NewCustomers.infinite', variables],
    (metaData) =>
      fetcher<NewCustomersQuery, NewCustomersQueryVariables>(
        client,
        NewCustomersDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteNewCustomersQuery.getKey = (
  variables?: NewCustomersQueryVariables,
) =>
  variables === undefined
    ? ['NewCustomers.infinite']
    : ['NewCustomers.infinite', variables];
useNewCustomersQuery.fetcher = (
  client: GraphQLClient,
  variables?: NewCustomersQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<NewCustomersQuery, NewCustomersQueryVariables>(
    client,
    NewCustomersDocument,
    variables,
    headers,
  );

useNewCustomersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: NewCustomersQueryVariables) =>
  (mutator: (cacheEntry: NewCustomersQuery) => NewCustomersQuery) => {
    const cacheKey = useNewCustomersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<NewCustomersQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<NewCustomersQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteNewCustomersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: NewCustomersQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<NewCustomersQuery>,
    ) => InfiniteData<NewCustomersQuery>,
  ) => {
    const cacheKey = useInfiniteNewCustomersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<NewCustomersQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<NewCustomersQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
