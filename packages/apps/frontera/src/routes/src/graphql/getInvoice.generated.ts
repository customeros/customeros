// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, useInfiniteQuery, UseQueryOptions, UseInfiniteQueryOptions, InfiniteData } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type GetInvoiceQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type GetInvoiceQuery = { __typename?: 'Query', invoice: { __typename?: 'Invoice', status?: Types.InvoiceStatus | null, invoiceNumber: string, invoicePeriodStart: any, invoicePeriodEnd: any, invoiceUrl: string, due: any, issued: any, subtotal: number, taxDue: number, amountDue: number, currency: string, note?: string | null, repositoryFileId: string, metadata: { __typename?: 'Metadata', id: string, created: any }, contract: { __typename?: 'Contract', billingDetails?: { __typename?: 'BillingDetails', canPayWithBankTransfer?: boolean | null } | null }, customer: { __typename?: 'InvoiceCustomer', name?: string | null, email?: string | null, addressLine1?: string | null, addressLine2?: string | null, addressZip?: string | null, addressLocality?: string | null, addressCountry?: string | null }, provider: { __typename?: 'InvoiceProvider', name?: string | null, logoUrl?: string | null, addressLine1?: string | null, addressLine2?: string | null, addressZip?: string | null, addressLocality?: string | null, addressCountry?: string | null }, invoiceLineItems: Array<{ __typename?: 'InvoiceLine', quantity: any, subtotal: number, taxDue: number, total: number, price: number, description: string, metadata: { __typename?: 'Metadata', id: string, created: any, lastUpdated: any, source: Types.DataSource, sourceOfTruth: Types.DataSource, appSource: string }, contractLineItem: { __typename?: 'ServiceLineItem', serviceStarted: any, billingCycle: Types.BilledType } }> } };



export const GetInvoiceDocument = `
    query GetInvoice($id: ID!) {
  invoice(id: $id) {
    metadata {
      id
      created
    }
    contract {
      billingDetails {
        canPayWithBankTransfer
      }
    }
    status
    invoiceNumber
    invoicePeriodStart
    invoicePeriodEnd
    invoiceUrl
    due
    issued
    subtotal
    taxDue
    amountDue
    currency
    note
    repositoryFileId
    customer {
      name
      email
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
    }
    provider {
      name
      logoUrl
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
    }
    invoiceLineItems {
      metadata {
        id
        created
        lastUpdated
        source
        sourceOfTruth
        appSource
      }
      contractLineItem {
        serviceStarted
        billingCycle
      }
      quantity
      subtotal
      taxDue
      total
      price
      description
    }
  }
}
    `;

export const useGetInvoiceQuery = <
      TData = GetInvoiceQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetInvoiceQueryVariables,
      options?: Omit<UseQueryOptions<GetInvoiceQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetInvoiceQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetInvoiceQuery, TError, TData>(
      {
    queryKey: ['GetInvoice', variables],
    queryFn: fetcher<GetInvoiceQuery, GetInvoiceQueryVariables>(client, GetInvoiceDocument, variables, headers),
    ...options
  }
    )};

useGetInvoiceQuery.document = GetInvoiceDocument;

useGetInvoiceQuery.getKey = (variables: GetInvoiceQueryVariables) => ['GetInvoice', variables];

export const useInfiniteGetInvoiceQuery = <
      TData = InfiniteData<GetInvoiceQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetInvoiceQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetInvoiceQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetInvoiceQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetInvoiceQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['GetInvoice.infinite', variables],
      queryFn: (metaData) => fetcher<GetInvoiceQuery, GetInvoiceQueryVariables>(client, GetInvoiceDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetInvoiceQuery.getKey = (variables: GetInvoiceQueryVariables) => ['GetInvoice.infinite', variables];


useGetInvoiceQuery.fetcher = (client: GraphQLClient, variables: GetInvoiceQueryVariables, headers?: RequestInit['headers']) => fetcher<GetInvoiceQuery, GetInvoiceQueryVariables>(client, GetInvoiceDocument, variables, headers);


useGetInvoiceQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoiceQueryVariables) =>
  (mutator: (cacheEntry: GetInvoiceQuery) => GetInvoiceQuery) => {
    const cacheKey = useGetInvoiceQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetInvoiceQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetInvoiceQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetInvoiceQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoiceQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetInvoiceQuery>) => InfiniteData<GetInvoiceQuery>) => {
    const cacheKey = useInfiniteGetInvoiceQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetInvoiceQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetInvoiceQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }