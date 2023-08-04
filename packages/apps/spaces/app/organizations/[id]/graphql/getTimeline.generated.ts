// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

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
export type GetTimelineQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  from: Types.Scalars['Time'];
  size: Types.Scalars['Int'];
}>;

export type GetTimelineQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    timelineEvents: Array<
      | { __typename: 'Action' }
      | { __typename: 'Analysis' }
      | { __typename: 'Conversation' }
      | {
          __typename: 'InteractionEvent';
          id: string;
          channel?: string | null;
          content?: string | null;
          contentType?: string | null;
          source: Types.DataSource;
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
          sentBy: Array<
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                };
              }
            | {
                __typename: 'EmailParticipant';
                type?: string | null;
                emailParticipant: {
                  __typename?: 'Email';
                  email?: string | null;
                  id: string;
                  contacts: Array<{
                    __typename?: 'Contact';
                    id: string;
                    name?: string | null;
                    firstName?: string | null;
                    lastName?: string | null;
                    emails: Array<{
                      __typename?: 'Email';
                      email?: string | null;
                    }>;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    id: string;
                    firstName: string;
                    lastName: string;
                  }>;
                  organizations: Array<{
                    __typename?: 'Organization';
                    id: string;
                    name: string;
                  }>;
                };
              }
            | { __typename?: 'OrganizationParticipant' }
            | { __typename?: 'PhoneNumberParticipant' }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  firstName: string;
                  lastName: string;
                };
              }
          >;
          sentTo: Array<
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  name?: string | null;
                  id: string;
                  firstName?: string | null;
                  lastName?: string | null;
                };
              }
            | {
                __typename: 'EmailParticipant';
                type?: string | null;
                emailParticipant: {
                  __typename?: 'Email';
                  email?: string | null;
                  id: string;
                  contacts: Array<{
                    __typename?: 'Contact';
                    id: string;
                    name?: string | null;
                    firstName?: string | null;
                    lastName?: string | null;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    id: string;
                    firstName: string;
                    lastName: string;
                  }>;
                  organizations: Array<{
                    __typename?: 'Organization';
                    id: string;
                    name: string;
                  }>;
                };
              }
            | { __typename: 'OrganizationParticipant' }
            | { __typename: 'PhoneNumberParticipant' }
            | {
                __typename: 'UserParticipant';
                type?: string | null;
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  firstName: string;
                  lastName: string;
                };
              }
          >;
        }
      | { __typename: 'InteractionSession' }
      | { __typename: 'Issue' }
      | { __typename: 'Meeting' }
      | { __typename: 'Note' }
      | { __typename: 'PageView' }
    >;
  } | null;
};

export const GetTimelineDocument = `
    query GetTimeline($organizationId: ID!, $from: Time!, $size: Int!) {
  organization(id: $organizationId) {
    timelineEvents(
      from: $from
      size: $size
      timelineEventTypes: [INTERACTION_EVENT]
    ) {
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
        sentBy {
          ... on EmailParticipant {
            __typename
            type
            emailParticipant {
              email
              id
              contacts {
                id
                name
                firstName
                lastName
                emails {
                  email
                }
              }
              users {
                id
                firstName
                lastName
              }
              organizations {
                id
                name
              }
            }
          }
          ... on ContactParticipant {
            __typename
            contactParticipant {
              id
              name
              firstName
              lastName
            }
          }
          ... on UserParticipant {
            __typename
            userParticipant {
              id
              firstName
              lastName
            }
          }
        }
        sentTo {
          __typename
          ... on EmailParticipant {
            __typename
            type
            emailParticipant {
              email
              contacts {
                id
                name
                firstName
                lastName
              }
              users {
                id
                firstName
                lastName
              }
              organizations {
                id
                name
              }
              id
            }
          }
          ... on ContactParticipant {
            __typename
            contactParticipant {
              name
              id
              firstName
              lastName
            }
          }
          ... on UserParticipant {
            __typename
            type
            userParticipant {
              id
              firstName
              lastName
            }
          }
        }
        source
      }
    }
  }
}
    `;
export const useGetTimelineQuery = <TData = GetTimelineQuery, TError = unknown>(
  client: GraphQLClient,
  variables: GetTimelineQueryVariables,
  options?: UseQueryOptions<GetTimelineQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetTimelineQuery, TError, TData>(
    ['GetTimeline', variables],
    fetcher<GetTimelineQuery, GetTimelineQueryVariables>(
      client,
      GetTimelineDocument,
      variables,
      headers,
    ),
    options,
  );
useGetTimelineQuery.document = GetTimelineDocument;

useGetTimelineQuery.getKey = (variables: GetTimelineQueryVariables) => [
  'GetTimeline',
  variables,
];
export const useInfiniteGetTimelineQuery = <
  TData = GetTimelineQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetTimelineQueryVariables,
  client: GraphQLClient,
  variables: GetTimelineQueryVariables,
  options?: UseInfiniteQueryOptions<GetTimelineQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetTimelineQuery, TError, TData>(
    ['GetTimeline.infinite', variables],
    (metaData) =>
      fetcher<GetTimelineQuery, GetTimelineQueryVariables>(
        client,
        GetTimelineDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetTimelineQuery.getKey = (variables: GetTimelineQueryVariables) => [
  'GetTimeline.infinite',
  variables,
];
useGetTimelineQuery.fetcher = (
  client: GraphQLClient,
  variables: GetTimelineQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetTimelineQuery, GetTimelineQueryVariables>(
    client,
    GetTimelineDocument,
    variables,
    headers,
  );
