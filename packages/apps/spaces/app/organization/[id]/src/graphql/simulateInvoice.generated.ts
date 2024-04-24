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
export type SimulateInvoiceMutationVariables = Types.Exact<{
  input: Types.InvoiceSimulateInput;
}>;

export type SimulateInvoiceMutation = {
  __typename?: 'Mutation';
  invoice_Simulate: Array<{
    __typename?: 'InvoiceSimulate';
    postpaid: boolean;
    offCycle: boolean;
    invoiceNumber: string;
    invoicePeriodEnd: any;
    invoicePeriodStart: any;
    amount: number;
    subtotal: number;
    total: number;
    due: any;
    issued: any;
    note: string;
    currency: string;
    invoiceLineItems: Array<{
      __typename?: 'InvoiceLineSimulate';
      key: string;
      description: string;
      price: number;
      quantity: any;
      subtotal: number;
      taxDue: number;
      total: number;
    }>;
    customer: {
      __typename?: 'InvoiceCustomer';
      name?: string | null;
      email?: string | null;
      addressLine1?: string | null;
      addressLine2?: string | null;
      addressZip?: string | null;
      addressLocality?: string | null;
      addressCountry?: string | null;
      addressRegion?: string | null;
    };
    provider: {
      __typename?: 'InvoiceProvider';
      name?: string | null;
      addressLine1?: string | null;
      addressLine2?: string | null;
      addressZip?: string | null;
      addressLocality?: string | null;
      addressCountry?: string | null;
      addressRegion?: string | null;
    };
  }>;
};

export const SimulateInvoiceDocument = `
    mutation simulateInvoice($input: InvoiceSimulateInput!) {
  invoice_Simulate(input: $input) {
    postpaid
    offCycle
    invoiceNumber
    invoicePeriodEnd
    invoicePeriodStart
    amount
    subtotal
    total
    due
    issued
    invoiceLineItems {
      key
      description
      price
      quantity
      subtotal
      taxDue
      total
    }
    note
    currency
    customer {
      name
      email
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
      addressRegion
    }
    provider {
      name
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
      addressRegion
    }
  }
}
    `;

export const useSimulateInvoiceMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    SimulateInvoiceMutation,
    TError,
    SimulateInvoiceMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    SimulateInvoiceMutation,
    TError,
    SimulateInvoiceMutationVariables,
    TContext
  >({
    mutationKey: ['simulateInvoice'],
    mutationFn: (variables?: SimulateInvoiceMutationVariables) =>
      fetcher<SimulateInvoiceMutation, SimulateInvoiceMutationVariables>(
        client,
        SimulateInvoiceDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useSimulateInvoiceMutation.getKey = () => ['simulateInvoice'];

useSimulateInvoiceMutation.fetcher = (
  client: GraphQLClient,
  variables: SimulateInvoiceMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<SimulateInvoiceMutation, SimulateInvoiceMutationVariables>(
    client,
    SimulateInvoiceDocument,
    variables,
    headers,
  );
