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
export type MrrPerCustomerQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;


export type MrrPerCustomerQuery = { __typename?: 'Query', dashboard_MRRPerCustomer?: { __typename?: 'DashboardMRRPerCustomer', mrrPerCustomer: number, increasePercentage: string, perMonth: Array<{ __typename?: 'DashboardMRRPerCustomerPerMonth', month: number, value: number } | null> } | null };



export const MrrPerCustomerDocument = `
    query MrrPerCustomer($period: DashboardPeriodInput) {
  dashboard_MRRPerCustomer(period: $period) {
    mrrPerCustomer
    increasePercentage
    perMonth {
      month
      value
    }
  }
}
    `;

export const useMrrPerCustomerQuery = <
      TData = MrrPerCustomerQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: MrrPerCustomerQueryVariables,
      options?: Omit<UseQueryOptions<MrrPerCustomerQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<MrrPerCustomerQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<MrrPerCustomerQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['MrrPerCustomer'] : ['MrrPerCustomer', variables],
    queryFn: fetcher<MrrPerCustomerQuery, MrrPerCustomerQueryVariables>(client, MrrPerCustomerDocument, variables, headers),
    ...options
  }
    )};

useMrrPerCustomerQuery.document = MrrPerCustomerDocument;

useMrrPerCustomerQuery.getKey = (variables?: MrrPerCustomerQueryVariables) => variables === undefined ? ['MrrPerCustomer'] : ['MrrPerCustomer', variables];

export const useInfiniteMrrPerCustomerQuery = <
      TData = InfiniteData<MrrPerCustomerQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: MrrPerCustomerQueryVariables,
      options: Omit<UseInfiniteQueryOptions<MrrPerCustomerQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<MrrPerCustomerQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<MrrPerCustomerQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['MrrPerCustomer.infinite'] : ['MrrPerCustomer.infinite', variables],
      queryFn: (metaData) => fetcher<MrrPerCustomerQuery, MrrPerCustomerQueryVariables>(client, MrrPerCustomerDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteMrrPerCustomerQuery.getKey = (variables?: MrrPerCustomerQueryVariables) => variables === undefined ? ['MrrPerCustomer.infinite'] : ['MrrPerCustomer.infinite', variables];


useMrrPerCustomerQuery.fetcher = (client: GraphQLClient, variables?: MrrPerCustomerQueryVariables, headers?: RequestInit['headers']) => fetcher<MrrPerCustomerQuery, MrrPerCustomerQueryVariables>(client, MrrPerCustomerDocument, variables, headers);


useMrrPerCustomerQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: MrrPerCustomerQueryVariables) =>
  (mutator: (cacheEntry: MrrPerCustomerQuery) => MrrPerCustomerQuery) => {
    const cacheKey = useMrrPerCustomerQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<MrrPerCustomerQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<MrrPerCustomerQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteMrrPerCustomerQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: MrrPerCustomerQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<MrrPerCustomerQuery>) => InfiniteData<MrrPerCustomerQuery>) => {
    const cacheKey = useInfiniteMrrPerCustomerQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<MrrPerCustomerQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<MrrPerCustomerQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }