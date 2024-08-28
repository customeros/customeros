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
export type BaseCurrencyQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type BaseCurrencyQuery = { __typename?: 'Query', tenantSettings: { __typename?: 'TenantSettings', baseCurrency?: Types.Currency | null } };



export const BaseCurrencyDocument = `
    query BaseCurrency {
  tenantSettings {
    baseCurrency
  }
}
    `;

export const useBaseCurrencyQuery = <
      TData = BaseCurrencyQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: BaseCurrencyQueryVariables,
      options?: Omit<UseQueryOptions<BaseCurrencyQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<BaseCurrencyQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<BaseCurrencyQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['BaseCurrency'] : ['BaseCurrency', variables],
    queryFn: fetcher<BaseCurrencyQuery, BaseCurrencyQueryVariables>(client, BaseCurrencyDocument, variables, headers),
    ...options
  }
    )};

useBaseCurrencyQuery.document = BaseCurrencyDocument;

useBaseCurrencyQuery.getKey = (variables?: BaseCurrencyQueryVariables) => variables === undefined ? ['BaseCurrency'] : ['BaseCurrency', variables];

export const useInfiniteBaseCurrencyQuery = <
      TData = InfiniteData<BaseCurrencyQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: BaseCurrencyQueryVariables,
      options: Omit<UseInfiniteQueryOptions<BaseCurrencyQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<BaseCurrencyQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<BaseCurrencyQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['BaseCurrency.infinite'] : ['BaseCurrency.infinite', variables],
      queryFn: (metaData) => fetcher<BaseCurrencyQuery, BaseCurrencyQueryVariables>(client, BaseCurrencyDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteBaseCurrencyQuery.getKey = (variables?: BaseCurrencyQueryVariables) => variables === undefined ? ['BaseCurrency.infinite'] : ['BaseCurrency.infinite', variables];


useBaseCurrencyQuery.fetcher = (client: GraphQLClient, variables?: BaseCurrencyQueryVariables, headers?: RequestInit['headers']) => fetcher<BaseCurrencyQuery, BaseCurrencyQueryVariables>(client, BaseCurrencyDocument, variables, headers);


useBaseCurrencyQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: BaseCurrencyQueryVariables) =>
  (mutator: (cacheEntry: BaseCurrencyQuery) => BaseCurrencyQuery) => {
    const cacheKey = useBaseCurrencyQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<BaseCurrencyQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<BaseCurrencyQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteBaseCurrencyQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: BaseCurrencyQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<BaseCurrencyQuery>) => InfiniteData<BaseCurrencyQuery>) => {
    const cacheKey = useInfiniteBaseCurrencyQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<BaseCurrencyQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<BaseCurrencyQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }