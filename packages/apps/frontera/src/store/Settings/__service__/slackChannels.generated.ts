// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../routes/src/types/__generated__/graphql.types';

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
export type SlackChannelsQueryVariables = Types.Exact<{
  pagination?: Types.InputMaybe<Types.Pagination>;
}>;

export type SlackChannelsQuery = {
  __typename?: 'Query';
  slack_Channels: {
    __typename?: 'SlackChannelPage';
    totalElements: any;
    content: Array<{
      __typename?: 'SlackChannel';
      channelId: string;
      channelName: string;
      metadata: {
        __typename?: 'Metadata';
        id: string;
        appSource: string;
        source: Types.DataSource;
        sourceOfTruth: Types.DataSource;
      };
      organization?: {
        __typename?: 'Organization';
        metadata: { __typename?: 'Metadata'; id: string };
      } | null;
    }>;
  };
};

export const SlackChannelsDocument = `
    query slackChannels($pagination: Pagination) {
  slack_Channels(pagination: $pagination) {
    content {
      channelId
      channelName
      metadata {
        id
        appSource
        source
        sourceOfTruth
      }
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

export const useSlackChannelsQuery = <
  TData = SlackChannelsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: SlackChannelsQueryVariables,
  options?: Omit<
    UseQueryOptions<SlackChannelsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<SlackChannelsQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<SlackChannelsQuery, TError, TData>({
    queryKey:
      variables === undefined
        ? ['slackChannels']
        : ['slackChannels', variables],
    queryFn: fetcher<SlackChannelsQuery, SlackChannelsQueryVariables>(
      client,
      SlackChannelsDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useSlackChannelsQuery.document = SlackChannelsDocument;

useSlackChannelsQuery.getKey = (variables?: SlackChannelsQueryVariables) =>
  variables === undefined ? ['slackChannels'] : ['slackChannels', variables];

export const useInfiniteSlackChannelsQuery = <
  TData = InfiniteData<SlackChannelsQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: SlackChannelsQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<SlackChannelsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      SlackChannelsQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<SlackChannelsQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey:
          optionsQueryKey ?? variables === undefined
            ? ['slackChannels.infinite']
            : ['slackChannels.infinite', variables],
        queryFn: (metaData) =>
          fetcher<SlackChannelsQuery, SlackChannelsQueryVariables>(
            client,
            SlackChannelsDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteSlackChannelsQuery.getKey = (
  variables?: SlackChannelsQueryVariables,
) =>
  variables === undefined
    ? ['slackChannels.infinite']
    : ['slackChannels.infinite', variables];

useSlackChannelsQuery.fetcher = (
  client: GraphQLClient,
  variables?: SlackChannelsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<SlackChannelsQuery, SlackChannelsQueryVariables>(
    client,
    SlackChannelsDocument,
    variables,
    headers,
  );

useSlackChannelsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: SlackChannelsQueryVariables) =>
  (mutator: (cacheEntry: SlackChannelsQuery) => SlackChannelsQuery) => {
    const cacheKey = useSlackChannelsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<SlackChannelsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<SlackChannelsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteSlackChannelsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: SlackChannelsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<SlackChannelsQuery>,
    ) => InfiniteData<SlackChannelsQuery>,
  ) => {
    const cacheKey = useInfiniteSlackChannelsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<SlackChannelsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<SlackChannelsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
