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
export type GetCanAccessOrganizationQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type GetCanAccessOrganizationQuery = {
  __typename?: 'Query';
  organization?: { __typename?: 'Organization'; id: string } | null;
};

export const GetCanAccessOrganizationDocument = `
    query GetCanAccessOrganization($id: ID!) {
  organization(id: $id) {
    id
  }
}
    `;
export const useGetCanAccessOrganizationQuery = <
  TData = GetCanAccessOrganizationQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetCanAccessOrganizationQueryVariables,
  options?: UseQueryOptions<GetCanAccessOrganizationQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetCanAccessOrganizationQuery, TError, TData>(
    ['GetCanAccessOrganization', variables],
    fetcher<
      GetCanAccessOrganizationQuery,
      GetCanAccessOrganizationQueryVariables
    >(client, GetCanAccessOrganizationDocument, variables, headers),
    options,
  );
useGetCanAccessOrganizationQuery.document = GetCanAccessOrganizationDocument;

useGetCanAccessOrganizationQuery.getKey = (
  variables: GetCanAccessOrganizationQueryVariables,
) => ['GetCanAccessOrganization', variables];
export const useInfiniteGetCanAccessOrganizationQuery = <
  TData = GetCanAccessOrganizationQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetCanAccessOrganizationQueryVariables,
  client: GraphQLClient,
  variables: GetCanAccessOrganizationQueryVariables,
  options?: UseInfiniteQueryOptions<
    GetCanAccessOrganizationQuery,
    TError,
    TData
  >,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetCanAccessOrganizationQuery, TError, TData>(
    ['GetCanAccessOrganization.infinite', variables],
    (metaData) =>
      fetcher<
        GetCanAccessOrganizationQuery,
        GetCanAccessOrganizationQueryVariables
      >(
        client,
        GetCanAccessOrganizationDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetCanAccessOrganizationQuery.getKey = (
  variables: GetCanAccessOrganizationQueryVariables,
) => ['GetCanAccessOrganization.infinite', variables];
useGetCanAccessOrganizationQuery.fetcher = (
  client: GraphQLClient,
  variables: GetCanAccessOrganizationQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    GetCanAccessOrganizationQuery,
    GetCanAccessOrganizationQueryVariables
  >(client, GetCanAccessOrganizationDocument, variables, headers);
