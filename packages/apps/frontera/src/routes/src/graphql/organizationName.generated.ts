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
export type GetOrganizationNameQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type GetOrganizationNameQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', name: string } | null };



export const GetOrganizationNameDocument = `
    query GetOrganizationName($id: ID!) {
  organization(id: $id) {
    name
  }
}
    `;

export const useGetOrganizationNameQuery = <
      TData = GetOrganizationNameQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetOrganizationNameQueryVariables,
      options?: Omit<UseQueryOptions<GetOrganizationNameQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetOrganizationNameQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetOrganizationNameQuery, TError, TData>(
      {
    queryKey: ['GetOrganizationName', variables],
    queryFn: fetcher<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>(client, GetOrganizationNameDocument, variables, headers),
    ...options
  }
    )};

useGetOrganizationNameQuery.document = GetOrganizationNameDocument;

useGetOrganizationNameQuery.getKey = (variables: GetOrganizationNameQueryVariables) => ['GetOrganizationName', variables];

export const useInfiniteGetOrganizationNameQuery = <
      TData = InfiniteData<GetOrganizationNameQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetOrganizationNameQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetOrganizationNameQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetOrganizationNameQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetOrganizationNameQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['GetOrganizationName.infinite', variables],
      queryFn: (metaData) => fetcher<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>(client, GetOrganizationNameDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetOrganizationNameQuery.getKey = (variables: GetOrganizationNameQueryVariables) => ['GetOrganizationName.infinite', variables];


useGetOrganizationNameQuery.fetcher = (client: GraphQLClient, variables: GetOrganizationNameQueryVariables, headers?: RequestInit['headers']) => fetcher<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>(client, GetOrganizationNameDocument, variables, headers);


useGetOrganizationNameQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetOrganizationNameQueryVariables) =>
  (mutator: (cacheEntry: GetOrganizationNameQuery) => GetOrganizationNameQuery) => {
    const cacheKey = useGetOrganizationNameQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetOrganizationNameQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetOrganizationNameQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetOrganizationNameQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetOrganizationNameQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetOrganizationNameQuery>) => InfiniteData<GetOrganizationNameQuery>) => {
    const cacheKey = useInfiniteGetOrganizationNameQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetOrganizationNameQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetOrganizationNameQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }