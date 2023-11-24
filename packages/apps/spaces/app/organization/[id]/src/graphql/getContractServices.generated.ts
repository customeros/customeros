// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

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
export type GetContractServicesQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type GetContractServicesQuery = {
  __typename?: 'Query';
  serviceLineItem: { __typename?: 'ServiceLineItem'; id: string; name: string };
};

export const GetContractServicesDocument = `
    query getContractServices($id: ID!) {
  serviceLineItem(id: $id) {
    id
    name
  }
}
    `;
export const useGetContractServicesQuery = <
  TData = GetContractServicesQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetContractServicesQueryVariables,
  options?: UseQueryOptions<GetContractServicesQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetContractServicesQuery, TError, TData>(
    ['getContractServices', variables],
    fetcher<GetContractServicesQuery, GetContractServicesQueryVariables>(
      client,
      GetContractServicesDocument,
      variables,
      headers,
    ),
    options,
  );
useGetContractServicesQuery.document = GetContractServicesDocument;

useGetContractServicesQuery.getKey = (
  variables: GetContractServicesQueryVariables,
) => ['getContractServices', variables];
export const useInfiniteGetContractServicesQuery = <
  TData = GetContractServicesQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetContractServicesQueryVariables,
  client: GraphQLClient,
  variables: GetContractServicesQueryVariables,
  options?: UseInfiniteQueryOptions<GetContractServicesQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetContractServicesQuery, TError, TData>(
    ['getContractServices.infinite', variables],
    (metaData) =>
      fetcher<GetContractServicesQuery, GetContractServicesQueryVariables>(
        client,
        GetContractServicesDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetContractServicesQuery.getKey = (
  variables: GetContractServicesQueryVariables,
) => ['getContractServices.infinite', variables];
useGetContractServicesQuery.fetcher = (
  client: GraphQLClient,
  variables: GetContractServicesQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetContractServicesQuery, GetContractServicesQueryVariables>(
    client,
    GetContractServicesDocument,
    variables,
    headers,
  );
