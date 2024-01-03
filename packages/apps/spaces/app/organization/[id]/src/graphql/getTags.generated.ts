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
export type GetTagsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GetTagsQuery = {
  __typename?: 'Query';
  tags: Array<{ __typename?: 'Tag'; value: string; label: string }>;
};

export const GetTagsDocument = `
    query getTags {
  tags {
    value: id
    label: name
  }
}
    `;
export const useGetTagsQuery = <TData = GetTagsQuery, TError = unknown>(
  client: GraphQLClient,
  variables?: GetTagsQueryVariables,
  options?: UseQueryOptions<GetTagsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetTagsQuery, TError, TData>(
    variables === undefined ? ['getTags'] : ['getTags', variables],
    fetcher<GetTagsQuery, GetTagsQueryVariables>(
      client,
      GetTagsDocument,
      variables,
      headers,
    ),
    options,
  );
useGetTagsQuery.document = GetTagsDocument;

useGetTagsQuery.getKey = (variables?: GetTagsQueryVariables) =>
  variables === undefined ? ['getTags'] : ['getTags', variables];
export const useInfiniteGetTagsQuery = <TData = GetTagsQuery, TError = unknown>(
  pageParamKey: keyof GetTagsQueryVariables,
  client: GraphQLClient,
  variables?: GetTagsQueryVariables,
  options?: UseInfiniteQueryOptions<GetTagsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetTagsQuery, TError, TData>(
    variables === undefined
      ? ['getTags.infinite']
      : ['getTags.infinite', variables],
    (metaData) =>
      fetcher<GetTagsQuery, GetTagsQueryVariables>(
        client,
        GetTagsDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetTagsQuery.getKey = (variables?: GetTagsQueryVariables) =>
  variables === undefined
    ? ['getTags.infinite']
    : ['getTags.infinite', variables];
useGetTagsQuery.fetcher = (
  client: GraphQLClient,
  variables?: GetTagsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetTagsQuery, GetTagsQueryVariables>(
    client,
    GetTagsDocument,
    variables,
    headers,
  );

useGetTagsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetTagsQueryVariables) =>
  (mutator: (cacheEntry: GetTagsQuery) => GetTagsQuery) => {
    const cacheKey = useGetTagsQuery.getKey(variables);
    const previousEntries = queryClient.getQueryData<GetTagsQuery>(cacheKey);
    if (previousEntry) {
      queryClient.setQueryData<GetTagsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetTagsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetTagsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetTagsQuery>,
    ) => InfiniteData<GetTagsQuery>,
  ) => {
    const cacheKey = useInfiniteGetTagsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetTagsQuery>>(cacheKey);
    if (previousEntry) {
      queryClient.setQueryData<InfiniteData<GetTagsQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  };
