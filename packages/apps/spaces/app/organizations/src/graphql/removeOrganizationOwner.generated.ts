// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

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
export type RemoveOrganizationOwnerMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
}>;

export type RemoveOrganizationOwnerMutation = {
  __typename?: 'Mutation';
  organization_UnsetOwner: {
    __typename?: 'Organization';
    id: string;
    owner?: { __typename?: 'User'; id: string } | null;
  };
};

export const RemoveOrganizationOwnerDocument = `
    mutation removeOrganizationOwner($organizationId: ID!) {
  organization_UnsetOwner(organizationId: $organizationId) {
    id
    owner {
      id
    }
  }
}
    `;
export const useRemoveOrganizationOwnerMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveOrganizationOwnerMutation,
    TError,
    RemoveOrganizationOwnerMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    RemoveOrganizationOwnerMutation,
    TError,
    RemoveOrganizationOwnerMutationVariables,
    TContext
  >(
    ['removeOrganizationOwner'],
    (variables?: RemoveOrganizationOwnerMutationVariables) =>
      fetcher<
        RemoveOrganizationOwnerMutation,
        RemoveOrganizationOwnerMutationVariables
      >(client, RemoveOrganizationOwnerDocument, variables, headers)(),
    options,
  );
useRemoveOrganizationOwnerMutation.getKey = () => ['removeOrganizationOwner'];

useRemoveOrganizationOwnerMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveOrganizationOwnerMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    RemoveOrganizationOwnerMutation,
    RemoveOrganizationOwnerMutationVariables
  >(client, RemoveOrganizationOwnerDocument, variables, headers);
