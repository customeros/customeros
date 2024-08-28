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
export type GetContactsEmailListQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Array<Types.SortBy> | Types.SortBy>;
}>;


export type GetContactsEmailListQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, emails: Array<{ __typename?: 'Email', id: string, email?: string | null }> }> } } | null };



export const GetContactsEmailListDocument = `
    query GetContactsEmailList($id: ID!, $pagination: Pagination!, $where: Filter, $sort: [SortBy!]) {
  organization(id: $id) {
    id
    contacts(pagination: $pagination, where: $where, sort: $sort) {
      content {
        id
        firstName
        lastName
        emails {
          id
          email
        }
      }
    }
  }
}
    `;

export const useGetContactsEmailListQuery = <
      TData = GetContactsEmailListQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetContactsEmailListQueryVariables,
      options?: Omit<UseQueryOptions<GetContactsEmailListQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetContactsEmailListQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetContactsEmailListQuery, TError, TData>(
      {
    queryKey: ['GetContactsEmailList', variables],
    queryFn: fetcher<GetContactsEmailListQuery, GetContactsEmailListQueryVariables>(client, GetContactsEmailListDocument, variables, headers),
    ...options
  }
    )};

useGetContactsEmailListQuery.document = GetContactsEmailListDocument;

useGetContactsEmailListQuery.getKey = (variables: GetContactsEmailListQueryVariables) => ['GetContactsEmailList', variables];

export const useInfiniteGetContactsEmailListQuery = <
      TData = InfiniteData<GetContactsEmailListQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetContactsEmailListQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetContactsEmailListQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetContactsEmailListQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetContactsEmailListQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['GetContactsEmailList.infinite', variables],
      queryFn: (metaData) => fetcher<GetContactsEmailListQuery, GetContactsEmailListQueryVariables>(client, GetContactsEmailListDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetContactsEmailListQuery.getKey = (variables: GetContactsEmailListQueryVariables) => ['GetContactsEmailList.infinite', variables];


useGetContactsEmailListQuery.fetcher = (client: GraphQLClient, variables: GetContactsEmailListQueryVariables, headers?: RequestInit['headers']) => fetcher<GetContactsEmailListQuery, GetContactsEmailListQueryVariables>(client, GetContactsEmailListDocument, variables, headers);


useGetContactsEmailListQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetContactsEmailListQueryVariables) =>
  (mutator: (cacheEntry: GetContactsEmailListQuery) => GetContactsEmailListQuery) => {
    const cacheKey = useGetContactsEmailListQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetContactsEmailListQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetContactsEmailListQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetContactsEmailListQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetContactsEmailListQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetContactsEmailListQuery>) => InfiniteData<GetContactsEmailListQuery>) => {
    const cacheKey = useInfiniteGetContactsEmailListQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetContactsEmailListQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetContactsEmailListQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }