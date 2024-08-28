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
export type BankAccountsCurrenciesQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type BankAccountsCurrenciesQuery = { __typename?: 'Query', bankAccounts: Array<{ __typename?: 'BankAccount', currency?: Types.Currency | null }> };



export const BankAccountsCurrenciesDocument = `
    query BankAccountsCurrencies {
  bankAccounts {
    currency
  }
}
    `;

export const useBankAccountsCurrenciesQuery = <
      TData = BankAccountsCurrenciesQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: BankAccountsCurrenciesQueryVariables,
      options?: Omit<UseQueryOptions<BankAccountsCurrenciesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<BankAccountsCurrenciesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<BankAccountsCurrenciesQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['BankAccountsCurrencies'] : ['BankAccountsCurrencies', variables],
    queryFn: fetcher<BankAccountsCurrenciesQuery, BankAccountsCurrenciesQueryVariables>(client, BankAccountsCurrenciesDocument, variables, headers),
    ...options
  }
    )};

useBankAccountsCurrenciesQuery.document = BankAccountsCurrenciesDocument;

useBankAccountsCurrenciesQuery.getKey = (variables?: BankAccountsCurrenciesQueryVariables) => variables === undefined ? ['BankAccountsCurrencies'] : ['BankAccountsCurrencies', variables];

export const useInfiniteBankAccountsCurrenciesQuery = <
      TData = InfiniteData<BankAccountsCurrenciesQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: BankAccountsCurrenciesQueryVariables,
      options: Omit<UseInfiniteQueryOptions<BankAccountsCurrenciesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<BankAccountsCurrenciesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<BankAccountsCurrenciesQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['BankAccountsCurrencies.infinite'] : ['BankAccountsCurrencies.infinite', variables],
      queryFn: (metaData) => fetcher<BankAccountsCurrenciesQuery, BankAccountsCurrenciesQueryVariables>(client, BankAccountsCurrenciesDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteBankAccountsCurrenciesQuery.getKey = (variables?: BankAccountsCurrenciesQueryVariables) => variables === undefined ? ['BankAccountsCurrencies.infinite'] : ['BankAccountsCurrencies.infinite', variables];


useBankAccountsCurrenciesQuery.fetcher = (client: GraphQLClient, variables?: BankAccountsCurrenciesQueryVariables, headers?: RequestInit['headers']) => fetcher<BankAccountsCurrenciesQuery, BankAccountsCurrenciesQueryVariables>(client, BankAccountsCurrenciesDocument, variables, headers);


useBankAccountsCurrenciesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: BankAccountsCurrenciesQueryVariables) =>
  (mutator: (cacheEntry: BankAccountsCurrenciesQuery) => BankAccountsCurrenciesQuery) => {
    const cacheKey = useBankAccountsCurrenciesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<BankAccountsCurrenciesQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<BankAccountsCurrenciesQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteBankAccountsCurrenciesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: BankAccountsCurrenciesQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<BankAccountsCurrenciesQuery>) => InfiniteData<BankAccountsCurrenciesQuery>) => {
    const cacheKey = useInfiniteBankAccountsCurrenciesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<BankAccountsCurrenciesQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<BankAccountsCurrenciesQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }