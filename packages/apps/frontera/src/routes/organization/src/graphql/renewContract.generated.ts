// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

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
export type RenewContractMutationVariables = Types.Exact<{
  input: Types.ContractRenewalInput;
}>;


export type RenewContractMutation = { __typename?: 'Mutation', contract_Renew: { __typename?: 'Contract', id: string } };



export const RenewContractDocument = `
    mutation renewContract($input: ContractRenewalInput!) {
  contract_Renew(input: $input) {
    id
  }
}
    `;

export const useRenewContractMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<RenewContractMutation, TError, RenewContractMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<RenewContractMutation, TError, RenewContractMutationVariables, TContext>(
      {
    mutationKey: ['renewContract'],
    mutationFn: (variables?: RenewContractMutationVariables) => fetcher<RenewContractMutation, RenewContractMutationVariables>(client, RenewContractDocument, variables, headers)(),
    ...options
  }
    )};

useRenewContractMutation.getKey = () => ['renewContract'];


useRenewContractMutation.fetcher = (client: GraphQLClient, variables: RenewContractMutationVariables, headers?: RequestInit['headers']) => fetcher<RenewContractMutation, RenewContractMutationVariables>(client, RenewContractDocument, variables, headers);
