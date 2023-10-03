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
export type SetRelationshipStageMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  relationship: Types.OrganizationRelationship;
  stage: Types.Scalars['String'];
}>;

export type SetRelationshipStageMutation = {
  __typename?: 'Mutation';
  organization_SetRelationshipStage: {
    __typename?: 'Organization';
    id: string;
  };
};

export const SetRelationshipStageDocument = `
    mutation setRelationshipStage($organizationId: ID!, $relationship: OrganizationRelationship!, $stage: String!) {
  organization_SetRelationshipStage(
    organizationId: $organizationId
    relationship: $relationship
    stage: $stage
  ) {
    id
  }
}
    `;
export const useSetRelationshipStageMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    SetRelationshipStageMutation,
    TError,
    SetRelationshipStageMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    SetRelationshipStageMutation,
    TError,
    SetRelationshipStageMutationVariables,
    TContext
  >(
    ['setRelationshipStage'],
    (variables?: SetRelationshipStageMutationVariables) =>
      fetcher<
        SetRelationshipStageMutation,
        SetRelationshipStageMutationVariables
      >(client, SetRelationshipStageDocument, variables, headers)(),
    options,
  );
useSetRelationshipStageMutation.getKey = () => ['setRelationshipStage'];

useSetRelationshipStageMutation.fetcher = (
  client: GraphQLClient,
  variables: SetRelationshipStageMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<SetRelationshipStageMutation, SetRelationshipStageMutationVariables>(
    client,
    SetRelationshipStageDocument,
    variables,
    headers,
  );
