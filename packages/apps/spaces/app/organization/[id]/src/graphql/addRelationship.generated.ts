// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../types/__generated__/graphql.types';

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
export type AddRelationshipMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  relationship: Types.OrganizationRelationship;
}>;

export type AddRelationshipMutation = {
  __typename?: 'Mutation';
  organization_AddRelationship: { __typename?: 'Organization'; id: string };
};

export const AddRelationshipDocument = `
    mutation addRelationship($organizationId: ID!, $relationship: OrganizationRelationship!) {
  organization_AddRelationship(
    organizationId: $organizationId
    relationship: $relationship
  ) {
    id
  }
}
    `;
export const useAddRelationshipMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddRelationshipMutation,
    TError,
    AddRelationshipMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddRelationshipMutation,
    TError,
    AddRelationshipMutationVariables,
    TContext
  >(
    ['addRelationship'],
    (variables?: AddRelationshipMutationVariables) =>
      fetcher<AddRelationshipMutation, AddRelationshipMutationVariables>(
        client,
        AddRelationshipDocument,
        variables,
        headers,
      )(),
    options,
  );
useAddRelationshipMutation.getKey = () => ['addRelationship'];

useAddRelationshipMutation.fetcher = (
  client: GraphQLClient,
  variables: AddRelationshipMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<AddRelationshipMutation, AddRelationshipMutationVariables>(
    client,
    AddRelationshipDocument,
    variables,
    headers,
  );
