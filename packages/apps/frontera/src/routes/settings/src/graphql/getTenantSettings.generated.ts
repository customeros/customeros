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
export type TenantSettingsQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type TenantSettingsQuery = { __typename?: 'Query', tenantSettings: { __typename?: 'TenantSettings', logoRepositoryFileId?: string | null, baseCurrency?: Types.Currency | null, billingEnabled: boolean } };



export const TenantSettingsDocument = `
    query TenantSettings {
  tenantSettings {
    logoRepositoryFileId
    baseCurrency
    billingEnabled
  }
}
    `;

export const useTenantSettingsQuery = <
      TData = TenantSettingsQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: TenantSettingsQueryVariables,
      options?: Omit<UseQueryOptions<TenantSettingsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<TenantSettingsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<TenantSettingsQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['TenantSettings'] : ['TenantSettings', variables],
    queryFn: fetcher<TenantSettingsQuery, TenantSettingsQueryVariables>(client, TenantSettingsDocument, variables, headers),
    ...options
  }
    )};

useTenantSettingsQuery.document = TenantSettingsDocument;

useTenantSettingsQuery.getKey = (variables?: TenantSettingsQueryVariables) => variables === undefined ? ['TenantSettings'] : ['TenantSettings', variables];

export const useInfiniteTenantSettingsQuery = <
      TData = InfiniteData<TenantSettingsQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: TenantSettingsQueryVariables,
      options: Omit<UseInfiniteQueryOptions<TenantSettingsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<TenantSettingsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<TenantSettingsQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['TenantSettings.infinite'] : ['TenantSettings.infinite', variables],
      queryFn: (metaData) => fetcher<TenantSettingsQuery, TenantSettingsQueryVariables>(client, TenantSettingsDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteTenantSettingsQuery.getKey = (variables?: TenantSettingsQueryVariables) => variables === undefined ? ['TenantSettings.infinite'] : ['TenantSettings.infinite', variables];


useTenantSettingsQuery.fetcher = (client: GraphQLClient, variables?: TenantSettingsQueryVariables, headers?: RequestInit['headers']) => fetcher<TenantSettingsQuery, TenantSettingsQueryVariables>(client, TenantSettingsDocument, variables, headers);


useTenantSettingsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantSettingsQueryVariables) =>
  (mutator: (cacheEntry: TenantSettingsQuery) => TenantSettingsQuery) => {
    const cacheKey = useTenantSettingsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TenantSettingsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TenantSettingsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteTenantSettingsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantSettingsQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<TenantSettingsQuery>) => InfiniteData<TenantSettingsQuery>) => {
    const cacheKey = useInfiniteTenantSettingsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TenantSettingsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TenantSettingsQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }