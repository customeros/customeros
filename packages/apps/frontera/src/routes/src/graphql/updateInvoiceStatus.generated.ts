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
export type UpdateInvoiceStatusMutationVariables = Types.Exact<{
  input: Types.InvoiceUpdateInput;
}>;


export type UpdateInvoiceStatusMutation = { __typename?: 'Mutation', invoice_Update: { __typename?: 'Invoice', metadata: { __typename?: 'Metadata', id: string } } };



export const UpdateInvoiceStatusDocument = `
    mutation UpdateInvoiceStatus($input: InvoiceUpdateInput!) {
  invoice_Update(input: $input) {
    metadata {
      id
    }
  }
}
    `;

export const useUpdateInvoiceStatusMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<UpdateInvoiceStatusMutation, TError, UpdateInvoiceStatusMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<UpdateInvoiceStatusMutation, TError, UpdateInvoiceStatusMutationVariables, TContext>(
      {
    mutationKey: ['UpdateInvoiceStatus'],
    mutationFn: (variables?: UpdateInvoiceStatusMutationVariables) => fetcher<UpdateInvoiceStatusMutation, UpdateInvoiceStatusMutationVariables>(client, UpdateInvoiceStatusDocument, variables, headers)(),
    ...options
  }
    )};

useUpdateInvoiceStatusMutation.getKey = () => ['UpdateInvoiceStatus'];


useUpdateInvoiceStatusMutation.fetcher = (client: GraphQLClient, variables: UpdateInvoiceStatusMutationVariables, headers?: RequestInit['headers']) => fetcher<UpdateInvoiceStatusMutation, UpdateInvoiceStatusMutationVariables>(client, UpdateInvoiceStatusDocument, variables, headers);
