// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

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
export type GetIssuesQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  from: Types.Scalars['Time'];
  size: Types.Scalars['Int'];
}>;

export type GetIssuesQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    name: string;
    timelineEventsTotalCount: any;
    timelineEvents: Array<
      | { __typename?: 'Action' }
      | { __typename?: 'Analysis' }
      | { __typename?: 'InteractionEvent' }
      | { __typename?: 'InteractionSession' }
      | {
          __typename?: 'Issue';
          id: string;
          subject?: string | null;
          status: string;
          appSource: string;
          createdAt: any;
          externalLinks: Array<{
            __typename?: 'ExternalSystem';
            externalId?: string | null;
          }>;
        }
      | { __typename?: 'LogEntry' }
      | { __typename?: 'Meeting' }
      | { __typename?: 'Note' }
      | { __typename?: 'PageView' }
    >;
  } | null;
};

export const GetIssuesDocument = `
    query GetIssues($organizationId: ID!, $from: Time!, $size: Int!) {
  organization(id: $organizationId) {
    name
    timelineEventsTotalCount(timelineEventTypes: [ISSUE])
    timelineEvents(from: $from, size: $size, timelineEventTypes: [ISSUE]) {
      ... on Issue {
        id
        subject
        status
        appSource
        externalLinks {
          externalId
        }
        createdAt
      }
    }
  }
}
    `;
export const useGetIssuesQuery = <TData = GetIssuesQuery, TError = unknown>(
  client: GraphQLClient,
  variables: GetIssuesQueryVariables,
  options?: UseQueryOptions<GetIssuesQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetIssuesQuery, TError, TData>(
    ['GetIssues', variables],
    fetcher<GetIssuesQuery, GetIssuesQueryVariables>(
      client,
      GetIssuesDocument,
      variables,
      headers,
    ),
    options,
  );
useGetIssuesQuery.document = GetIssuesDocument;

useGetIssuesQuery.getKey = (variables: GetIssuesQueryVariables) => [
  'GetIssues',
  variables,
];
export const useInfiniteGetIssuesQuery = <
  TData = GetIssuesQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetIssuesQueryVariables,
  client: GraphQLClient,
  variables: GetIssuesQueryVariables,
  options?: UseInfiniteQueryOptions<GetIssuesQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetIssuesQuery, TError, TData>(
    ['GetIssues.infinite', variables],
    (metaData) =>
      fetcher<GetIssuesQuery, GetIssuesQueryVariables>(
        client,
        GetIssuesDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetIssuesQuery.getKey = (variables: GetIssuesQueryVariables) => [
  'GetIssues.infinite',
  variables,
];
useGetIssuesQuery.fetcher = (
  client: GraphQLClient,
  variables: GetIssuesQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetIssuesQuery, GetIssuesQueryVariables>(
    client,
    GetIssuesDocument,
    variables,
    headers,
  );
