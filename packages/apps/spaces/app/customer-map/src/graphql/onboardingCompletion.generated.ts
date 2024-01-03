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
export type OnboardingCompletionQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;

export type OnboardingCompletionQuery = {
  __typename?: 'Query';
  dashboard_OnboardingCompletion?: {
    __typename?: 'DashboardOnboardingCompletion';
    completionPercentage: number;
    increasePercentage: number;
    perMonth: Array<{
      __typename?: 'DashboardOnboardingCompletionPerMonth';
      month: number;
      value: number;
    }>;
  } | null;
};

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
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: OnboardingCompletionQueryVariables,
  options?: UseQueryOptions<OnboardingCompletionQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<OnboardingCompletionQuery, TError, TData>(
    variables === undefined
      ? ['OnboardingCompletion']
      : ['OnboardingCompletion', variables],
    fetcher<OnboardingCompletionQuery, OnboardingCompletionQueryVariables>(
      client,
      OnboardingCompletionDocument,
      variables,
      headers,
    ),
    options,
  );
useOnboardingCompletionQuery.document = OnboardingCompletionDocument;

useOnboardingCompletionQuery.getKey = (
  variables?: OnboardingCompletionQueryVariables,
) =>
  variables === undefined
    ? ['OnboardingCompletion']
    : ['OnboardingCompletion', variables];
export const useInfiniteOnboardingCompletionQuery = <
  TData = OnboardingCompletionQuery,
  TError = unknown,
>(
  pageParamKey: keyof OnboardingCompletionQueryVariables,
  client: GraphQLClient,
  variables?: OnboardingCompletionQueryVariables,
  options?: UseInfiniteQueryOptions<OnboardingCompletionQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<OnboardingCompletionQuery, TError, TData>(
    variables === undefined
      ? ['OnboardingCompletion.infinite']
      : ['OnboardingCompletion.infinite', variables],
    (metaData) =>
      fetcher<OnboardingCompletionQuery, OnboardingCompletionQueryVariables>(
        client,
        OnboardingCompletionDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteOnboardingCompletionQuery.getKey = (
  variables?: OnboardingCompletionQueryVariables,
) =>
  variables === undefined
    ? ['OnboardingCompletion.infinite']
    : ['OnboardingCompletion.infinite', variables];
useOnboardingCompletionQuery.fetcher = (
  client: GraphQLClient,
  variables?: OnboardingCompletionQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<OnboardingCompletionQuery, OnboardingCompletionQueryVariables>(
    client,
    OnboardingCompletionDocument,
    variables,
    headers,
  );

useOnboardingCompletionQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: OnboardingCompletionQueryVariables) =>
  (
    mutator: (
      cacheEntry: OnboardingCompletionQuery,
    ) => OnboardingCompletionQuery,
  ) => {
    const cacheKey = useOnboardingCompletionQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<OnboardingCompletionQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<OnboardingCompletionQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteOnboardingCompletionQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: OnboardingCompletionQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<OnboardingCompletionQuery>,
    ) => InfiniteData<OnboardingCompletionQuery>,
  ) => {
    const cacheKey = useInfiniteOnboardingCompletionQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<OnboardingCompletionQuery>>(
        cacheKey,
      );
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<OnboardingCompletionQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
