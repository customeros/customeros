// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, useInfiniteQuery, UseQueryOptions, UseInfiniteQueryOptions, InfiniteData } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type NewCustomersQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;


export type NewCustomersQuery = { __typename?: 'Query', dashboard_NewCustomers?: { __typename?: 'DashboardNewCustomers', thisMonthCount: number, thisMonthIncreasePercentage: string, perMonth: Array<{ __typename?: 'DashboardNewCustomersPerMonth', year: number, month: number, count: number } | null> } | null };



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
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: NewCustomersQueryVariables,
      options?: Omit<UseQueryOptions<NewCustomersQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<NewCustomersQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<NewCustomersQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['NewCustomers'] : ['NewCustomers', variables],
    queryFn: fetcher<NewCustomersQuery, NewCustomersQueryVariables>(client, NewCustomersDocument, variables, headers),
    ...options
  }
    )};

useNewCustomersQuery.document = NewCustomersDocument;

useNewCustomersQuery.getKey = (variables?: NewCustomersQueryVariables) => variables === undefined ? ['NewCustomers'] : ['NewCustomers', variables];

export const useInfiniteNewCustomersQuery = <
      TData = InfiniteData<NewCustomersQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: NewCustomersQueryVariables,
      options: Omit<UseInfiniteQueryOptions<NewCustomersQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<NewCustomersQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<NewCustomersQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['NewCustomers.infinite'] : ['NewCustomers.infinite', variables],
      queryFn: (metaData) => fetcher<NewCustomersQuery, NewCustomersQueryVariables>(client, NewCustomersDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteNewCustomersQuery.getKey = (variables?: NewCustomersQueryVariables) => variables === undefined ? ['NewCustomers.infinite'] : ['NewCustomers.infinite', variables];


useNewCustomersQuery.fetcher = (client: GraphQLClient, variables?: NewCustomersQueryVariables, headers?: RequestInit['headers']) => fetcher<NewCustomersQuery, NewCustomersQueryVariables>(client, NewCustomersDocument, variables, headers);


useNewCustomersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: NewCustomersQueryVariables) =>
  (mutator: (cacheEntry: NewCustomersQuery) => NewCustomersQuery) => {
    const cacheKey = useNewCustomersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<NewCustomersQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<NewCustomersQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteNewCustomersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: NewCustomersQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<NewCustomersQuery>) => InfiniteData<NewCustomersQuery>) => {
    const cacheKey = useInfiniteNewCustomersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<NewCustomersQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<NewCustomersQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }