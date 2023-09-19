// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

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
