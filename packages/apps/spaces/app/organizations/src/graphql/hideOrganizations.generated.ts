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
export type HideOrganizationsMutationVariables = Types.Exact<{
  ids: Array<Types.Scalars['ID']> | Types.Scalars['ID'];
}>;

export type HideOrganizationsMutation = {
  __typename?: 'Mutation';
  organization_HideAll?: { __typename?: 'Result'; result: boolean } | null;
};

export const HideOrganizationsDocument = `
    mutation hideOrganizations($ids: [ID!]!) {
  organization_HideAll(ids: $ids) {
    result
  }
}
    `;
export const useHideOrganizationsMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    HideOrganizationsMutation,
    TError,
    HideOrganizationsMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    HideOrganizationsMutation,
    TError,
    HideOrganizationsMutationVariables,
    TContext
  >(
    ['hideOrganizations'],
    (variables?: HideOrganizationsMutationVariables) =>
      fetcher<HideOrganizationsMutation, HideOrganizationsMutationVariables>(
        client,
        HideOrganizationsDocument,
        variables,
        headers,
      )(),
    options,
  );
useHideOrganizationsMutation.getKey = () => ['hideOrganizations'];

useHideOrganizationsMutation.fetcher = (
  client: GraphQLClient,
  variables: HideOrganizationsMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<HideOrganizationsMutation, HideOrganizationsMutationVariables>(
    client,
    HideOrganizationsDocument,
    variables,
    headers,
  );
