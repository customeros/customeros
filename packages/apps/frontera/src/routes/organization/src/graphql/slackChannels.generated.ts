// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

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
export type SlackChannelssQueryVariables = Types.Exact<{
  pagination?: Types.InputMaybe<Types.Pagination>;
}>;

export type SlackChannelssQuery = {
  __typename?: 'Query';
  slack_Channels: {
    __typename?: 'SlackChannelPage';
    totalElements: any;
    content: Array<{
      __typename?: 'SlackChannel';
      channelId: string;
      channelName: string;
      organization?: {
        __typename?: 'Organization';
        metadata: { __typename?: 'Metadata'; id: string };
      } | null;
    }>;
  };
};

export const SlackChannelssDocument = `
    query slackChannelss($pagination: Pagination) {
  slack_Channels(pagination: $pagination) {
    content {
      channelId
      channelName
      organization {
        metadata {
          id
        }
      }
    }
    totalElements
  }
}
    `;

export const useSlackChannelssQuery = <
  TData = SlackChannelssQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: SlackChannelssQueryVariables,
  options?: Omit<
    UseQueryOptions<SlackChannelssQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<SlackChannelssQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<SlackChannelssQuery, TError, TData>({
    queryKey:
      variables === undefined
        ? ['slackChannelss']
        : ['slackChannelss', variables],
    queryFn: fetcher<SlackChannelssQuery, SlackChannelssQueryVariables>(
      client,
      SlackChannelssDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useSlackChannelssQuery.document = SlackChannelssDocument;

useSlackChannelssQuery.getKey = (variables?: SlackChannelssQueryVariables) =>
  variables === undefined ? ['slackChannelss'] : ['slackChannelss', variables];

export const useInfiniteSlackChannelssQuery = <
  TData = InfiniteData<SlackChannelssQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: SlackChannelssQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<SlackChannelssQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      SlackChannelssQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<SlackChannelssQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey:
          optionsQueryKey ?? variables === undefined
            ? ['slackChannelss.infinite']
            : ['slackChannelss.infinite', variables],
        queryFn: (metaData) =>
          fetcher<SlackChannelssQuery, SlackChannelssQueryVariables>(
            client,
            SlackChannelssDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteSlackChannelssQuery.getKey = (
  variables?: SlackChannelssQueryVariables,
) =>
  variables === undefined
    ? ['slackChannelss.infinite']
    : ['slackChannelss.infinite', variables];

useSlackChannelssQuery.fetcher = (
  client: GraphQLClient,
  variables?: SlackChannelssQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<SlackChannelssQuery, SlackChannelssQueryVariables>(
    client,
    SlackChannelssDocument,
    variables,
    headers,
  );

useSlackChannelssQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: SlackChannelssQueryVariables) =>
  (mutator: (cacheEntry: SlackChannelssQuery) => SlackChannelssQuery) => {
    const cacheKey = useSlackChannelssQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<SlackChannelssQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<SlackChannelssQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteSlackChannelssQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: SlackChannelssQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<SlackChannelssQuery>,
    ) => InfiniteData<SlackChannelssQuery>,
  ) => {
    const cacheKey = useInfiniteSlackChannelssQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<SlackChannelssQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<SlackChannelssQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
