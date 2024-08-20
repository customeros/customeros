// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

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
export type GlobalCacheQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GlobalCacheQuery = {
  __typename?: 'Query';
  global_Cache: {
    __typename?: 'GlobalCache';
    cdnLogoUrl: string;
    mailboxes: Array<string>;
    isOwner: boolean;
    minARRForecastValue: number;
    maxARRForecastValue: number;
    contractsExist: boolean;
    user: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
      emails?: Array<{
        __typename?: 'Email';
        email?: string | null;
        rawEmail?: string | null;
        primary: boolean;
      }> | null;
    };
    inactiveEmailTokens: Array<{
      __typename?: 'GlobalCacheEmailToken';
      email: string;
      provider: string;
    }>;
    activeEmailTokens: Array<{
      __typename?: 'GlobalCacheEmailToken';
      email: string;
      provider: string;
    }>;
    gCliCache: Array<{
      __typename?: 'GCliItem';
      id: string;
      type: Types.GCliSearchResultType;
      display: string;
      data?: Array<{
        __typename?: 'GCliAttributeKeyValuePair';
        key: string;
        value: string;
        display?: string | null;
      }> | null;
    }>;
  };
};

export const GlobalCacheDocument = `
    query global_Cache {
  global_Cache {
    cdnLogoUrl
    user {
      id
      emails {
        email
        rawEmail
        primary
      }
      firstName
      lastName
    }
    inactiveEmailTokens {
      email
      provider
    }
    activeEmailTokens {
      email
      provider
    }
    mailboxes
    isOwner
    gCliCache {
      id
      type
      display
      data {
        key
        value
        display
      }
    }
    minARRForecastValue
    maxARRForecastValue
    contractsExist
  }
}
    `;

export const useGlobalCacheQuery = <TData = GlobalCacheQuery, TError = unknown>(
  client: GraphQLClient,
  variables?: GlobalCacheQueryVariables,
  options?: Omit<
    UseQueryOptions<GlobalCacheQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<GlobalCacheQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GlobalCacheQuery, TError, TData>({
    queryKey:
      variables === undefined ? ['global_Cache'] : ['global_Cache', variables],
    queryFn: fetcher<GlobalCacheQuery, GlobalCacheQueryVariables>(
      client,
      GlobalCacheDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGlobalCacheQuery.document = GlobalCacheDocument;

useGlobalCacheQuery.getKey = (variables?: GlobalCacheQueryVariables) =>
  variables === undefined ? ['global_Cache'] : ['global_Cache', variables];

export const useInfiniteGlobalCacheQuery = <
  TData = InfiniteData<GlobalCacheQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GlobalCacheQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GlobalCacheQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GlobalCacheQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GlobalCacheQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey:
          optionsQueryKey ?? variables === undefined
            ? ['global_Cache.infinite']
            : ['global_Cache.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GlobalCacheQuery, GlobalCacheQueryVariables>(
            client,
            GlobalCacheDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGlobalCacheQuery.getKey = (variables?: GlobalCacheQueryVariables) =>
  variables === undefined
    ? ['global_Cache.infinite']
    : ['global_Cache.infinite', variables];

useGlobalCacheQuery.fetcher = (
  client: GraphQLClient,
  variables?: GlobalCacheQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GlobalCacheQuery, GlobalCacheQueryVariables>(
    client,
    GlobalCacheDocument,
    variables,
    headers,
  );

useGlobalCacheQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GlobalCacheQueryVariables) =>
  (mutator: (cacheEntry: GlobalCacheQuery) => GlobalCacheQuery) => {
    const cacheKey = useGlobalCacheQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GlobalCacheQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GlobalCacheQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGlobalCacheQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GlobalCacheQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GlobalCacheQuery>,
    ) => InfiniteData<GlobalCacheQuery>,
  ) => {
    const cacheKey = useInfiniteGlobalCacheQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GlobalCacheQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GlobalCacheQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
