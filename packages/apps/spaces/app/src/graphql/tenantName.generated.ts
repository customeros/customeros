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
export type TenantNameQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type TenantNameQuery = { __typename?: 'Query', tenant: string };



export const TenantNameDocument = `
    query TenantName {
  tenant
}
    `;

export const useTenantNameQuery = <
      TData = TenantNameQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: TenantNameQueryVariables,
      options?: Omit<UseQueryOptions<TenantNameQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<TenantNameQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<TenantNameQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['TenantName'] : ['TenantName', variables],
    queryFn: fetcher<TenantNameQuery, TenantNameQueryVariables>(client, TenantNameDocument, variables, headers),
    ...options
  }
    )};

useTenantNameQuery.document = TenantNameDocument;

useTenantNameQuery.getKey = (variables?: TenantNameQueryVariables) => variables === undefined ? ['TenantName'] : ['TenantName', variables];

export const useInfiniteTenantNameQuery = <
      TData = InfiniteData<TenantNameQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: TenantNameQueryVariables,
      options: Omit<UseInfiniteQueryOptions<TenantNameQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<TenantNameQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<TenantNameQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['TenantName.infinite'] : ['TenantName.infinite', variables],
      queryFn: (metaData) => fetcher<TenantNameQuery, TenantNameQueryVariables>(client, TenantNameDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteTenantNameQuery.getKey = (variables?: TenantNameQueryVariables) => variables === undefined ? ['TenantName.infinite'] : ['TenantName.infinite', variables];


useTenantNameQuery.fetcher = (client: GraphQLClient, variables?: TenantNameQueryVariables, headers?: RequestInit['headers']) => fetcher<TenantNameQuery, TenantNameQueryVariables>(client, TenantNameDocument, variables, headers);


useTenantNameQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantNameQueryVariables) =>
  (mutator: (cacheEntry: TenantNameQuery) => TenantNameQuery) => {
    const cacheKey = useTenantNameQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TenantNameQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TenantNameQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteTenantNameQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantNameQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<TenantNameQuery>) => InfiniteData<TenantNameQuery>) => {
    const cacheKey = useInfiniteTenantNameQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TenantNameQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TenantNameQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }