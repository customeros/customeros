// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useMutation, UseMutationOptions } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type BulkUpdateOpportunityRenewalMutationVariables = Types.Exact<{
  input: Types.OpportunityRenewalUpdateAllForOrganizationInput;
}>;


export type BulkUpdateOpportunityRenewalMutation = { __typename?: 'Mutation', opportunityRenewal_UpdateAllForOrganization: { __typename?: 'Organization', metadata: { __typename?: 'Metadata', id: string } } };



export const BulkUpdateOpportunityRenewalDocument = `
    mutation bulkUpdateOpportunityRenewal($input: OpportunityRenewalUpdateAllForOrganizationInput!) {
  opportunityRenewal_UpdateAllForOrganization(input: $input) {
    metadata {
      id
    }
  }
}
    `;

export const useBulkUpdateOpportunityRenewalMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<BulkUpdateOpportunityRenewalMutation, TError, BulkUpdateOpportunityRenewalMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<BulkUpdateOpportunityRenewalMutation, TError, BulkUpdateOpportunityRenewalMutationVariables, TContext>(
      {
    mutationKey: ['bulkUpdateOpportunityRenewal'],
    mutationFn: (variables?: BulkUpdateOpportunityRenewalMutationVariables) => fetcher<BulkUpdateOpportunityRenewalMutation, BulkUpdateOpportunityRenewalMutationVariables>(client, BulkUpdateOpportunityRenewalDocument, variables, headers)(),
    ...options
  }
    )};

useBulkUpdateOpportunityRenewalMutation.getKey = () => ['bulkUpdateOpportunityRenewal'];


useBulkUpdateOpportunityRenewalMutation.fetcher = (client: GraphQLClient, variables: BulkUpdateOpportunityRenewalMutationVariables, headers?: RequestInit['headers']) => fetcher<BulkUpdateOpportunityRenewalMutation, BulkUpdateOpportunityRenewalMutationVariables>(client, BulkUpdateOpportunityRenewalDocument, variables, headers);
