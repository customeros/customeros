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
export type RemoveRelationshipStageMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  relationship: Types.OrganizationRelationship;
}>;

export type RemoveRelationshipStageMutation = {
  __typename?: 'Mutation';
  organization_RemoveRelationshipStage: {
    __typename?: 'Organization';
    id: string;
  };
};

export const RemoveRelationshipStageDocument = `
    mutation removeRelationshipStage($organizationId: ID!, $relationship: OrganizationRelationship!) {
  organization_RemoveRelationshipStage(
    organizationId: $organizationId
    relationship: $relationship
  ) {
    id
  }
}
    `;
export const useRemoveRelationshipStageMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveRelationshipStageMutation,
    TError,
    RemoveRelationshipStageMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    RemoveRelationshipStageMutation,
    TError,
    RemoveRelationshipStageMutationVariables,
    TContext
  >(
    ['removeRelationshipStage'],
    (variables?: RemoveRelationshipStageMutationVariables) =>
      fetcher<
        RemoveRelationshipStageMutation,
        RemoveRelationshipStageMutationVariables
      >(client, RemoveRelationshipStageDocument, variables, headers)(),
    options,
  );
useRemoveRelationshipStageMutation.getKey = () => ['removeRelationshipStage'];

useRemoveRelationshipStageMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveRelationshipStageMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    RemoveRelationshipStageMutation,
    RemoveRelationshipStageMutationVariables
  >(client, RemoveRelationshipStageDocument, variables, headers);
