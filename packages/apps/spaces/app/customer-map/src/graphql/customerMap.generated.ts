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
export type CustomerMapQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type CustomerMapQuery = { __typename?: 'Query', dashboard_CustomerMap?: Array<{ __typename?: 'DashboardCustomerMap', state: Types.DashboardCustomerMapState, arr: number, contractSignedDate: any, organization: { __typename?: 'Organization', id: string, name: string } }> | null };



export const CustomerMapDocument = `
    query CustomerMap {
  dashboard_CustomerMap {
    organization {
      id
      name
    }
    state
    arr
    contractSignedDate
  }
}
    `;

export const useCustomerMapQuery = <
      TData = CustomerMapQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: CustomerMapQueryVariables,
      options?: Omit<UseQueryOptions<CustomerMapQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<CustomerMapQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<CustomerMapQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['CustomerMap'] : ['CustomerMap', variables],
    queryFn: fetcher<CustomerMapQuery, CustomerMapQueryVariables>(client, CustomerMapDocument, variables, headers),
    ...options
  }
    )};

useCustomerMapQuery.document = CustomerMapDocument;

useCustomerMapQuery.getKey = (variables?: CustomerMapQueryVariables) => variables === undefined ? ['CustomerMap'] : ['CustomerMap', variables];

export const useInfiniteCustomerMapQuery = <
      TData = InfiniteData<CustomerMapQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: CustomerMapQueryVariables,
      options: Omit<UseInfiniteQueryOptions<CustomerMapQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<CustomerMapQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<CustomerMapQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['CustomerMap.infinite'] : ['CustomerMap.infinite', variables],
      queryFn: (metaData) => fetcher<CustomerMapQuery, CustomerMapQueryVariables>(client, CustomerMapDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteCustomerMapQuery.getKey = (variables?: CustomerMapQueryVariables) => variables === undefined ? ['CustomerMap.infinite'] : ['CustomerMap.infinite', variables];


useCustomerMapQuery.fetcher = (client: GraphQLClient, variables?: CustomerMapQueryVariables, headers?: RequestInit['headers']) => fetcher<CustomerMapQuery, CustomerMapQueryVariables>(client, CustomerMapDocument, variables, headers);


useCustomerMapQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: CustomerMapQueryVariables) =>
  (mutator: (cacheEntry: CustomerMapQuery) => CustomerMapQuery) => {
    const cacheKey = useCustomerMapQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<CustomerMapQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<CustomerMapQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteCustomerMapQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: CustomerMapQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<CustomerMapQuery>) => InfiniteData<CustomerMapQuery>) => {
    const cacheKey = useInfiniteCustomerMapQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<CustomerMapQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<CustomerMapQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }