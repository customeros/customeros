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
export type PayInvoiceMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type PayInvoiceMutation = { __typename?: 'Mutation', invoice_Pay: { __typename?: 'Invoice', metadata: { __typename?: 'Metadata', id: string } } };



export const PayInvoiceDocument = `
    mutation PayInvoice($id: ID!) {
  invoice_Pay(id: $id) {
    metadata {
      id
    }
  }
}
    `;

export const usePayInvoiceMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<PayInvoiceMutation, TError, PayInvoiceMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<PayInvoiceMutation, TError, PayInvoiceMutationVariables, TContext>(
      {
    mutationKey: ['PayInvoice'],
    mutationFn: (variables?: PayInvoiceMutationVariables) => fetcher<PayInvoiceMutation, PayInvoiceMutationVariables>(client, PayInvoiceDocument, variables, headers)(),
    ...options
  }
    )};

usePayInvoiceMutation.getKey = () => ['PayInvoice'];


usePayInvoiceMutation.fetcher = (client: GraphQLClient, variables: PayInvoiceMutationVariables, headers?: RequestInit['headers']) => fetcher<PayInvoiceMutation, PayInvoiceMutationVariables>(client, PayInvoiceDocument, variables, headers);
