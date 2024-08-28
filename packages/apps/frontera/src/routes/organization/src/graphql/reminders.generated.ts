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
export type RemindersQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
}>;


export type RemindersQuery = { __typename?: 'Query', remindersForOrganization: Array<{ __typename?: 'Reminder', content?: string | null, dueDate?: any | null, dismissed?: boolean | null, metadata: { __typename?: 'Metadata', id: string, created: any, lastUpdated: any }, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string, name?: string | null } | null }> };



export const RemindersDocument = `
    query reminders($organizationId: ID!) {
  remindersForOrganization(organizationId: $organizationId) {
    metadata {
      id
      created
      lastUpdated
    }
    content
    owner {
      id
      firstName
      lastName
      name
    }
    dueDate
    dismissed
  }
}
    `;

export const useRemindersQuery = <
      TData = RemindersQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: RemindersQueryVariables,
      options?: Omit<UseQueryOptions<RemindersQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<RemindersQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<RemindersQuery, TError, TData>(
      {
    queryKey: ['reminders', variables],
    queryFn: fetcher<RemindersQuery, RemindersQueryVariables>(client, RemindersDocument, variables, headers),
    ...options
  }
    )};

useRemindersQuery.document = RemindersDocument;

useRemindersQuery.getKey = (variables: RemindersQueryVariables) => ['reminders', variables];

export const useInfiniteRemindersQuery = <
      TData = InfiniteData<RemindersQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: RemindersQueryVariables,
      options: Omit<UseInfiniteQueryOptions<RemindersQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<RemindersQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<RemindersQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['reminders.infinite', variables],
      queryFn: (metaData) => fetcher<RemindersQuery, RemindersQueryVariables>(client, RemindersDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteRemindersQuery.getKey = (variables: RemindersQueryVariables) => ['reminders.infinite', variables];


useRemindersQuery.fetcher = (client: GraphQLClient, variables: RemindersQueryVariables, headers?: RequestInit['headers']) => fetcher<RemindersQuery, RemindersQueryVariables>(client, RemindersDocument, variables, headers);


useRemindersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: RemindersQueryVariables) =>
  (mutator: (cacheEntry: RemindersQuery) => RemindersQuery) => {
    const cacheKey = useRemindersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<RemindersQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<RemindersQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteRemindersQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: RemindersQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<RemindersQuery>) => InfiniteData<RemindersQuery>) => {
    const cacheKey = useInfiniteRemindersQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<RemindersQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<RemindersQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }