// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

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
export type BankAccountsQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type BankAccountsQuery = { __typename?: 'Query', bankAccounts: Array<{ __typename?: 'BankAccount', bankName?: string | null, currency?: Types.Currency | null, bankTransferEnabled: boolean, iban?: string | null, bic?: string | null, sortCode?: string | null, accountNumber?: string | null, routingNumber?: string | null, allowInternational: boolean, otherDetails?: string | null, metadata: { __typename?: 'Metadata', id: string, created: any, lastUpdated: any } }> };



export const BankAccountsDocument = `
    query BankAccounts {
  bankAccounts {
    metadata {
      id
      created
      lastUpdated
    }
    bankName
    currency
    bankTransferEnabled
    iban
    bic
    sortCode
    accountNumber
    routingNumber
    allowInternational
    otherDetails
  }
}
    `;

export const useBankAccountsQuery = <
      TData = BankAccountsQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: BankAccountsQueryVariables,
      options?: Omit<UseQueryOptions<BankAccountsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<BankAccountsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<BankAccountsQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['BankAccounts'] : ['BankAccounts', variables],
    queryFn: fetcher<BankAccountsQuery, BankAccountsQueryVariables>(client, BankAccountsDocument, variables, headers),
    ...options
  }
    )};

useBankAccountsQuery.document = BankAccountsDocument;

useBankAccountsQuery.getKey = (variables?: BankAccountsQueryVariables) => variables === undefined ? ['BankAccounts'] : ['BankAccounts', variables];

export const useInfiniteBankAccountsQuery = <
      TData = InfiniteData<BankAccountsQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: BankAccountsQueryVariables,
      options: Omit<UseInfiniteQueryOptions<BankAccountsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<BankAccountsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<BankAccountsQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['BankAccounts.infinite'] : ['BankAccounts.infinite', variables],
      queryFn: (metaData) => fetcher<BankAccountsQuery, BankAccountsQueryVariables>(client, BankAccountsDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteBankAccountsQuery.getKey = (variables?: BankAccountsQueryVariables) => variables === undefined ? ['BankAccounts.infinite'] : ['BankAccounts.infinite', variables];


useBankAccountsQuery.fetcher = (client: GraphQLClient, variables?: BankAccountsQueryVariables, headers?: RequestInit['headers']) => fetcher<BankAccountsQuery, BankAccountsQueryVariables>(client, BankAccountsDocument, variables, headers);


useBankAccountsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: BankAccountsQueryVariables) =>
  (mutator: (cacheEntry: BankAccountsQuery) => BankAccountsQuery) => {
    const cacheKey = useBankAccountsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<BankAccountsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<BankAccountsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteBankAccountsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: BankAccountsQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<BankAccountsQuery>) => InfiniteData<BankAccountsQuery>) => {
    const cacheKey = useInfiniteBankAccountsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<BankAccountsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<BankAccountsQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }