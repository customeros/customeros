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
export type ArrBreakdownQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;


export type ArrBreakdownQuery = { __typename?: 'Query', dashboard_ARRBreakdown?: { __typename?: 'DashboardARRBreakdown', arrBreakdown: number, increasePercentage: string, perMonth: Array<{ __typename?: 'DashboardARRBreakdownPerMonth', year: number, month: number, newlyContracted: number, renewals: number, upsells: number, downgrades: number, cancellations: number, churned: number } | null> } | null };



export const ArrBreakdownDocument = `
    query ArrBreakdown($period: DashboardPeriodInput) {
  dashboard_ARRBreakdown(period: $period) {
    arrBreakdown
    increasePercentage
    perMonth {
      year
      month
      newlyContracted
      renewals
      upsells
      downgrades
      cancellations
      churned
    }
  }
}
    `;

export const useArrBreakdownQuery = <
      TData = ArrBreakdownQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: ArrBreakdownQueryVariables,
      options?: Omit<UseQueryOptions<ArrBreakdownQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<ArrBreakdownQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<ArrBreakdownQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['ArrBreakdown'] : ['ArrBreakdown', variables],
    queryFn: fetcher<ArrBreakdownQuery, ArrBreakdownQueryVariables>(client, ArrBreakdownDocument, variables, headers),
    ...options
  }
    )};

useArrBreakdownQuery.document = ArrBreakdownDocument;

useArrBreakdownQuery.getKey = (variables?: ArrBreakdownQueryVariables) => variables === undefined ? ['ArrBreakdown'] : ['ArrBreakdown', variables];

export const useInfiniteArrBreakdownQuery = <
      TData = InfiniteData<ArrBreakdownQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: ArrBreakdownQueryVariables,
      options: Omit<UseInfiniteQueryOptions<ArrBreakdownQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<ArrBreakdownQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<ArrBreakdownQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['ArrBreakdown.infinite'] : ['ArrBreakdown.infinite', variables],
      queryFn: (metaData) => fetcher<ArrBreakdownQuery, ArrBreakdownQueryVariables>(client, ArrBreakdownDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteArrBreakdownQuery.getKey = (variables?: ArrBreakdownQueryVariables) => variables === undefined ? ['ArrBreakdown.infinite'] : ['ArrBreakdown.infinite', variables];


useArrBreakdownQuery.fetcher = (client: GraphQLClient, variables?: ArrBreakdownQueryVariables, headers?: RequestInit['headers']) => fetcher<ArrBreakdownQuery, ArrBreakdownQueryVariables>(client, ArrBreakdownDocument, variables, headers);


useArrBreakdownQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: ArrBreakdownQueryVariables) =>
  (mutator: (cacheEntry: ArrBreakdownQuery) => ArrBreakdownQuery) => {
    const cacheKey = useArrBreakdownQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<ArrBreakdownQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<ArrBreakdownQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteArrBreakdownQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: ArrBreakdownQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<ArrBreakdownQuery>) => InfiniteData<ArrBreakdownQuery>) => {
    const cacheKey = useInfiniteArrBreakdownQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<ArrBreakdownQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<ArrBreakdownQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }