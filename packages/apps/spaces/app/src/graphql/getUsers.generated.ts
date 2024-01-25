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
export type GetUsersQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
}>;


export type GetUsersQuery = { __typename?: 'Query', users: { __typename?: 'UserPage', totalElements: any, content: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string, name?: string | null }> } };



export const GetUsersDocument = `
    query getUsers($pagination: Pagination!, $where: Filter) {
  users(pagination: $pagination, where: $where) {
    content {
      id
      firstName
      lastName
      name
    }
    totalElements
  }
}
    `;

export const useGetUsersQuery = <
      TData = GetUsersQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetUsersQueryVariables,
      options?: Omit<UseQueryOptions<GetUsersQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetUsersQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetUsersQuery, TError, TData>(
      {
    queryKey: ['getUsers', variables],
    queryFn: fetcher<GetUsersQuery, GetUsersQueryVariables>(client, GetUsersDocument, variables, headers),
    ...options
  }
    )};

useGetUsersQuery.document = GetUsersDocument;

useGetUsersQuery.getKey = (variables: GetUsersQueryVariables) => ['getUsers', variables];

export const useInfiniteGetUsersQuery = <
      TData = InfiniteData<GetUsersQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetUsersQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetUsersQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetUsersQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetUsersQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['getUsers.infinite', variables],
      queryFn: (metaData) => fetcher<GetUsersQuery, GetUsersQueryVariables>(client, GetUsersDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetUsersQuery.getKey = (variables: GetUsersQueryVariables) => ['getUsers.infinite', variables];


useGetUsersQuery.fetcher = (client: GraphQLClient, variables: GetUsersQueryVariables, headers?: RequestInit['headers']) => fetcher<GetUsersQuery, GetUsersQueryVariables>(client, GetUsersDocument, variables, headers);


useGetUsersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetUsersQueryVariables) =>
  (mutator: (cacheEntry: GetUsersQuery) => GetUsersQuery) => {
    const cacheKey = useGetUsersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetUsersQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetUsersQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetUsersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetUsersQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetUsersQuery>) => InfiniteData<GetUsersQuery>) => {
    const cacheKey = useInfiniteGetUsersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetUsersQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetUsersQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }