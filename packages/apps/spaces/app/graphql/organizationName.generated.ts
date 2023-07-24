// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';

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
export type GetOrganizationNameQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type GetOrganizationNameQuery = {
  __typename?: 'Query';
  organization?: { __typename?: 'Organization'; name: string } | null;
};

export const GetOrganizationNameDocument = `
    query GetOrganizationName($id: ID!) {
  organization(id: $id) {
    name
  }
}
    `;
export const useGetOrganizationNameQuery = <
  TData = GetOrganizationNameQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationNameQueryVariables,
  options?: UseQueryOptions<GetOrganizationNameQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetOrganizationNameQuery, TError, TData>(
    ['GetOrganizationName', variables],
    fetcher<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>(
      client,
      GetOrganizationNameDocument,
      variables,
      headers,
    ),
    options,
  );
useGetOrganizationNameQuery.document = GetOrganizationNameDocument;

useGetOrganizationNameQuery.getKey = (
  variables: GetOrganizationNameQueryVariables,
) => ['GetOrganizationName', variables];
useGetOrganizationNameQuery.fetcher = (
  client: GraphQLClient,
  variables: GetOrganizationNameQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>(
    client,
    GetOrganizationNameDocument,
    variables,
    headers,
  );
