// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
  InfiniteData,
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
export type GetLogEntryQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type GetLogEntryQuery = {
  __typename?: 'Query';
  logEntry: {
    __typename?: 'LogEntry';
    id: string;
    content?: string | null;
    contentType?: string | null;
    createdAt: any;
    updatedAt: any;
    tags: Array<{ __typename?: 'Tag'; id: string; name: string }>;
    createdBy?: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
      name?: string | null;
    } | null;
  };
};

export const GetLogEntryDocument = `
    query getLogEntry($id: ID!) {
  logEntry(id: $id) {
    id
    content
    contentType
    createdAt
    updatedAt
    tags {
      id
      name
    }
    createdBy {
      id
      firstName
      lastName
      name
    }
  }
}
    `;

export const useGetLogEntryQuery = <TData = GetLogEntryQuery, TError = unknown>(
  client: GraphQLClient,
  variables: GetLogEntryQueryVariables,
  options?: Omit<
    UseQueryOptions<GetLogEntryQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<GetLogEntryQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetLogEntryQuery, TError, TData>({
    queryKey: ['getLogEntry', variables],
    queryFn: fetcher<GetLogEntryQuery, GetLogEntryQueryVariables>(
      client,
      GetLogEntryDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGetLogEntryQuery.document = GetLogEntryDocument;

useGetLogEntryQuery.getKey = (variables: GetLogEntryQueryVariables) => [
  'getLogEntry',
  variables,
];

export const useInfiniteGetLogEntryQuery = <
  TData = InfiniteData<GetLogEntryQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetLogEntryQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetLogEntryQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetLogEntryQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetLogEntryQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? ['getLogEntry.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GetLogEntryQuery, GetLogEntryQueryVariables>(
            client,
            GetLogEntryDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetLogEntryQuery.getKey = (variables: GetLogEntryQueryVariables) => [
  'getLogEntry.infinite',
  variables,
];

useGetLogEntryQuery.fetcher = (
  client: GraphQLClient,
  variables: GetLogEntryQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetLogEntryQuery, GetLogEntryQueryVariables>(
    client,
    GetLogEntryDocument,
    variables,
    headers,
  );

useGetLogEntryQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetLogEntryQueryVariables) =>
  (mutator: (cacheEntry: GetLogEntryQuery) => GetLogEntryQuery) => {
    const cacheKey = useGetLogEntryQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetLogEntryQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetLogEntryQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetLogEntryQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetLogEntryQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetLogEntryQuery>,
    ) => InfiniteData<GetLogEntryQuery>,
  ) => {
    const cacheKey = useInfiniteGetLogEntryQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetLogEntryQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetLogEntryQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
