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
export type NextDryRunInvoiceForContractMutationVariables = Types.Exact<{
  contractId: Types.Scalars['ID']['input'];
}>;


export type NextDryRunInvoiceForContractMutation = { __typename?: 'Mutation', invoice_NextDryRunForContract: string };



export const NextDryRunInvoiceForContractDocument = `
    mutation NextDryRunInvoiceForContract($contractId: ID!) {
  invoice_NextDryRunForContract(contractId: $contractId)
}
    `;

export const useNextDryRunInvoiceForContractMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<NextDryRunInvoiceForContractMutation, TError, NextDryRunInvoiceForContractMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<NextDryRunInvoiceForContractMutation, TError, NextDryRunInvoiceForContractMutationVariables, TContext>(
      {
    mutationKey: ['NextDryRunInvoiceForContract'],
    mutationFn: (variables?: NextDryRunInvoiceForContractMutationVariables) => fetcher<NextDryRunInvoiceForContractMutation, NextDryRunInvoiceForContractMutationVariables>(client, NextDryRunInvoiceForContractDocument, variables, headers)(),
    ...options
  }
    )};

useNextDryRunInvoiceForContractMutation.getKey = () => ['NextDryRunInvoiceForContract'];


useNextDryRunInvoiceForContractMutation.fetcher = (client: GraphQLClient, variables: NextDryRunInvoiceForContractMutationVariables, headers?: RequestInit['headers']) => fetcher<NextDryRunInvoiceForContractMutation, NextDryRunInvoiceForContractMutationVariables>(client, NextDryRunInvoiceForContractDocument, variables, headers);
