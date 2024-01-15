// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import type { InfiniteData } from '@tanstack/react-query';
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
export type TableViewDefsQueryVariables = Types.Exact<{
  pagination?: Types.InputMaybe<Types.Pagination>;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;

export type TableViewDefsQuery = {
  __typename?: 'Query';
  tableViewDefs: {
    __typename?: 'TableViewDefPage';
    content: Array<{
      __typename?: 'TableViewDef';
      id: string;
      name: string;
      order?: number | null;
      filters?: string | null;
      sorting?: string | null;
      type?: {
        __typename?: 'ViewType';
        id: string;
        name?: string | null;
      } | null;
      columns?: Array<{
        __typename?: 'ColumnDef';
        id: string;
        isFilterable?: boolean | null;
        isSortable?: boolean | null;
        isDefaultSort?: boolean | null;
        isVisible?: boolean | null;
        columnType?: {
          __typename?: 'ColumnType';
          id: string;
          name?: string | null;
        } | null;
      } | null> | null;
    }>;
  };
};

export const TableViewDefsDocument = `
    query tableViewDefs($pagination: Pagination, $where: Filter, $sort: SortBy) {
  tableViewDefs(pagination: $pagination, where: $where, sort: $sort) {
    content {
      id
      name
      order
      type {
        id
        name
      }
      columns {
        id
        columnType {
          id
          name
        }
        isFilterable
        isSortable
        isDefaultSort
        isVisible
      }
      filters
      sorting
    }
  }
}
    `;

export const useTableViewDefsQuery = <
  TData = TableViewDefsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables?: TableViewDefsQueryVariables,
  options?: Omit<
    UseQueryOptions<TableViewDefsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<TableViewDefsQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<TableViewDefsQuery, TError, TData>({
    queryKey:
      variables === undefined
        ? ['tableViewDefs']
        : ['tableViewDefs', variables],
    queryFn: fetcher<TableViewDefsQuery, TableViewDefsQueryVariables>(
      client,
      TableViewDefsDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useTableViewDefsQuery.document = TableViewDefsDocument;

useTableViewDefsQuery.getKey = (variables?: TableViewDefsQueryVariables) =>
  variables === undefined ? ['tableViewDefs'] : ['tableViewDefs', variables];

export const useInfiniteTableViewDefsQuery = <
  TData = InfiniteData<TableViewDefsQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: TableViewDefsQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<TableViewDefsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      TableViewDefsQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<TableViewDefsQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey:
          optionsQueryKey ?? variables === undefined
            ? ['tableViewDefs.infinite']
            : ['tableViewDefs.infinite', variables],
        queryFn: (metaData) =>
          fetcher<TableViewDefsQuery, TableViewDefsQueryVariables>(
            client,
            TableViewDefsDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteTableViewDefsQuery.getKey = (
  variables?: TableViewDefsQueryVariables,
) =>
  variables === undefined
    ? ['tableViewDefs.infinite']
    : ['tableViewDefs.infinite', variables];

useTableViewDefsQuery.fetcher = (
  client: GraphQLClient,
  variables?: TableViewDefsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<TableViewDefsQuery, TableViewDefsQueryVariables>(
    client,
    TableViewDefsDocument,
    variables,
    headers,
  );

useTableViewDefsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: TableViewDefsQueryVariables) =>
  (mutator: (cacheEntry: TableViewDefsQuery) => TableViewDefsQuery) => {
    const cacheKey = useTableViewDefsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TableViewDefsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TableViewDefsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteTableViewDefsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: TableViewDefsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<TableViewDefsQuery>,
    ) => InfiniteData<TableViewDefsQuery>,
  ) => {
    const cacheKey = useInfiniteTableViewDefsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TableViewDefsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TableViewDefsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
