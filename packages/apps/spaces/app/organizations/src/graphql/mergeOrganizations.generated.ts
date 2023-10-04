// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useMutation, UseMutationOptions } from '@tanstack/react-query';

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
export type MergeOrganizationsMutationVariables = Types.Exact<{
  primaryOrganizationId: Types.Scalars['ID'];
  mergedOrganizationIds: Array<Types.Scalars['ID']> | Types.Scalars['ID'];
}>;

export type MergeOrganizationsMutation = {
  __typename?: 'Mutation';
  organization_Merge: { __typename?: 'Organization'; id: string };
};

export const MergeOrganizationsDocument = `
    mutation mergeOrganizations($primaryOrganizationId: ID!, $mergedOrganizationIds: [ID!]!) {
  organization_Merge(
    primaryOrganizationId: $primaryOrganizationId
    mergedOrganizationIds: $mergedOrganizationIds
  ) {
    id
  }
}
    `;
export const useMergeOrganizationsMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    MergeOrganizationsMutation,
    TError,
    MergeOrganizationsMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    MergeOrganizationsMutation,
    TError,
    MergeOrganizationsMutationVariables,
    TContext
  >(
    ['mergeOrganizations'],
    (variables?: MergeOrganizationsMutationVariables) =>
      fetcher<MergeOrganizationsMutation, MergeOrganizationsMutationVariables>(
        client,
        MergeOrganizationsDocument,
        variables,
        headers,
      )(),
    options,
  );
useMergeOrganizationsMutation.getKey = () => ['mergeOrganizations'];

useMergeOrganizationsMutation.fetcher = (
  client: GraphQLClient,
  variables: MergeOrganizationsMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<MergeOrganizationsMutation, MergeOrganizationsMutationVariables>(
    client,
    MergeOrganizationsDocument,
    variables,
    headers,
  );
