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
export type RevenueAtRiskQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;

export type RevenueAtRiskQuery = {
  __typename?: 'Query';
  dashboard_RevenueAtRisk?: {
    __typename?: 'DashboardRevenueAtRisk';
    highConfidence: number;
    atRisk: number;
  } | null;
};

export const RevenueAtRiskDocument = `
    query RevenueAtRisk($period: DashboardPeriodInput) {
  dashboard_RevenueAtRisk(period: $period) {
    highConfidence
    atRisk
  }
}
    `;
export const useRevenueAtRiskQuery = <
  TData = RevenueAtRiskQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: RevenueAtRiskQueryVariables,
  options?: UseQueryOptions<RevenueAtRiskQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<RevenueAtRiskQuery, TError, TData>(
    variables === undefined ? ['RevenueAtRisk'] : ['RevenueAtRisk', variables],
    fetcher<RevenueAtRiskQuery, RevenueAtRiskQueryVariables>(
      client,
      RevenueAtRiskDocument,
      variables,
      headers,
    ),
    options,
  );
useRevenueAtRiskQuery.document = RevenueAtRiskDocument;

useRevenueAtRiskQuery.getKey = (variables?: RevenueAtRiskQueryVariables) =>
  variables === undefined ? ['RevenueAtRisk'] : ['RevenueAtRisk', variables];
export const useInfiniteRevenueAtRiskQuery = <
  TData = RevenueAtRiskQuery,
  TError = unknown,
>(
  pageParamKey: keyof RevenueAtRiskQueryVariables,
  client: GraphQLClient,
  variables?: RevenueAtRiskQueryVariables,
  options?: UseInfiniteQueryOptions<RevenueAtRiskQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<RevenueAtRiskQuery, TError, TData>(
    variables === undefined
      ? ['RevenueAtRisk.infinite']
      : ['RevenueAtRisk.infinite', variables],
    (metaData) =>
      fetcher<RevenueAtRiskQuery, RevenueAtRiskQueryVariables>(
        client,
        RevenueAtRiskDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteRevenueAtRiskQuery.getKey = (
  variables?: RevenueAtRiskQueryVariables,
) =>
  variables === undefined
    ? ['RevenueAtRisk.infinite']
    : ['RevenueAtRisk.infinite', variables];
useRevenueAtRiskQuery.fetcher = (
  client: GraphQLClient,
  variables?: RevenueAtRiskQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<RevenueAtRiskQuery, RevenueAtRiskQueryVariables>(
    client,
    RevenueAtRiskDocument,
    variables,
    headers,
  );

useRevenueAtRiskQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: RevenueAtRiskQueryVariables) =>
  (mutator: (cacheEntry: RevenueAtRiskQuery) => RevenueAtRiskQuery) => {
    const cacheKey = useRevenueAtRiskQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<RevenueAtRiskQuery>(cacheKey);
    if (previousEntry) {
      queryClient.setQueryData<RevenueAtRiskQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteRevenueAtRiskQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: RevenueAtRiskQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<RevenueAtRiskQuery>,
    ) => InfiniteData<RevenueAtRiskQuery>,
  ) => {
    const cacheKey = useInfiniteRevenueAtRiskQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<RevenueAtRiskQuery>>(cacheKey);
    if (previousEntry) {
      queryClient.setQueryData<InfiniteData<RevenueAtRiskQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
