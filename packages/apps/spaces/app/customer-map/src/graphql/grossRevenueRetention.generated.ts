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
export type GrossRevenueRetentionQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;

export type GrossRevenueRetentionQuery = {
  __typename?: 'Query';
  dashboard_GrossRevenueRetention?: {
    __typename?: 'DashboardGrossRevenueRetention';
    grossRevenueRetention: number;
    increasePercentage: number;
    perMonth: Array<{
      __typename?: 'DashboardGrossRevenueRetentionPerMonth';
      month: number;
      percentage: number;
    } | null>;
  } | null;
};

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
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: GrossRevenueRetentionQueryVariables,
  options?: UseQueryOptions<GrossRevenueRetentionQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GrossRevenueRetentionQuery, TError, TData>(
    variables === undefined
      ? ['GrossRevenueRetention']
      : ['GrossRevenueRetention', variables],
    fetcher<GrossRevenueRetentionQuery, GrossRevenueRetentionQueryVariables>(
      client,
      GrossRevenueRetentionDocument,
      variables,
      headers,
    ),
    options,
  );
useGrossRevenueRetentionQuery.document = GrossRevenueRetentionDocument;

useGrossRevenueRetentionQuery.getKey = (
  variables?: GrossRevenueRetentionQueryVariables,
) =>
  variables === undefined
    ? ['GrossRevenueRetention']
    : ['GrossRevenueRetention', variables];
export const useInfiniteGrossRevenueRetentionQuery = <
  TData = GrossRevenueRetentionQuery,
  TError = unknown,
>(
  pageParamKey: keyof GrossRevenueRetentionQueryVariables,
  client: GraphQLClient,
  variables?: GrossRevenueRetentionQueryVariables,
  options?: UseInfiniteQueryOptions<GrossRevenueRetentionQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GrossRevenueRetentionQuery, TError, TData>(
    variables === undefined
      ? ['GrossRevenueRetention.infinite']
      : ['GrossRevenueRetention.infinite', variables],
    (metaData) =>
      fetcher<GrossRevenueRetentionQuery, GrossRevenueRetentionQueryVariables>(
        client,
        GrossRevenueRetentionDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGrossRevenueRetentionQuery.getKey = (
  variables?: GrossRevenueRetentionQueryVariables,
) =>
  variables === undefined
    ? ['GrossRevenueRetention.infinite']
    : ['GrossRevenueRetention.infinite', variables];
useGrossRevenueRetentionQuery.fetcher = (
  client: GraphQLClient,
  variables?: GrossRevenueRetentionQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GrossRevenueRetentionQuery, GrossRevenueRetentionQueryVariables>(
    client,
    GrossRevenueRetentionDocument,
    variables,
    headers,
  );
