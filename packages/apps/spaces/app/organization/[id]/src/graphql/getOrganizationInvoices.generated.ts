// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
  InfiniteData,
} from '@tanstack/react-query';

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
export type GetOrganizationInvoicesQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  pagination: Types.Pagination;
}>;

export type GetOrganizationInvoicesQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    invoices: {
      __typename?: 'InvoicesPage';
      totalPages: number;
      totalElements: any;
      content: Array<{
        __typename?: 'Invoice';
        id: string;
        number: string;
        dueDate: any;
        total: number;
        currency: string;
        createdAt: any;
        status?: Types.InvoiceStatus | null;
        dryRun: boolean;
        invoiceLines: Array<{
          __typename?: 'InvoiceLine';
          id: string;
          name: string;
          price: number;
          quantity: number;
          amount: number;
          vat: number;
          total: number;
        }>;
      }>;
    };
  } | null;
};

export const GetOrganizationInvoicesDocument = `
    query getOrganizationInvoices($id: ID!, $pagination: Pagination!) {
  organization(id: $id) {
    invoices(pagination: $pagination) {
      content {
        id
        number
        dueDate
        total
        currency
        createdAt
        status
        dryRun
        invoiceLines {
          id
          name
          price
          quantity
          amount
          vat
          total
        }
      }
      totalPages
      totalElements
    }
  }
}
    `;

export const useGetOrganizationInvoicesQuery = <
  TData = GetOrganizationInvoicesQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationInvoicesQueryVariables,
  options?: Omit<
    UseQueryOptions<GetOrganizationInvoicesQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<
      GetOrganizationInvoicesQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetOrganizationInvoicesQuery, TError, TData>({
    queryKey: ['getOrganizationInvoices', variables],
    queryFn: fetcher<
      GetOrganizationInvoicesQuery,
      GetOrganizationInvoicesQueryVariables
    >(client, GetOrganizationInvoicesDocument, variables, headers),
    ...options,
  });
};

useGetOrganizationInvoicesQuery.document = GetOrganizationInvoicesDocument;

useGetOrganizationInvoicesQuery.getKey = (
  variables: GetOrganizationInvoicesQueryVariables,
) => ['getOrganizationInvoices', variables];

export const useInfiniteGetOrganizationInvoicesQuery = <
  TData = InfiniteData<GetOrganizationInvoicesQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationInvoicesQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetOrganizationInvoicesQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetOrganizationInvoicesQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetOrganizationInvoicesQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? [
          'getOrganizationInvoices.infinite',
          variables,
        ],
        queryFn: (metaData) =>
          fetcher<
            GetOrganizationInvoicesQuery,
            GetOrganizationInvoicesQueryVariables
          >(
            client,
            GetOrganizationInvoicesDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetOrganizationInvoicesQuery.getKey = (
  variables: GetOrganizationInvoicesQueryVariables,
) => ['getOrganizationInvoices.infinite', variables];

useGetOrganizationInvoicesQuery.fetcher = (
  client: GraphQLClient,
  variables: GetOrganizationInvoicesQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetOrganizationInvoicesQuery, GetOrganizationInvoicesQueryVariables>(
    client,
    GetOrganizationInvoicesDocument,
    variables,
    headers,
  );

useGetOrganizationInvoicesQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables?: GetOrganizationInvoicesQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: GetOrganizationInvoicesQuery,
    ) => GetOrganizationInvoicesQuery,
  ) => {
    const cacheKey = useGetOrganizationInvoicesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetOrganizationInvoicesQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetOrganizationInvoicesQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetOrganizationInvoicesQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables?: GetOrganizationInvoicesQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetOrganizationInvoicesQuery>,
    ) => InfiniteData<GetOrganizationInvoicesQuery>,
  ) => {
    const cacheKey = useInfiniteGetOrganizationInvoicesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetOrganizationInvoicesQuery>>(
        cacheKey,
      );
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetOrganizationInvoicesQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
