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
export type GetBillableInfoQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetBillableInfoQuery = {
  __typename?: 'Query';
  billableInfo: {
    __typename?: 'TenantBillableInfo';
    whitelistedOrganizations: any;
    whitelistedContacts: any;
    greylistedOrganizations: any;
    greylistedContacts: any;
  };
};

export const GetBillableInfoDocument = `
    query getBillableInfo {
  billableInfo {
    whitelistedOrganizations
    whitelistedContacts
    greylistedOrganizations
    greylistedContacts
  }
}
    `;

export const useGetBillableInfoQuery = <
  TData = GetBillableInfoQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: GetBillableInfoQueryVariables,
  options?: Omit<
    UseQueryOptions<GetBillableInfoQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<GetBillableInfoQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetBillableInfoQuery, TError, TData>({
    queryKey:
      variables === undefined
        ? ['getBillableInfo']
        : ['getBillableInfo', variables],
    queryFn: fetcher<GetBillableInfoQuery, GetBillableInfoQueryVariables>(
      client,
      GetBillableInfoDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGetBillableInfoQuery.document = GetBillableInfoDocument;

useGetBillableInfoQuery.getKey = (variables?: GetBillableInfoQueryVariables) =>
  variables === undefined
    ? ['getBillableInfo']
    : ['getBillableInfo', variables];

export const useInfiniteGetBillableInfoQuery = <
  TData = InfiniteData<GetBillableInfoQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetBillableInfoQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetBillableInfoQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetBillableInfoQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetBillableInfoQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey:
          optionsQueryKey ?? variables === undefined
            ? ['getBillableInfo.infinite']
            : ['getBillableInfo.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GetBillableInfoQuery, GetBillableInfoQueryVariables>(
            client,
            GetBillableInfoDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetBillableInfoQuery.getKey = (
  variables?: GetBillableInfoQueryVariables,
) =>
  variables === undefined
    ? ['getBillableInfo.infinite']
    : ['getBillableInfo.infinite', variables];

useGetBillableInfoQuery.fetcher = (
  client: GraphQLClient,
  variables?: GetBillableInfoQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetBillableInfoQuery, GetBillableInfoQueryVariables>(
    client,
    GetBillableInfoDocument,
    variables,
    headers,
  );

useGetBillableInfoQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetBillableInfoQueryVariables) =>
  (mutator: (cacheEntry: GetBillableInfoQuery) => GetBillableInfoQuery) => {
    const cacheKey = useGetBillableInfoQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetBillableInfoQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetBillableInfoQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetBillableInfoQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetBillableInfoQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetBillableInfoQuery>,
    ) => InfiniteData<GetBillableInfoQuery>,
  ) => {
    const cacheKey = useInfiniteGetBillableInfoQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetBillableInfoQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetBillableInfoQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
