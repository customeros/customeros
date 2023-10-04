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
export type GetTimelineEventsQueryVariables = Types.Exact<{
  ids: Array<Types.Scalars['ID']> | Types.Scalars['ID'];
}>;

export type GetTimelineEventsQuery = {
  __typename?: 'Query';
  timelineEvents: Array<
    | { __typename: 'Action' }
    | { __typename: 'Analysis' }
    | {
        __typename: 'InteractionEvent';
        id: string;
        channel?: string | null;
        content?: string | null;
        contentType?: string | null;
        date: any;
        interactionSession?: {
          __typename?: 'InteractionSession';
          name: string;
        } | null;
        issue?: {
          __typename?: 'Issue';
          externalLinks: Array<{
            __typename?: 'ExternalSystem';
            type: Types.ExternalSystemType;
            externalId?: string | null;
            externalUrl?: string | null;
          }>;
        } | null;
        repliesTo?: { __typename?: 'InteractionEvent'; id: string } | null;
        summary?: {
          __typename?: 'Analysis';
          id: string;
          content?: string | null;
          contentType?: string | null;
        } | null;
        actionItems?: Array<{
          __typename?: 'ActionItem';
          id: string;
          content: string;
        }> | null;
      }
    | { __typename: 'InteractionSession' }
    | { __typename: 'Issue' }
    | { __typename: 'LogEntry' }
    | { __typename: 'Meeting' }
    | { __typename: 'Note' }
    | { __typename: 'PageView' }
  >;
};

export const GetTimelineEventsDocument = `
    query getTimelineEvents($ids: [ID!]!) {
  timelineEvents(ids: $ids) {
    __typename
    ... on InteractionEvent {
      id
      date: createdAt
      channel
      interactionSession {
        name
      }
      content
      contentType
      issue {
        externalLinks {
          type
          externalId
          externalUrl
        }
      }
      repliesTo {
        id
      }
      summary {
        id
        content
        contentType
      }
      actionItems {
        id
        content
      }
    }
  }
}
    `;
export const useGetTimelineEventsQuery = <
  TData = GetTimelineEventsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetTimelineEventsQueryVariables,
  options?: UseQueryOptions<GetTimelineEventsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetTimelineEventsQuery, TError, TData>(
    ['getTimelineEvents', variables],
    fetcher<GetTimelineEventsQuery, GetTimelineEventsQueryVariables>(
      client,
      GetTimelineEventsDocument,
      variables,
      headers,
    ),
    options,
  );
useGetTimelineEventsQuery.document = GetTimelineEventsDocument;

useGetTimelineEventsQuery.getKey = (
  variables: GetTimelineEventsQueryVariables,
) => ['getTimelineEvents', variables];
export const useInfiniteGetTimelineEventsQuery = <
  TData = GetTimelineEventsQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetTimelineEventsQueryVariables,
  client: GraphQLClient,
  variables: GetTimelineEventsQueryVariables,
  options?: UseInfiniteQueryOptions<GetTimelineEventsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetTimelineEventsQuery, TError, TData>(
    ['getTimelineEvents.infinite', variables],
    (metaData) =>
      fetcher<GetTimelineEventsQuery, GetTimelineEventsQueryVariables>(
        client,
        GetTimelineEventsDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetTimelineEventsQuery.getKey = (
  variables: GetTimelineEventsQueryVariables,
) => ['getTimelineEvents.infinite', variables];
useGetTimelineEventsQuery.fetcher = (
  client: GraphQLClient,
  variables: GetTimelineEventsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetTimelineEventsQuery, GetTimelineEventsQueryVariables>(
    client,
    GetTimelineEventsDocument,
    variables,
    headers,
  );
