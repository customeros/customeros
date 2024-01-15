// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

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
export type GetContractServicesQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type GetContractServicesQuery = {
  __typename?: 'Query';
  serviceLineItem: { __typename?: 'ServiceLineItem'; id: string; name: string };
};

export const GetContractServicesDocument = `
    query getContractServices($id: ID!) {
  serviceLineItem(id: $id) {
    id
    name
  }
}
    `;

export const useGetContractServicesQuery = <
  TData = GetContractServicesQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetContractServicesQueryVariables,
  options?: Omit<
    UseQueryOptions<GetContractServicesQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<
      GetContractServicesQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetContractServicesQuery, TError, TData>({
    queryKey: ['getContractServices', variables],
    queryFn: fetcher<
      GetContractServicesQuery,
      GetContractServicesQueryVariables
    >(client, GetContractServicesDocument, variables, headers),
    ...options,
  });
};

useGetContractServicesQuery.document = GetContractServicesDocument;

useGetContractServicesQuery.getKey = (
  variables: GetContractServicesQueryVariables,
) => ['getContractServices', variables];

export const useInfiniteGetContractServicesQuery = <
  TData = InfiniteData<GetContractServicesQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetContractServicesQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetContractServicesQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetContractServicesQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetContractServicesQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? [
          'getContractServices.infinite',
          variables,
        ],
        queryFn: (metaData) =>
          fetcher<GetContractServicesQuery, GetContractServicesQueryVariables>(
            client,
            GetContractServicesDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetContractServicesQuery.getKey = (
  variables: GetContractServicesQueryVariables,
) => ['getContractServices.infinite', variables];

useGetContractServicesQuery.fetcher = (
  client: GraphQLClient,
  variables: GetContractServicesQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetContractServicesQuery, GetContractServicesQueryVariables>(
    client,
    GetContractServicesDocument,
    variables,
    headers,
  );

useGetContractServicesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetContractServicesQueryVariables) =>
  (
    mutator: (cacheEntry: GetContractServicesQuery) => GetContractServicesQuery,
  ) => {
    const cacheKey = useGetContractServicesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetContractServicesQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetContractServicesQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetContractServicesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetContractServicesQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetContractServicesQuery>,
    ) => InfiniteData<GetContractServicesQuery>,
  ) => {
    const cacheKey = useInfiniteGetContractServicesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetContractServicesQuery>>(
        cacheKey,
      );
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetContractServicesQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
