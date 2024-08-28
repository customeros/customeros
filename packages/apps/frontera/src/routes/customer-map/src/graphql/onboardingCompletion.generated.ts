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
export type OnboardingCompletionQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;


export type OnboardingCompletionQuery = { __typename?: 'Query', dashboard_OnboardingCompletion?: { __typename?: 'DashboardOnboardingCompletion', completionPercentage: number, increasePercentage: number, perMonth: Array<{ __typename?: 'DashboardOnboardingCompletionPerMonth', month: number, value: number }> } | null };



export const OnboardingCompletionDocument = `
    query OnboardingCompletion($period: DashboardPeriodInput) {
  dashboard_OnboardingCompletion(period: $period) {
    completionPercentage
    increasePercentage
    perMonth {
      month
      value
    }
  }
}
    `;

export const useOnboardingCompletionQuery = <
      TData = OnboardingCompletionQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: OnboardingCompletionQueryVariables,
      options?: Omit<UseQueryOptions<OnboardingCompletionQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<OnboardingCompletionQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<OnboardingCompletionQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['OnboardingCompletion'] : ['OnboardingCompletion', variables],
    queryFn: fetcher<OnboardingCompletionQuery, OnboardingCompletionQueryVariables>(client, OnboardingCompletionDocument, variables, headers),
    ...options
  }
    )};

useOnboardingCompletionQuery.document = OnboardingCompletionDocument;

useOnboardingCompletionQuery.getKey = (variables?: OnboardingCompletionQueryVariables) => variables === undefined ? ['OnboardingCompletion'] : ['OnboardingCompletion', variables];

export const useInfiniteOnboardingCompletionQuery = <
      TData = InfiniteData<OnboardingCompletionQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: OnboardingCompletionQueryVariables,
      options: Omit<UseInfiniteQueryOptions<OnboardingCompletionQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<OnboardingCompletionQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<OnboardingCompletionQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['OnboardingCompletion.infinite'] : ['OnboardingCompletion.infinite', variables],
      queryFn: (metaData) => fetcher<OnboardingCompletionQuery, OnboardingCompletionQueryVariables>(client, OnboardingCompletionDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteOnboardingCompletionQuery.getKey = (variables?: OnboardingCompletionQueryVariables) => variables === undefined ? ['OnboardingCompletion.infinite'] : ['OnboardingCompletion.infinite', variables];


useOnboardingCompletionQuery.fetcher = (client: GraphQLClient, variables?: OnboardingCompletionQueryVariables, headers?: RequestInit['headers']) => fetcher<OnboardingCompletionQuery, OnboardingCompletionQueryVariables>(client, OnboardingCompletionDocument, variables, headers);


useOnboardingCompletionQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: OnboardingCompletionQueryVariables) =>
  (mutator: (cacheEntry: OnboardingCompletionQuery) => OnboardingCompletionQuery) => {
    const cacheKey = useOnboardingCompletionQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<OnboardingCompletionQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<OnboardingCompletionQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteOnboardingCompletionQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: OnboardingCompletionQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<OnboardingCompletionQuery>) => InfiniteData<OnboardingCompletionQuery>) => {
    const cacheKey = useInfiniteOnboardingCompletionQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<OnboardingCompletionQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<OnboardingCompletionQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }