// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

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
export type GetInvoiceQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type GetInvoiceQuery = {
  __typename?: 'Query';
  invoice: {
    __typename?: 'Invoice';
    id: string;
    createdAt: any;
    status?: Types.InvoiceStatus | null;
    number: string;
    periodStartDate: any;
    periodEndDate: any;
    dueDate: any;
    amount: number;
    vat: number;
    total: number;
    currency: string;
    repositoryFileId: string;
    invoiceLines: Array<{
      __typename?: 'InvoiceLine';
      id: string;
      createdAt: any;
      quantity: number;
      amount: number;
      vat: number;
      total: number;
      price: number;
      name: string;
    }>;
  };
};

export const GetInvoiceDocument = `
    query GetInvoice($id: ID!) {
  invoice(id: $id) {
    id
    createdAt
    status
    number
    periodStartDate
    periodEndDate
    dueDate
    amount
    vat
    total
    currency
    repositoryFileId
    invoiceLines {
      id
      createdAt
      quantity
      amount
      vat
      total
      price
      name
    }
  }
}
    `;

export const useGetInvoiceQuery = <TData = GetInvoiceQuery, TError = unknown>(
  client: GraphQLClient,
  variables: GetInvoiceQueryVariables,
  options?: Omit<
    UseQueryOptions<GetInvoiceQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<GetInvoiceQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetInvoiceQuery, TError, TData>({
    queryKey: ['GetInvoice', variables],
    queryFn: fetcher<GetInvoiceQuery, GetInvoiceQueryVariables>(
      client,
      GetInvoiceDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGetInvoiceQuery.document = GetInvoiceDocument;

useGetInvoiceQuery.getKey = (variables: GetInvoiceQueryVariables) => [
  'GetInvoice',
  variables,
];

export const useInfiniteGetInvoiceQuery = <
  TData = InfiniteData<GetInvoiceQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetInvoiceQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetInvoiceQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetInvoiceQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetInvoiceQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? ['GetInvoice.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GetInvoiceQuery, GetInvoiceQueryVariables>(
            client,
            GetInvoiceDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetInvoiceQuery.getKey = (variables: GetInvoiceQueryVariables) => [
  'GetInvoice.infinite',
  variables,
];

useGetInvoiceQuery.fetcher = (
  client: GraphQLClient,
  variables: GetInvoiceQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetInvoiceQuery, GetInvoiceQueryVariables>(
    client,
    GetInvoiceDocument,
    variables,
    headers,
  );

useGetInvoiceQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoiceQueryVariables) =>
  (mutator: (cacheEntry: GetInvoiceQuery) => GetInvoiceQuery) => {
    const cacheKey = useGetInvoiceQuery.getKey(variables);
    const previousEntries = queryClient.getQueryData<GetInvoiceQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetInvoiceQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetInvoiceQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoiceQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetInvoiceQuery>,
    ) => InfiniteData<GetInvoiceQuery>,
  ) => {
    const cacheKey = useInfiniteGetInvoiceQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetInvoiceQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetInvoiceQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
