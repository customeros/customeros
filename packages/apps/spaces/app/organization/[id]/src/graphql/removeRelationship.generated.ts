// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

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
export type RemoveRelationshipMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  relationship: Types.OrganizationRelationship;
}>;

export type RemoveRelationshipMutation = {
  __typename?: 'Mutation';
  organization_RemoveRelationship: {
    __typename?: 'Organization';
    id: string;
    relationshipStages: Array<{
      __typename?: 'OrganizationRelationshipStage';
      relationship: Types.OrganizationRelationship;
      stage?: string | null;
    }>;
  };
};

export const RemoveRelationshipDocument = `
    mutation removeRelationship($organizationId: ID!, $relationship: OrganizationRelationship!) {
  organization_RemoveRelationship(
    organizationId: $organizationId
    relationship: $relationship
  ) {
    id
    relationshipStages {
      relationship
      stage
    }
  }
}
    `;
export const useRemoveRelationshipMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveRelationshipMutation,
    TError,
    RemoveRelationshipMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    RemoveRelationshipMutation,
    TError,
    RemoveRelationshipMutationVariables,
    TContext
  >(
    ['removeRelationship'],
    (variables?: RemoveRelationshipMutationVariables) =>
      fetcher<RemoveRelationshipMutation, RemoveRelationshipMutationVariables>(
        client,
        RemoveRelationshipDocument,
        variables,
        headers,
      )(),
    options,
  );
useRemoveRelationshipMutation.getKey = () => ['removeRelationship'];

useRemoveRelationshipMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveRelationshipMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<RemoveRelationshipMutation, RemoveRelationshipMutationVariables>(
    client,
    RemoveRelationshipDocument,
    variables,
    headers,
  );
