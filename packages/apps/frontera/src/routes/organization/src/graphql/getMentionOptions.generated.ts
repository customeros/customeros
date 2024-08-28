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
export type GetMentionOptionsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type GetMentionOptionsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, emails: Array<{ __typename?: 'Email', email?: string | null }> }> } } | null };



export const GetMentionOptionsDocument = `
    query getMentionOptions($id: ID!) {
  organization(id: $id) {
    contacts(pagination: {page: 0, limit: 100}) {
      content {
        id
        name
        firstName
        lastName
        emails {
          email
        }
      }
    }
  }
}
    `;

export const useGetMentionOptionsQuery = <
      TData = GetMentionOptionsQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetMentionOptionsQueryVariables,
      options?: Omit<UseQueryOptions<GetMentionOptionsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetMentionOptionsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetMentionOptionsQuery, TError, TData>(
      {
    queryKey: ['getMentionOptions', variables],
    queryFn: fetcher<GetMentionOptionsQuery, GetMentionOptionsQueryVariables>(client, GetMentionOptionsDocument, variables, headers),
    ...options
  }
    )};

useGetMentionOptionsQuery.document = GetMentionOptionsDocument;

useGetMentionOptionsQuery.getKey = (variables: GetMentionOptionsQueryVariables) => ['getMentionOptions', variables];

export const useInfiniteGetMentionOptionsQuery = <
      TData = InfiniteData<GetMentionOptionsQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetMentionOptionsQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetMentionOptionsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetMentionOptionsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetMentionOptionsQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['getMentionOptions.infinite', variables],
      queryFn: (metaData) => fetcher<GetMentionOptionsQuery, GetMentionOptionsQueryVariables>(client, GetMentionOptionsDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetMentionOptionsQuery.getKey = (variables: GetMentionOptionsQueryVariables) => ['getMentionOptions.infinite', variables];


useGetMentionOptionsQuery.fetcher = (client: GraphQLClient, variables: GetMentionOptionsQueryVariables, headers?: RequestInit['headers']) => fetcher<GetMentionOptionsQuery, GetMentionOptionsQueryVariables>(client, GetMentionOptionsDocument, variables, headers);


useGetMentionOptionsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetMentionOptionsQueryVariables) =>
  (mutator: (cacheEntry: GetMentionOptionsQuery) => GetMentionOptionsQuery) => {
    const cacheKey = useGetMentionOptionsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetMentionOptionsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetMentionOptionsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetMentionOptionsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetMentionOptionsQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetMentionOptionsQuery>) => InfiniteData<GetMentionOptionsQuery>) => {
    const cacheKey = useInfiniteGetMentionOptionsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetMentionOptionsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetMentionOptionsQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }