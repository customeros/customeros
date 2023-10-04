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
export type SetOrganizationOwnerMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  userId: Types.Scalars['ID'];
}>;

export type SetOrganizationOwnerMutation = {
  __typename?: 'Mutation';
  organization_SetOwner: {
    __typename?: 'Organization';
    id: string;
    owner?: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
    } | null;
  };
};

export const SetOrganizationOwnerDocument = `
    mutation setOrganizationOwner($organizationId: ID!, $userId: ID!) {
  organization_SetOwner(organizationId: $organizationId, userId: $userId) {
    id
    owner {
      id
      firstName
      lastName
    }
  }
}
    `;
export const useSetOrganizationOwnerMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    SetOrganizationOwnerMutation,
    TError,
    SetOrganizationOwnerMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    SetOrganizationOwnerMutation,
    TError,
    SetOrganizationOwnerMutationVariables,
    TContext
  >(
    ['setOrganizationOwner'],
    (variables?: SetOrganizationOwnerMutationVariables) =>
      fetcher<
        SetOrganizationOwnerMutation,
        SetOrganizationOwnerMutationVariables
      >(client, SetOrganizationOwnerDocument, variables, headers)(),
    options,
  );
useSetOrganizationOwnerMutation.getKey = () => ['setOrganizationOwner'];

useSetOrganizationOwnerMutation.fetcher = (
  client: GraphQLClient,
  variables: SetOrganizationOwnerMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<SetOrganizationOwnerMutation, SetOrganizationOwnerMutationVariables>(
    client,
    SetOrganizationOwnerDocument,
    variables,
    headers,
  );
