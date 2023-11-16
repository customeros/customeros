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
export type GetContractsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type GetContractsQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    id: string;
    name: string;
    contracts?: Array<{
      __typename?: 'Contract';
      id: string;
      name: string;
      createdAt: any;
      serviceStartedAt?: any | null;
      signedAt?: any | null;
      endedAt?: any | null;
      renewalCycle: Types.ContractRenewalCycle;
      status: Types.ContractStatus;
      contractUrl?: string | null;
    }> | null;
  } | null;
};

export const GetContractsDocument = `
    query getContracts($id: ID!) {
  organization(id: $id) {
    id
    name
    contracts {
      id
      name
      createdAt
      serviceStartedAt
      signedAt
      endedAt
      renewalCycle
      status
      contractUrl
    }
  }
}
    `;
export const useGetContractsQuery = <
  TData = GetContractsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetContractsQueryVariables,
  options?: UseQueryOptions<GetContractsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetContractsQuery, TError, TData>(
    ['getContracts', variables],
    fetcher<GetContractsQuery, GetContractsQueryVariables>(
      client,
      GetContractsDocument,
      variables,
      headers,
    ),
    options,
  );
useGetContractsQuery.document = GetContractsDocument;

useGetContractsQuery.getKey = (variables: GetContractsQueryVariables) => [
  'getContracts',
  variables,
];
export const useInfiniteGetContractsQuery = <
  TData = GetContractsQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetContractsQueryVariables,
  client: GraphQLClient,
  variables: GetContractsQueryVariables,
  options?: UseInfiniteQueryOptions<GetContractsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetContractsQuery, TError, TData>(
    ['getContracts.infinite', variables],
    (metaData) =>
      fetcher<GetContractsQuery, GetContractsQueryVariables>(
        client,
        GetContractsDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetContractsQuery.getKey = (
  variables: GetContractsQueryVariables,
) => ['getContracts.infinite', variables];
useGetContractsQuery.fetcher = (
  client: GraphQLClient,
  variables: GetContractsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetContractsQuery, GetContractsQueryVariables>(
    client,
    GetContractsDocument,
    variables,
    headers,
  );
