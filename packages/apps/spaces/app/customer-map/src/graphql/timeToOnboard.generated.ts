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
export type TimeToOnboardQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;


export type TimeToOnboardQuery = { __typename?: 'Query', dashboard_TimeToOnboard?: { __typename?: 'DashboardTimeToOnboard', timeToOnboard?: number | null, increasePercentage?: number | null, perMonth: Array<{ __typename?: 'DashboardTimeToOnboardPerMonth', month: number, value: number }> } | null };



export const TimeToOnboardDocument = `
    query TimeToOnboard($period: DashboardPeriodInput) {
  dashboard_TimeToOnboard(period: $period) {
    timeToOnboard
    increasePercentage
    perMonth {
      month
      value
    }
  }
}
    `;

export const useTimeToOnboardQuery = <
      TData = TimeToOnboardQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: TimeToOnboardQueryVariables,
      options?: Omit<UseQueryOptions<TimeToOnboardQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<TimeToOnboardQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<TimeToOnboardQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['TimeToOnboard'] : ['TimeToOnboard', variables],
    queryFn: fetcher<TimeToOnboardQuery, TimeToOnboardQueryVariables>(client, TimeToOnboardDocument, variables, headers),
    ...options
  }
    )};

useTimeToOnboardQuery.document = TimeToOnboardDocument;

useTimeToOnboardQuery.getKey = (variables?: TimeToOnboardQueryVariables) => variables === undefined ? ['TimeToOnboard'] : ['TimeToOnboard', variables];

export const useInfiniteTimeToOnboardQuery = <
      TData = InfiniteData<TimeToOnboardQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: TimeToOnboardQueryVariables,
      options: Omit<UseInfiniteQueryOptions<TimeToOnboardQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<TimeToOnboardQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<TimeToOnboardQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['TimeToOnboard.infinite'] : ['TimeToOnboard.infinite', variables],
      queryFn: (metaData) => fetcher<TimeToOnboardQuery, TimeToOnboardQueryVariables>(client, TimeToOnboardDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteTimeToOnboardQuery.getKey = (variables?: TimeToOnboardQueryVariables) => variables === undefined ? ['TimeToOnboard.infinite'] : ['TimeToOnboard.infinite', variables];


useTimeToOnboardQuery.fetcher = (client: GraphQLClient, variables?: TimeToOnboardQueryVariables, headers?: RequestInit['headers']) => fetcher<TimeToOnboardQuery, TimeToOnboardQueryVariables>(client, TimeToOnboardDocument, variables, headers);


useTimeToOnboardQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TimeToOnboardQueryVariables) =>
  (mutator: (cacheEntry: TimeToOnboardQuery) => TimeToOnboardQuery) => {
    const cacheKey = useTimeToOnboardQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TimeToOnboardQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TimeToOnboardQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteTimeToOnboardQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TimeToOnboardQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<TimeToOnboardQuery>) => InfiniteData<TimeToOnboardQuery>) => {
    const cacheKey = useInfiniteTimeToOnboardQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TimeToOnboardQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TimeToOnboardQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }