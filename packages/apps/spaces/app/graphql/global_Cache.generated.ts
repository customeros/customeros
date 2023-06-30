// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';

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
    isOwner: boolean;
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
  }
}
    `;
export const useGlobalCacheQuery = <TData = GlobalCacheQuery, TError = unknown>(
  client: GraphQLClient,
  variables?: GlobalCacheQueryVariables,
  options?: UseQueryOptions<GlobalCacheQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GlobalCacheQuery, TError, TData>(
    variables === undefined ? ['global_Cache'] : ['global_Cache', variables],
    fetcher<GlobalCacheQuery, GlobalCacheQueryVariables>(
      client,
      GlobalCacheDocument,
      variables,
      headers,
    ),
    options,
  );
