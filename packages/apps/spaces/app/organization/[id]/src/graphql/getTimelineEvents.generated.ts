// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  InteractionEventParticipantFragmentFragmentDoc,
  MeetingParticipantFragmentFragmentDoc,
} from './participantsFragment.generated';
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
    | {
        __typename: 'Action';
        id: string;
        actionType: Types.ActionType;
        appSource: string;
        createdAt: any;
        metadata?: string | null;
        content?: string | null;
        actionCreatedBy?: {
          __typename: 'User';
          id: string;
          firstName: string;
          lastName: string;
          profilePhotoUrl?: string | null;
        } | null;
      }
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
    | {
        __typename: 'Issue';
        id: string;
        subject?: string | null;
        priority?: string | null;
        appSource: string;
        createdAt: any;
        description?: string | null;
        issueStatus: string;
        externalLinks: Array<{
          __typename?: 'ExternalSystem';
          externalId?: string | null;
          externalUrl?: string | null;
        }>;
        interactionEvents: Array<{
          __typename?: 'InteractionEvent';
          id: string;
          content?: string | null;
          sentBy: Array<
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                  profilePhotoUrl?: string | null;
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
            | {
                __typename: 'JobRoleParticipant';
                jobRoleParticipant: {
                  __typename?: 'JobRole';
                  id: string;
                  contact?: {
                    __typename?: 'Contact';
                    id: string;
                    name?: string | null;
                    firstName?: string | null;
                    lastName?: string | null;
                    profilePhotoUrl?: string | null;
                  } | null;
                };
              }
            | {
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | { __typename?: 'PhoneNumberParticipant' }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
        }>;
        issueTags?: Array<{
          __typename?: 'Tag';
          id: string;
          name: string;
        } | null> | null;
      }
    | {
        __typename: 'LogEntry';
        id: string;
        createdAt: any;
        updatedAt: any;
        source: Types.DataSource;
        content?: string | null;
        contentType?: string | null;
        logEntryStartedAt: any;
        logEntryCreatedBy?: {
          __typename: 'User';
          id: string;
          firstName: string;
          lastName: string;
          profilePhotoUrl?: string | null;
          emails?: Array<{
            __typename?: 'Email';
            email?: string | null;
          }> | null;
        } | null;
        tags: Array<{ __typename?: 'Tag'; id: string; name: string }>;
        externalLinks: Array<{
          __typename?: 'ExternalSystem';
          type: Types.ExternalSystemType;
          externalUrl?: string | null;
          externalSource?: string | null;
        }>;
      }
    | {
        __typename: 'Meeting';
        id: string;
        name?: string | null;
        createdAt: any;
        updatedAt: any;
        startedAt?: any | null;
        endedAt?: any | null;
        agenda?: string | null;
        status: Types.MeetingStatus;
        attendedBy: Array<
          | {
              __typename: 'ContactParticipant';
              contactParticipant: {
                __typename?: 'Contact';
                id: string;
                name?: string | null;
                firstName?: string | null;
                lastName?: string | null;
                profilePhotoUrl?: string | null;
                timezone?: string | null;
                emails: Array<{
                  __typename?: 'Email';
                  id: string;
                  email?: string | null;
                  rawEmail?: string | null;
                  primary: boolean;
                }>;
              };
            }
          | { __typename?: 'EmailParticipant' }
          | {
              __typename: 'OrganizationParticipant';
              organizationParticipant: {
                __typename?: 'Organization';
                id: string;
                name: string;
                emails: Array<{
                  __typename?: 'Email';
                  id: string;
                  email?: string | null;
                  rawEmail?: string | null;
                  primary: boolean;
                }>;
              };
            }
          | {
              __typename: 'UserParticipant';
              userParticipant: {
                __typename?: 'User';
                id: string;
                firstName: string;
                lastName: string;
                profilePhotoUrl?: string | null;
                emails?: Array<{
                  __typename?: 'Email';
                  id: string;
                  email?: string | null;
                  rawEmail?: string | null;
                  primary: boolean;
                }> | null;
              };
            }
        >;
        createdBy: Array<
          | {
              __typename: 'ContactParticipant';
              contactParticipant: {
                __typename?: 'Contact';
                id: string;
                name?: string | null;
                firstName?: string | null;
                lastName?: string | null;
                profilePhotoUrl?: string | null;
                timezone?: string | null;
                emails: Array<{
                  __typename?: 'Email';
                  id: string;
                  email?: string | null;
                  rawEmail?: string | null;
                  primary: boolean;
                }>;
              };
            }
          | { __typename?: 'EmailParticipant' }
          | {
              __typename: 'OrganizationParticipant';
              organizationParticipant: {
                __typename?: 'Organization';
                id: string;
                name: string;
                emails: Array<{
                  __typename?: 'Email';
                  id: string;
                  email?: string | null;
                  rawEmail?: string | null;
                  primary: boolean;
                }>;
              };
            }
          | {
              __typename: 'UserParticipant';
              userParticipant: {
                __typename?: 'User';
                id: string;
                firstName: string;
                lastName: string;
                profilePhotoUrl?: string | null;
                emails?: Array<{
                  __typename?: 'Email';
                  id: string;
                  email?: string | null;
                  rawEmail?: string | null;
                  primary: boolean;
                }> | null;
              };
            }
        >;
        note: Array<{
          __typename?: 'Note';
          id: string;
          content?: string | null;
        }>;
      }
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
    ... on Issue {
      __typename
      id
      subject
      priority
      issueStatus: status
      appSource
      createdAt
      description
      externalLinks {
        externalId
        externalUrl
      }
      interactionEvents {
        id
        sentBy {
          ...InteractionEventParticipantFragment
        }
        content
      }
      issueTags: tags {
        id
        name
      }
    }
    ... on Action {
      __typename
      id
      actionType
      appSource
      createdAt
      metadata
      actionCreatedBy: createdBy {
        ... on User {
          __typename
          id
          firstName
          lastName
          profilePhotoUrl
        }
      }
      content
    }
    ... on Meeting {
      id
      name
      createdAt
      updatedAt
      startedAt
      endedAt
      attendedBy {
        ...MeetingParticipantFragment
      }
      createdBy {
        ...MeetingParticipantFragment
      }
      note {
        id
        content
      }
      agenda
      status
    }
    ... on LogEntry {
      id
      createdAt
      updatedAt
      logEntryStartedAt: startedAt
      logEntryCreatedBy: createdBy {
        ... on User {
          __typename
          id
          firstName
          lastName
          profilePhotoUrl
          emails {
            email
          }
        }
      }
      tags {
        id
        name
      }
      source
      content
      contentType
      externalLinks {
        type
        externalUrl
        externalSource
      }
    }
  }
}
    ${InteractionEventParticipantFragmentFragmentDoc}
${MeetingParticipantFragmentFragmentDoc}`;
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
