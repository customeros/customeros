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
export type TableViewDefsQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type TableViewDefsQuery = { __typename?: 'Query', tableViewDefs: Array<{ __typename?: 'TableViewDef', id: string }> };



export const TableViewDefsDocument = `
    query tableViewDefs {
  tableViewDefs {
    id
  }
}
    `;

export const useTableViewDefsQuery = <
      TData = TableViewDefsQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: TableViewDefsQueryVariables,
      options?: Omit<UseQueryOptions<TableViewDefsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<TableViewDefsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<TableViewDefsQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['tableViewDefs'] : ['tableViewDefs', variables],
    queryFn: fetcher<TableViewDefsQuery, TableViewDefsQueryVariables>(client, TableViewDefsDocument, variables, headers),
    ...options
  }
    )};

useTableViewDefsQuery.document = TableViewDefsDocument;

useTableViewDefsQuery.getKey = (variables?: TableViewDefsQueryVariables) => variables === undefined ? ['tableViewDefs'] : ['tableViewDefs', variables];

export const useInfiniteTableViewDefsQuery = <
      TData = InfiniteData<TableViewDefsQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: TableViewDefsQueryVariables,
      options: Omit<UseInfiniteQueryOptions<TableViewDefsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<TableViewDefsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<TableViewDefsQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['tableViewDefs.infinite'] : ['tableViewDefs.infinite', variables],
      queryFn: (metaData) => fetcher<TableViewDefsQuery, TableViewDefsQueryVariables>(client, TableViewDefsDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteTableViewDefsQuery.getKey = (variables?: TableViewDefsQueryVariables) => variables === undefined ? ['tableViewDefs.infinite'] : ['tableViewDefs.infinite', variables];


useTableViewDefsQuery.fetcher = (client: GraphQLClient, variables?: TableViewDefsQueryVariables, headers?: RequestInit['headers']) => fetcher<TableViewDefsQuery, TableViewDefsQueryVariables>(client, TableViewDefsDocument, variables, headers);


useTableViewDefsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TableViewDefsQueryVariables) =>
  (mutator: (cacheEntry: TableViewDefsQuery) => TableViewDefsQuery) => {
    const cacheKey = useTableViewDefsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TableViewDefsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TableViewDefsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteTableViewDefsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TableViewDefsQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<TableViewDefsQuery>) => InfiniteData<TableViewDefsQuery>) => {
    const cacheKey = useInfiniteTableViewDefsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TableViewDefsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TableViewDefsQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }