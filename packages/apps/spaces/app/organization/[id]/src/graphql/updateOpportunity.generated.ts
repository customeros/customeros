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
export type UpdateOpportunityMutationVariables = Types.Exact<{
  input: Types.OpportunityUpdateInput;
}>;

export type UpdateOpportunityMutation = {
  __typename?: 'Mutation';
  opportunityUpdate: { __typename?: 'Opportunity'; id: string };
};

export const UpdateOpportunityDocument = `
    mutation updateOpportunity($input: OpportunityUpdateInput!) {
  opportunityUpdate(input: $input) {
    id
  }
}
    `;
export const useUpdateOpportunityMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateOpportunityMutation,
    TError,
    UpdateOpportunityMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateOpportunityMutation,
    TError,
    UpdateOpportunityMutationVariables,
    TContext
  >(
    ['updateOpportunity'],
    (variables?: UpdateOpportunityMutationVariables) =>
      fetcher<UpdateOpportunityMutation, UpdateOpportunityMutationVariables>(
        client,
        UpdateOpportunityDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateOpportunityMutation.getKey = () => ['updateOpportunity'];

useUpdateOpportunityMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateOpportunityMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateOpportunityMutation, UpdateOpportunityMutationVariables>(
    client,
    UpdateOpportunityDocument,
    variables,
    headers,
  );
