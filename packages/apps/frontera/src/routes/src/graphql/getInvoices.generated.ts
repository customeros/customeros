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
export type GetInvoicesQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  organizationId?: Types.InputMaybe<Types.Scalars['ID']['input']>;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Array<Types.SortBy> | Types.SortBy>;
}>;


export type GetInvoicesQuery = { __typename?: 'Query', invoices: { __typename?: 'InvoicesPage', totalPages: number, totalAvailable: any, totalElements: any, content: Array<{ __typename?: 'Invoice', issued: any, invoiceUrl: string, invoiceNumber: string, invoicePeriodStart: any, invoicePeriodEnd: any, due: any, amountDue: number, currency: string, dryRun: boolean, status?: Types.InvoiceStatus | null, metadata: { __typename?: 'Metadata', id: string, created: any }, organization: { __typename?: 'Organization', name: string, metadata: { __typename?: 'Metadata', id: string } }, customer: { __typename?: 'InvoiceCustomer', name?: string | null, email?: string | null }, contract: { __typename?: 'Contract', name: string, contractEnded?: any | null, metadata: { __typename?: 'Metadata', id: string }, billingDetails?: { __typename?: 'BillingDetails', billingCycle?: Types.ContractBillingCycle | null } | null }, invoiceLineItems: Array<{ __typename?: 'InvoiceLine', quantity: any, subtotal: number, taxDue: number, total: number, price: number, description: string, metadata: { __typename?: 'Metadata', id: string, created: any, lastUpdated: any, source: Types.DataSource, sourceOfTruth: Types.DataSource, appSource: string }, contractLineItem: { __typename?: 'ServiceLineItem', serviceStarted: any, billingCycle: Types.BilledType } }> }> } };



export const GetInvoicesDocument = `
    query getInvoices($pagination: Pagination!, $organizationId: ID, $where: Filter, $sort: [SortBy!]) {
  invoices(
    pagination: $pagination
    organizationId: $organizationId
    where: $where
    sort: $sort
  ) {
    content {
      issued
      metadata {
        id
        created
      }
      organization {
        metadata {
          id
        }
        name
      }
      customer {
        name
        email
      }
      contract {
        metadata {
          id
        }
        name
        billingDetails {
          billingCycle
        }
        contractEnded
      }
      invoiceUrl
      invoiceNumber
      invoicePeriodStart
      invoicePeriodEnd
      due
      issued
      amountDue
      currency
      dryRun
      status
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
    totalPages
    totalAvailable
    totalElements
  }
}
    `;

export const useGetInvoicesQuery = <
      TData = GetInvoicesQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetInvoicesQueryVariables,
      options?: Omit<UseQueryOptions<GetInvoicesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetInvoicesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetInvoicesQuery, TError, TData>(
      {
    queryKey: ['getInvoices', variables],
    queryFn: fetcher<GetInvoicesQuery, GetInvoicesQueryVariables>(client, GetInvoicesDocument, variables, headers),
    ...options
  }
    )};

useGetInvoicesQuery.document = GetInvoicesDocument;

useGetInvoicesQuery.getKey = (variables: GetInvoicesQueryVariables) => ['getInvoices', variables];

export const useInfiniteGetInvoicesQuery = <
      TData = InfiniteData<GetInvoicesQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetInvoicesQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetInvoicesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetInvoicesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetInvoicesQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['getInvoices.infinite', variables],
      queryFn: (metaData) => fetcher<GetInvoicesQuery, GetInvoicesQueryVariables>(client, GetInvoicesDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetInvoicesQuery.getKey = (variables: GetInvoicesQueryVariables) => ['getInvoices.infinite', variables];


useGetInvoicesQuery.fetcher = (client: GraphQLClient, variables: GetInvoicesQueryVariables, headers?: RequestInit['headers']) => fetcher<GetInvoicesQuery, GetInvoicesQueryVariables>(client, GetInvoicesDocument, variables, headers);


useGetInvoicesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoicesQueryVariables) =>
  (mutator: (cacheEntry: GetInvoicesQuery) => GetInvoicesQuery) => {
    const cacheKey = useGetInvoicesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetInvoicesQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetInvoicesQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetInvoicesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoicesQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetInvoicesQuery>) => InfiniteData<GetInvoicesQuery>) => {
    const cacheKey = useInfiniteGetInvoicesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetInvoicesQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetInvoicesQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }