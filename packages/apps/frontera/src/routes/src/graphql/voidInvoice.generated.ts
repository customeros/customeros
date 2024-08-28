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
export type VoidInvoiceMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type VoidInvoiceMutation = { __typename?: 'Mutation', invoice_Void: { __typename?: 'Invoice', metadata: { __typename?: 'Metadata', id: string } } };



export const VoidInvoiceDocument = `
    mutation VoidInvoice($id: ID!) {
  invoice_Void(id: $id) {
    metadata {
      id
    }
  }
}
    `;

export const useVoidInvoiceMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<VoidInvoiceMutation, TError, VoidInvoiceMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<VoidInvoiceMutation, TError, VoidInvoiceMutationVariables, TContext>(
      {
    mutationKey: ['VoidInvoice'],
    mutationFn: (variables?: VoidInvoiceMutationVariables) => fetcher<VoidInvoiceMutation, VoidInvoiceMutationVariables>(client, VoidInvoiceDocument, variables, headers)(),
    ...options
  }
    )};

useVoidInvoiceMutation.getKey = () => ['VoidInvoice'];


useVoidInvoiceMutation.fetcher = (client: GraphQLClient, variables: VoidInvoiceMutationVariables, headers?: RequestInit['headers']) => fetcher<VoidInvoiceMutation, VoidInvoiceMutationVariables>(client, VoidInvoiceDocument, variables, headers);
