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
export type MrrPerCustomerQueryVariables = Types.Exact<{
  period?: Types.InputMaybe<Types.DashboardPeriodInput>;
}>;

export type MrrPerCustomerQuery = {
  __typename?: 'Query';
  dashboard_MRRPerCustomer?: {
    __typename?: 'DashboardMRRPerCustomer';
    mrrPerCustomer: number;
    increasePercentage: string;
    perMonth: Array<{
      __typename?: 'DashboardMRRPerCustomerPerMonth';
      month: number;
      value: number;
    } | null>;
  } | null;
};

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
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: MrrPerCustomerQueryVariables,
  options?: UseQueryOptions<MrrPerCustomerQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<MrrPerCustomerQuery, TError, TData>(
    variables === undefined
      ? ['MrrPerCustomer']
      : ['MrrPerCustomer', variables],
    fetcher<MrrPerCustomerQuery, MrrPerCustomerQueryVariables>(
      client,
      MrrPerCustomerDocument,
      variables,
      headers,
    ),
    options,
  );
useMrrPerCustomerQuery.document = MrrPerCustomerDocument;

useMrrPerCustomerQuery.getKey = (variables?: MrrPerCustomerQueryVariables) =>
  variables === undefined ? ['MrrPerCustomer'] : ['MrrPerCustomer', variables];
export const useInfiniteMrrPerCustomerQuery = <
  TData = MrrPerCustomerQuery,
  TError = unknown,
>(
  pageParamKey: keyof MrrPerCustomerQueryVariables,
  client: GraphQLClient,
  variables?: MrrPerCustomerQueryVariables,
  options?: UseInfiniteQueryOptions<MrrPerCustomerQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<MrrPerCustomerQuery, TError, TData>(
    variables === undefined
      ? ['MrrPerCustomer.infinite']
      : ['MrrPerCustomer.infinite', variables],
    (metaData) =>
      fetcher<MrrPerCustomerQuery, MrrPerCustomerQueryVariables>(
        client,
        MrrPerCustomerDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteMrrPerCustomerQuery.getKey = (
  variables?: MrrPerCustomerQueryVariables,
) =>
  variables === undefined
    ? ['MrrPerCustomer.infinite']
    : ['MrrPerCustomer.infinite', variables];
useMrrPerCustomerQuery.fetcher = (
  client: GraphQLClient,
  variables?: MrrPerCustomerQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<MrrPerCustomerQuery, MrrPerCustomerQueryVariables>(
    client,
    MrrPerCustomerDocument,
    variables,
    headers,
  );
