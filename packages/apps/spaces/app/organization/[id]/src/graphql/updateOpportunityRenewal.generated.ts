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
export type UpdateOpportunityRenewalMutationVariables = Types.Exact<{
  input: Types.OpportunityRenewalUpdateInput;
}>;

export type UpdateOpportunityRenewalMutation = {
  __typename?: 'Mutation';
  opportunityRenewalUpdate: { __typename?: 'Opportunity'; id: string };
};

export const UpdateOpportunityRenewalDocument = `
    mutation updateOpportunityRenewal($input: OpportunityRenewalUpdateInput!) {
  opportunityRenewalUpdate(input: $input) {
    id
  }
}
    `;
export const useUpdateOpportunityRenewalMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateOpportunityRenewalMutation,
    TError,
    UpdateOpportunityRenewalMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateOpportunityRenewalMutation,
    TError,
    UpdateOpportunityRenewalMutationVariables,
    TContext
  >(
    ['updateOpportunityRenewal'],
    (variables?: UpdateOpportunityRenewalMutationVariables) =>
      fetcher<
        UpdateOpportunityRenewalMutation,
        UpdateOpportunityRenewalMutationVariables
      >(client, UpdateOpportunityRenewalDocument, variables, headers)(),
    options,
  );
useUpdateOpportunityRenewalMutation.getKey = () => ['updateOpportunityRenewal'];

useUpdateOpportunityRenewalMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateOpportunityRenewalMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    UpdateOpportunityRenewalMutation,
    UpdateOpportunityRenewalMutationVariables
  >(client, UpdateOpportunityRenewalDocument, variables, headers);
