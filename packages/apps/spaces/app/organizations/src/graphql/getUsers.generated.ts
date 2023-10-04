// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
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
export type GetUsersQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
}>;

export type GetUsersQuery = {
  __typename?: 'Query';
  users: {
    __typename?: 'UserPage';
    totalElements: any;
    content: Array<{
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
    }>;
  };
};

export const GetUsersDocument = `
    query getUsers($pagination: Pagination!, $where: Filter) {
  users(pagination: $pagination, where: $where) {
    content {
      id
      firstName
      lastName
    }
    totalElements
  }
}
    `;
export const useGetUsersQuery = <TData = GetUsersQuery, TError = unknown>(
  client: GraphQLClient,
  variables: GetUsersQueryVariables,
  options?: UseQueryOptions<GetUsersQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetUsersQuery, TError, TData>(
    ['getUsers', variables],
    fetcher<GetUsersQuery, GetUsersQueryVariables>(
      client,
      GetUsersDocument,
      variables,
      headers,
    ),
    options,
  );
useGetUsersQuery.document = GetUsersDocument;

useGetUsersQuery.getKey = (variables: GetUsersQueryVariables) => [
  'getUsers',
  variables,
];
export const useInfiniteGetUsersQuery = <
  TData = GetUsersQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetUsersQueryVariables,
  client: GraphQLClient,
  variables: GetUsersQueryVariables,
  options?: UseInfiniteQueryOptions<GetUsersQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetUsersQuery, TError, TData>(
    ['getUsers.infinite', variables],
    (metaData) =>
      fetcher<GetUsersQuery, GetUsersQueryVariables>(
        client,
        GetUsersDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetUsersQuery.getKey = (variables: GetUsersQueryVariables) => [
  'getUsers.infinite',
  variables,
];
useGetUsersQuery.fetcher = (
  client: GraphQLClient,
  variables: GetUsersQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetUsersQuery, GetUsersQueryVariables>(
    client,
    GetUsersDocument,
    variables,
    headers,
  );
