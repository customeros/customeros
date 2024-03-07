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
export type GrossRevenueRetentionQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;


export type GrossRevenueRetentionQuery = { __typename?: 'Query', dashboard_GrossRevenueRetention?: { __typename?: 'DashboardGrossRevenueRetention', grossRevenueRetention: number, increasePercentage: string, perMonth: Array<{ __typename?: 'DashboardGrossRevenueRetentionPerMonth', month: number, percentage: number } | null> } | null };



export const GrossRevenueRetentionDocument = `
    query GrossRevenueRetention($period: DashboardPeriodInput) {
  dashboard_GrossRevenueRetention(period: $period) {
    grossRevenueRetention
    increasePercentage
    perMonth {
      month
      percentage
    }
  }
}
    `;

export const useGrossRevenueRetentionQuery = <
      TData = GrossRevenueRetentionQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GrossRevenueRetentionQueryVariables,
      options?: Omit<UseQueryOptions<GrossRevenueRetentionQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GrossRevenueRetentionQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GrossRevenueRetentionQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['GrossRevenueRetention'] : ['GrossRevenueRetention', variables],
    queryFn: fetcher<GrossRevenueRetentionQuery, GrossRevenueRetentionQueryVariables>(client, GrossRevenueRetentionDocument, variables, headers),
    ...options
  }
    )};

useGrossRevenueRetentionQuery.document = GrossRevenueRetentionDocument;

useGrossRevenueRetentionQuery.getKey = (variables?: GrossRevenueRetentionQueryVariables) => variables === undefined ? ['GrossRevenueRetention'] : ['GrossRevenueRetention', variables];

export const useInfiniteGrossRevenueRetentionQuery = <
      TData = InfiniteData<GrossRevenueRetentionQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GrossRevenueRetentionQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GrossRevenueRetentionQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GrossRevenueRetentionQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GrossRevenueRetentionQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['GrossRevenueRetention.infinite'] : ['GrossRevenueRetention.infinite', variables],
      queryFn: (metaData) => fetcher<GrossRevenueRetentionQuery, GrossRevenueRetentionQueryVariables>(client, GrossRevenueRetentionDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGrossRevenueRetentionQuery.getKey = (variables?: GrossRevenueRetentionQueryVariables) => variables === undefined ? ['GrossRevenueRetention.infinite'] : ['GrossRevenueRetention.infinite', variables];


useGrossRevenueRetentionQuery.fetcher = (client: GraphQLClient, variables?: GrossRevenueRetentionQueryVariables, headers?: RequestInit['headers']) => fetcher<GrossRevenueRetentionQuery, GrossRevenueRetentionQueryVariables>(client, GrossRevenueRetentionDocument, variables, headers);


useGrossRevenueRetentionQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GrossRevenueRetentionQueryVariables) =>
  (mutator: (cacheEntry: GrossRevenueRetentionQuery) => GrossRevenueRetentionQuery) => {
    const cacheKey = useGrossRevenueRetentionQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GrossRevenueRetentionQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GrossRevenueRetentionQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGrossRevenueRetentionQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GrossRevenueRetentionQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GrossRevenueRetentionQuery>) => InfiniteData<GrossRevenueRetentionQuery>) => {
    const cacheKey = useInfiniteGrossRevenueRetentionQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GrossRevenueRetentionQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GrossRevenueRetentionQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }