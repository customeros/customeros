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
export type RetentionRateQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;


export type RetentionRateQuery = { __typename?: 'Query', dashboard_RetentionRate?: { __typename?: 'DashboardRetentionRate', retentionRate: number, increasePercentage: string, perMonth: Array<{ __typename?: 'DashboardRetentionRatePerMonth', month: number, renewCount: number, churnCount: number } | null> } | null };



export const RetentionRateDocument = `
    query RetentionRate($period: DashboardPeriodInput) {
  dashboard_RetentionRate(period: $period) {
    retentionRate
    increasePercentage
    perMonth {
      month
      renewCount
      churnCount
    }
  }
}
    `;

export const useRetentionRateQuery = <
      TData = RetentionRateQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: RetentionRateQueryVariables,
      options?: Omit<UseQueryOptions<RetentionRateQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<RetentionRateQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<RetentionRateQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['RetentionRate'] : ['RetentionRate', variables],
    queryFn: fetcher<RetentionRateQuery, RetentionRateQueryVariables>(client, RetentionRateDocument, variables, headers),
    ...options
  }
    )};

useRetentionRateQuery.document = RetentionRateDocument;

useRetentionRateQuery.getKey = (variables?: RetentionRateQueryVariables) => variables === undefined ? ['RetentionRate'] : ['RetentionRate', variables];

export const useInfiniteRetentionRateQuery = <
      TData = InfiniteData<RetentionRateQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: RetentionRateQueryVariables,
      options: Omit<UseInfiniteQueryOptions<RetentionRateQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<RetentionRateQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<RetentionRateQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['RetentionRate.infinite'] : ['RetentionRate.infinite', variables],
      queryFn: (metaData) => fetcher<RetentionRateQuery, RetentionRateQueryVariables>(client, RetentionRateDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteRetentionRateQuery.getKey = (variables?: RetentionRateQueryVariables) => variables === undefined ? ['RetentionRate.infinite'] : ['RetentionRate.infinite', variables];


useRetentionRateQuery.fetcher = (client: GraphQLClient, variables?: RetentionRateQueryVariables, headers?: RequestInit['headers']) => fetcher<RetentionRateQuery, RetentionRateQueryVariables>(client, RetentionRateDocument, variables, headers);


useRetentionRateQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: RetentionRateQueryVariables) =>
  (mutator: (cacheEntry: RetentionRateQuery) => RetentionRateQuery) => {
    const cacheKey = useRetentionRateQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<RetentionRateQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<RetentionRateQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteRetentionRateQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: RetentionRateQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<RetentionRateQuery>) => InfiniteData<RetentionRateQuery>) => {
    const cacheKey = useInfiniteRetentionRateQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<RetentionRateQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<RetentionRateQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }