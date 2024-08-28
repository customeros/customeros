// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, useInfiniteQuery, UseQueryOptions, UseInfiniteQueryOptions, InfiniteData } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type GetExternalSystemInstancesQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetExternalSystemInstancesQuery = { __typename?: 'Query', externalSystemInstances: Array<{ __typename?: 'ExternalSystemInstance', type: Types.ExternalSystemType, stripeDetails?: { __typename?: 'ExternalSystemStripeDetails', paymentMethodTypes: Array<string> } | null }> };



export const GetExternalSystemInstancesDocument = `
    query GetExternalSystemInstances {
  externalSystemInstances {
    type
    stripeDetails {
      paymentMethodTypes
    }
  }
}
    `;

export const useGetExternalSystemInstancesQuery = <
      TData = GetExternalSystemInstancesQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GetExternalSystemInstancesQueryVariables,
      options?: Omit<UseQueryOptions<GetExternalSystemInstancesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetExternalSystemInstancesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetExternalSystemInstancesQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['GetExternalSystemInstances'] : ['GetExternalSystemInstances', variables],
    queryFn: fetcher<GetExternalSystemInstancesQuery, GetExternalSystemInstancesQueryVariables>(client, GetExternalSystemInstancesDocument, variables, headers),
    ...options
  }
    )};

useGetExternalSystemInstancesQuery.document = GetExternalSystemInstancesDocument;

useGetExternalSystemInstancesQuery.getKey = (variables?: GetExternalSystemInstancesQueryVariables) => variables === undefined ? ['GetExternalSystemInstances'] : ['GetExternalSystemInstances', variables];

export const useInfiniteGetExternalSystemInstancesQuery = <
      TData = InfiniteData<GetExternalSystemInstancesQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetExternalSystemInstancesQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetExternalSystemInstancesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetExternalSystemInstancesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetExternalSystemInstancesQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['GetExternalSystemInstances.infinite'] : ['GetExternalSystemInstances.infinite', variables],
      queryFn: (metaData) => fetcher<GetExternalSystemInstancesQuery, GetExternalSystemInstancesQueryVariables>(client, GetExternalSystemInstancesDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetExternalSystemInstancesQuery.getKey = (variables?: GetExternalSystemInstancesQueryVariables) => variables === undefined ? ['GetExternalSystemInstances.infinite'] : ['GetExternalSystemInstances.infinite', variables];


useGetExternalSystemInstancesQuery.fetcher = (client: GraphQLClient, variables?: GetExternalSystemInstancesQueryVariables, headers?: RequestInit['headers']) => fetcher<GetExternalSystemInstancesQuery, GetExternalSystemInstancesQueryVariables>(client, GetExternalSystemInstancesDocument, variables, headers);


useGetExternalSystemInstancesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetExternalSystemInstancesQueryVariables) =>
  (mutator: (cacheEntry: GetExternalSystemInstancesQuery) => GetExternalSystemInstancesQuery) => {
    const cacheKey = useGetExternalSystemInstancesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetExternalSystemInstancesQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetExternalSystemInstancesQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetExternalSystemInstancesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetExternalSystemInstancesQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetExternalSystemInstancesQuery>) => InfiniteData<GetExternalSystemInstancesQuery>) => {
    const cacheKey = useInfiniteGetExternalSystemInstancesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetExternalSystemInstancesQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetExternalSystemInstancesQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }