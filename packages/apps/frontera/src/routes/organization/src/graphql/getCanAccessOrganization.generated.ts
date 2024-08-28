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
export type GetCanAccessOrganizationQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type GetCanAccessOrganizationQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string } | null };



export const GetCanAccessOrganizationDocument = `
    query GetCanAccessOrganization($id: ID!) {
  organization(id: $id) {
    id
  }
}
    `;

export const useGetCanAccessOrganizationQuery = <
      TData = GetCanAccessOrganizationQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetCanAccessOrganizationQueryVariables,
      options?: Omit<UseQueryOptions<GetCanAccessOrganizationQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetCanAccessOrganizationQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetCanAccessOrganizationQuery, TError, TData>(
      {
    queryKey: ['GetCanAccessOrganization', variables],
    queryFn: fetcher<GetCanAccessOrganizationQuery, GetCanAccessOrganizationQueryVariables>(client, GetCanAccessOrganizationDocument, variables, headers),
    ...options
  }
    )};

useGetCanAccessOrganizationQuery.document = GetCanAccessOrganizationDocument;

useGetCanAccessOrganizationQuery.getKey = (variables: GetCanAccessOrganizationQueryVariables) => ['GetCanAccessOrganization', variables];

export const useInfiniteGetCanAccessOrganizationQuery = <
      TData = InfiniteData<GetCanAccessOrganizationQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetCanAccessOrganizationQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetCanAccessOrganizationQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetCanAccessOrganizationQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetCanAccessOrganizationQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['GetCanAccessOrganization.infinite', variables],
      queryFn: (metaData) => fetcher<GetCanAccessOrganizationQuery, GetCanAccessOrganizationQueryVariables>(client, GetCanAccessOrganizationDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetCanAccessOrganizationQuery.getKey = (variables: GetCanAccessOrganizationQueryVariables) => ['GetCanAccessOrganization.infinite', variables];


useGetCanAccessOrganizationQuery.fetcher = (client: GraphQLClient, variables: GetCanAccessOrganizationQueryVariables, headers?: RequestInit['headers']) => fetcher<GetCanAccessOrganizationQuery, GetCanAccessOrganizationQueryVariables>(client, GetCanAccessOrganizationDocument, variables, headers);


useGetCanAccessOrganizationQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetCanAccessOrganizationQueryVariables) =>
  (mutator: (cacheEntry: GetCanAccessOrganizationQuery) => GetCanAccessOrganizationQuery) => {
    const cacheKey = useGetCanAccessOrganizationQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetCanAccessOrganizationQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetCanAccessOrganizationQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetCanAccessOrganizationQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetCanAccessOrganizationQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetCanAccessOrganizationQuery>) => InfiniteData<GetCanAccessOrganizationQuery>) => {
    const cacheKey = useInfiniteGetCanAccessOrganizationQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetCanAccessOrganizationQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetCanAccessOrganizationQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }