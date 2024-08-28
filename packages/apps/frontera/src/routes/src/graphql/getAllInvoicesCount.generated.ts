// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

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
export type GetAllInvoicesCountQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetAllInvoicesCountQuery = { __typename?: 'Query', invoices: { __typename?: 'InvoicesPage', totalElements: any } };



export const GetAllInvoicesCountDocument = `
    query getAllInvoicesCount {
  invoices(
    pagination: {page: 0, limit: 0}
    where: {AND: [{filter: {property: "DRY_RUN", operation: EQ, value: false}}]}
  ) {
    totalElements
  }
}
    `;

export const useGetAllInvoicesCountQuery = <
      TData = GetAllInvoicesCountQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GetAllInvoicesCountQueryVariables,
      options?: Omit<UseQueryOptions<GetAllInvoicesCountQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetAllInvoicesCountQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetAllInvoicesCountQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['getAllInvoicesCount'] : ['getAllInvoicesCount', variables],
    queryFn: fetcher<GetAllInvoicesCountQuery, GetAllInvoicesCountQueryVariables>(client, GetAllInvoicesCountDocument, variables, headers),
    ...options
  }
    )};

useGetAllInvoicesCountQuery.document = GetAllInvoicesCountDocument;

useGetAllInvoicesCountQuery.getKey = (variables?: GetAllInvoicesCountQueryVariables) => variables === undefined ? ['getAllInvoicesCount'] : ['getAllInvoicesCount', variables];

export const useInfiniteGetAllInvoicesCountQuery = <
      TData = InfiniteData<GetAllInvoicesCountQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetAllInvoicesCountQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetAllInvoicesCountQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetAllInvoicesCountQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetAllInvoicesCountQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['getAllInvoicesCount.infinite'] : ['getAllInvoicesCount.infinite', variables],
      queryFn: (metaData) => fetcher<GetAllInvoicesCountQuery, GetAllInvoicesCountQueryVariables>(client, GetAllInvoicesCountDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetAllInvoicesCountQuery.getKey = (variables?: GetAllInvoicesCountQueryVariables) => variables === undefined ? ['getAllInvoicesCount.infinite'] : ['getAllInvoicesCount.infinite', variables];


useGetAllInvoicesCountQuery.fetcher = (client: GraphQLClient, variables?: GetAllInvoicesCountQueryVariables, headers?: RequestInit['headers']) => fetcher<GetAllInvoicesCountQuery, GetAllInvoicesCountQueryVariables>(client, GetAllInvoicesCountDocument, variables, headers);


useGetAllInvoicesCountQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetAllInvoicesCountQueryVariables) =>
  (mutator: (cacheEntry: GetAllInvoicesCountQuery) => GetAllInvoicesCountQuery) => {
    const cacheKey = useGetAllInvoicesCountQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetAllInvoicesCountQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetAllInvoicesCountQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetAllInvoicesCountQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetAllInvoicesCountQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetAllInvoicesCountQuery>) => InfiniteData<GetAllInvoicesCountQuery>) => {
    const cacheKey = useInfiniteGetAllInvoicesCountQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetAllInvoicesCountQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetAllInvoicesCountQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }