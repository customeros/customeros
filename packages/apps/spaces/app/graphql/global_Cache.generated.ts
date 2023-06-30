import * as Types from '../types/__generated__/graphql.types';

import { useQuery, UseQueryOptions } from '@tanstack/react-query';

function fetcher<TData, TVariables>(
  endpoint: string,
  requestInit: RequestInit,
  query: string,
  variables?: TVariables,
) {
  return async (): Promise<TData> => {
    const res = await fetch(endpoint, {
      method: 'POST',
      ...requestInit,
      body: JSON.stringify({ query, variables }),
    });

    const json = await res.json();

    if (json.errors) {
      const { message } = json.errors[0];

      throw new Error(message);
    }

    return json.data;
  };
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
  dataSource: { endpoint: string; fetchParams?: RequestInit },
  variables?: GlobalCacheQueryVariables,
  options?: UseQueryOptions<GlobalCacheQuery, TError, TData>,
) =>
  useQuery<GlobalCacheQuery, TError, TData>(
    variables === undefined ? ['global_Cache'] : ['global_Cache', variables],
    fetcher<GlobalCacheQuery, GlobalCacheQueryVariables>(
      dataSource.endpoint,
      dataSource.fetchParams || {},
      GlobalCacheDocument,
      variables,
    ),
    options,
  );
