// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

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
export type RetentionRateQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;

export type RetentionRateQuery = {
  __typename?: 'Query';
  dashboard_RetentionRate?: {
    __typename?: 'DashboardRetentionRate';
    retentionRate: number;
    increasePercentage: string;
    perMonth: Array<{
      __typename?: 'DashboardRetentionRatePerMonth';
      month: number;
      renewCount: number;
      churnCount: number;
    } | null>;
  } | null;
};

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
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: RetentionRateQueryVariables,
  options?: UseQueryOptions<RetentionRateQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<RetentionRateQuery, TError, TData>(
    variables === undefined ? ['RetentionRate'] : ['RetentionRate', variables],
    fetcher<RetentionRateQuery, RetentionRateQueryVariables>(
      client,
      RetentionRateDocument,
      variables,
      headers,
    ),
    options,
  );
useRetentionRateQuery.document = RetentionRateDocument;

useRetentionRateQuery.getKey = (variables?: RetentionRateQueryVariables) =>
  variables === undefined ? ['RetentionRate'] : ['RetentionRate', variables];
export const useInfiniteRetentionRateQuery = <
  TData = RetentionRateQuery,
  TError = unknown,
>(
  pageParamKey: keyof RetentionRateQueryVariables,
  client: GraphQLClient,
  variables?: RetentionRateQueryVariables,
  options?: UseInfiniteQueryOptions<RetentionRateQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<RetentionRateQuery, TError, TData>(
    variables === undefined
      ? ['RetentionRate.infinite']
      : ['RetentionRate.infinite', variables],
    (metaData) =>
      fetcher<RetentionRateQuery, RetentionRateQueryVariables>(
        client,
        RetentionRateDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteRetentionRateQuery.getKey = (
  variables?: RetentionRateQueryVariables,
) =>
  variables === undefined
    ? ['RetentionRate.infinite']
    : ['RetentionRate.infinite', variables];
useRetentionRateQuery.fetcher = (
  client: GraphQLClient,
  variables?: RetentionRateQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<RetentionRateQuery, RetentionRateQueryVariables>(
    client,
    RetentionRateDocument,
    variables,
    headers,
  );
