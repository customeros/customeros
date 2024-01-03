// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

import type { InfiniteData } from '@tanstack/react-query';
import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
} from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(
  client: GraphQLClient,
  query: string,
  variables?: TVariables,
  requestHeaders?: RequestInit['headers'],
) {
  return async (): Promise<TData> =>
    client.request({
      document: query,
      variables,
      requestHeaders,
    });
}
export type GetMentionOptionsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type GetMentionOptionsQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    contacts: {
      __typename?: 'ContactsPage';
      content: Array<{
        __typename?: 'Contact';
        id: string;
        name?: string | null;
        firstName?: string | null;
        lastName?: string | null;
        emails: Array<{ __typename?: 'Email'; email?: string | null }>;
      }>;
    };
  } | null;
};

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
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetMentionOptionsQueryVariables,
  options?: UseQueryOptions<GetMentionOptionsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetMentionOptionsQuery, TError, TData>(
    ['getMentionOptions', variables],
    fetcher<GetMentionOptionsQuery, GetMentionOptionsQueryVariables>(
      client,
      GetMentionOptionsDocument,
      variables,
      headers,
    ),
    options,
  );
useGetMentionOptionsQuery.document = GetMentionOptionsDocument;

useGetMentionOptionsQuery.getKey = (
  variables: GetMentionOptionsQueryVariables,
) => ['getMentionOptions', variables];
export const useInfiniteGetMentionOptionsQuery = <
  TData = GetMentionOptionsQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetMentionOptionsQueryVariables,
  client: GraphQLClient,
  variables: GetMentionOptionsQueryVariables,
  options?: UseInfiniteQueryOptions<GetMentionOptionsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetMentionOptionsQuery, TError, TData>(
    ['getMentionOptions.infinite', variables],
    (metaData) =>
      fetcher<GetMentionOptionsQuery, GetMentionOptionsQueryVariables>(
        client,
        GetMentionOptionsDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetMentionOptionsQuery.getKey = (
  variables: GetMentionOptionsQueryVariables,
) => ['getMentionOptions.infinite', variables];
useGetMentionOptionsQuery.fetcher = (
  client: GraphQLClient,
  variables: GetMentionOptionsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetMentionOptionsQuery, GetMentionOptionsQueryVariables>(
    client,
    GetMentionOptionsDocument,
    variables,
    headers,
  );

useGetMentionOptionsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetMentionOptionsQueryVariables) =>
  (mutator: (cacheEntry: GetMentionOptionsQuery) => GetMentionOptionsQuery) => {
    const cacheKey = useGetMentionOptionsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetMentionOptionsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetMentionOptionsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetMentionOptionsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetMentionOptionsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetMentionOptionsQuery>,
    ) => InfiniteData<GetMentionOptionsQuery>,
  ) => {
    const cacheKey = useInfiniteGetMentionOptionsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetMentionOptionsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetMentionOptionsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
