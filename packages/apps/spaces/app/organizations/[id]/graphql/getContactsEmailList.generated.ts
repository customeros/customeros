// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

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
export type GetContactsEmailListQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Array<Types.SortBy> | Types.SortBy>;
}>;

export type GetContactsEmailListQuery = {
  __typename?: 'Query';
  contacts: {
    __typename?: 'ContactsPage';
    content: Array<{
      __typename?: 'Contact';
      id: string;
      firstName?: string | null;
      lastName?: string | null;
      emails: Array<{
        __typename?: 'Email';
        id: string;
        email?: string | null;
      }>;
    }>;
  };
};

export const GetContactsEmailListDocument = `
    query GetContactsEmailList($pagination: Pagination!, $where: Filter, $sort: [SortBy!]) {
  contacts(pagination: $pagination, where: $where, sort: $sort) {
    content {
      id
      firstName
      lastName
      emails {
        id
        email
      }
    }
  }
}
    `;
export const useGetContactsEmailListQuery = <
  TData = GetContactsEmailListQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetContactsEmailListQueryVariables,
  options?: UseQueryOptions<GetContactsEmailListQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetContactsEmailListQuery, TError, TData>(
    ['GetContactsEmailList', variables],
    fetcher<GetContactsEmailListQuery, GetContactsEmailListQueryVariables>(
      client,
      GetContactsEmailListDocument,
      variables,
      headers,
    ),
    options,
  );
useGetContactsEmailListQuery.document = GetContactsEmailListDocument;

useGetContactsEmailListQuery.getKey = (
  variables: GetContactsEmailListQueryVariables,
) => ['GetContactsEmailList', variables];
useGetContactsEmailListQuery.fetcher = (
  client: GraphQLClient,
  variables: GetContactsEmailListQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetContactsEmailListQuery, GetContactsEmailListQueryVariables>(
    client,
    GetContactsEmailListDocument,
    variables,
    headers,
  );
