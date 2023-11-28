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
export type CustomerMapQueryVariables = Types.Exact<{ [key: string]: never }>;

export type CustomerMapQuery = {
  __typename?: 'Query';
  dashboard_CustomerMap?: Array<{
    __typename?: 'DashboardCustomerMap';
    state: Types.DashboardCustomerMapState;
    arr: number;
    contractSignedDate: any;
    organization: { __typename?: 'Organization'; id: string; name: string };
  }> | null;
};

export const CustomerMapDocument = `
    query CustomerMap {
  dashboard_CustomerMap {
    organization {
      id
      name
    }
    state
    arr
    contractSignedDate
  }
}
    `;
export const useCustomerMapQuery = <TData = CustomerMapQuery, TError = unknown>(
  client: GraphQLClient,
  variables?: CustomerMapQueryVariables,
  options?: UseQueryOptions<CustomerMapQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<CustomerMapQuery, TError, TData>(
    variables === undefined ? ['CustomerMap'] : ['CustomerMap', variables],
    fetcher<CustomerMapQuery, CustomerMapQueryVariables>(
      client,
      CustomerMapDocument,
      variables,
      headers,
    ),
    options,
  );
useCustomerMapQuery.document = CustomerMapDocument;

useCustomerMapQuery.getKey = (variables?: CustomerMapQueryVariables) =>
  variables === undefined ? ['CustomerMap'] : ['CustomerMap', variables];
export const useInfiniteCustomerMapQuery = <
  TData = CustomerMapQuery,
  TError = unknown,
>(
  pageParamKey: keyof CustomerMapQueryVariables,
  client: GraphQLClient,
  variables?: CustomerMapQueryVariables,
  options?: UseInfiniteQueryOptions<CustomerMapQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<CustomerMapQuery, TError, TData>(
    variables === undefined
      ? ['CustomerMap.infinite']
      : ['CustomerMap.infinite', variables],
    (metaData) =>
      fetcher<CustomerMapQuery, CustomerMapQueryVariables>(
        client,
        CustomerMapDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteCustomerMapQuery.getKey = (variables?: CustomerMapQueryVariables) =>
  variables === undefined
    ? ['CustomerMap.infinite']
    : ['CustomerMap.infinite', variables];
useCustomerMapQuery.fetcher = (
  client: GraphQLClient,
  variables?: CustomerMapQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CustomerMapQuery, CustomerMapQueryVariables>(
    client,
    CustomerMapDocument,
    variables,
    headers,
  );
