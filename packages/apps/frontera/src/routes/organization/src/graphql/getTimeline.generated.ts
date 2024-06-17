// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

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
export type GetTimelineQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  from: Types.Scalars['Time']['input'];
  size: Types.Scalars['Int']['input'];
}>;

export type GetTimelineQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    timelineEventsTotalCount: any;
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
            name?: string | null;
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
          source: Types.DataSource;
          date: any;
          includes: Array<{
            __typename?: 'Attachment';
            id: string;
            mimeType: string;
            fileName: string;
            size: any;
          }>;
          issue?: {
            __typename?: 'Issue';
            externalLinks: Array<{
              __typename?: 'ExternalSystem';
              type: Types.ExternalSystemType;
              externalId?: string | null;
              externalUrl?: string | null;
            }>;
          } | null;
          externalLinks: Array<{
            __typename?: 'ExternalSystem';
            externalUrl?: string | null;
            type: Types.ExternalSystemType;
          }>;
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
                    profilePhotoUrl?: string | null;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    id: string;
                    firstName: string;
                    lastName: string;
                    profilePhotoUrl?: string | null;
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
            | { __typename: 'PhoneNumberParticipant' }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          sentTo: Array<
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
                    profilePhotoUrl?: string | null;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    id: string;
                    firstName: string;
                    lastName: string;
                    profilePhotoUrl?: string | null;
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
            | { __typename: 'PhoneNumberParticipant' }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          interactionSession?: {
            __typename?: 'InteractionSession';
            name: string;
            events: Array<{
              __typename?: 'InteractionEvent';
              id: string;
              channel?: string | null;
              date: any;
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
                        profilePhotoUrl?: string | null;
                      }>;
                      users: Array<{
                        __typename?: 'User';
                        id: string;
                        firstName: string;
                        lastName: string;
                        profilePhotoUrl?: string | null;
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
                      name?: string | null;
                      firstName: string;
                      lastName: string;
                      profilePhotoUrl?: string | null;
                    };
                  }
              >;
            }>;
          } | null;
        }
      | { __typename: 'InteractionSession' }
      | {
          __typename: 'Issue';
          id: string;
          subject?: string | null;
          priority?: string | null;
          appSource: string;
          updatedAt: any;
          createdAt: any;
          description?: string | null;
          issueStatus: string;
          externalLinks: Array<{
            __typename?: 'ExternalSystem';
            type: Types.ExternalSystemType;
            externalId?: string | null;
            externalUrl?: string | null;
          }>;
          submittedBy?:
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
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
            | null;
          followedBy: Array<
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
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          assignedTo: Array<
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
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          reportedBy?:
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
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
            | null;
          interactionEvents: Array<{
            __typename?: 'InteractionEvent';
            content?: string | null;
            contentType?: string | null;
            createdAt: any;
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
                      profilePhotoUrl?: string | null;
                    }>;
                    users: Array<{
                      __typename?: 'User';
                      id: string;
                      firstName: string;
                      lastName: string;
                      profilePhotoUrl?: string | null;
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
                    name?: string | null;
                    firstName: string;
                    lastName: string;
                    profilePhotoUrl?: string | null;
                  };
                }
            >;
          }>;
          comments: Array<{
            __typename?: 'Comment';
            content?: string | null;
            contentType?: string | null;
            createdAt: any;
            createdBy?: {
              __typename?: 'User';
              id: string;
              name?: string | null;
              firstName: string;
              lastName: string;
            } | null;
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
            name?: string | null;
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
            | {
                __typename: 'EmailParticipant';
                emailParticipant: {
                  __typename?: 'Email';
                  rawEmail?: string | null;
                  email?: string | null;
                  contacts: Array<{
                    __typename?: 'Contact';
                    firstName?: string | null;
                    lastName?: string | null;
                    name?: string | null;
                    timezone?: string | null;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    firstName: string;
                    lastName: string;
                  }>;
                };
              }
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
            | {
                __typename: 'EmailParticipant';
                emailParticipant: {
                  __typename?: 'Email';
                  rawEmail?: string | null;
                  email?: string | null;
                  contacts: Array<{
                    __typename?: 'Contact';
                    firstName?: string | null;
                    lastName?: string | null;
                    name?: string | null;
                    timezone?: string | null;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    firstName: string;
                    lastName: string;
                  }>;
                };
              }
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
      | {
          __typename: 'Order';
          id: string;
          confirmedAt?: any | null;
          fulfilledAt?: any | null;
          createdAt: any;
          cancelledAt?: any | null;
        }
      | { __typename: 'PageView' }
    >;
  } | null;
};

export const GetTimelineDocument = `
    query GetTimeline($organizationId: ID!, $from: Time!, $size: Int!) {
  organization(id: $organizationId) {
    timelineEventsTotalCount(
      timelineEventTypes: [INTERACTION_EVENT, MEETING, ACTION, LOG_ENTRY, ISSUE, ORDER]
    )
    timelineEvents(
      from: $from
      size: $size
      timelineEventTypes: [INTERACTION_EVENT, MEETING, ACTION, LOG_ENTRY, ISSUE, ORDER]
    ) {
      __typename
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
            name
            firstName
            lastName
            profilePhotoUrl
          }
        }
        content
      }
      ... on Order {
        id
        confirmedAt
        fulfilledAt
        createdAt
        cancelledAt
      }
      ... on InteractionEvent {
        id
        date: createdAt
        channel
        content
        contentType
        includes {
          id
          mimeType
          fileName
          size
        }
        issue {
          externalLinks {
            type
            externalId
            externalUrl
          }
        }
        externalLinks {
          externalUrl
          type
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
          __typename
          ...InteractionEventParticipantFragment
        }
        sentTo {
          __typename
          ...InteractionEventParticipantFragment
        }
        interactionSession {
          name
          events {
            ... on InteractionEvent {
              id
              date: createdAt
              channel
              sentBy {
                ...InteractionEventParticipantFragment
              }
            }
          }
        }
        source
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
            name
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
      ... on Issue {
        __typename
        id
        subject
        priority
        issueStatus: status
        appSource
        updatedAt
        createdAt
        description
        externalLinks {
          type
          externalId
          externalUrl
        }
        submittedBy {
          ...InteractionEventParticipantFragment
        }
        followedBy {
          ...InteractionEventParticipantFragment
        }
        assignedTo {
          ...InteractionEventParticipantFragment
        }
        reportedBy {
          ...InteractionEventParticipantFragment
        }
        interactionEvents {
          content
          contentType
          createdAt
          sentBy {
            ...InteractionEventParticipantFragment
          }
        }
        comments {
          content
          contentType
          createdAt
          createdBy {
            id
            name
            firstName
            lastName
          }
        }
        issueTags: tags {
          id
          name
        }
      }
    }
  }
}
    ${InteractionEventParticipantFragmentFragmentDoc}
${MeetingParticipantFragmentFragmentDoc}`;

export const useGetTimelineQuery = <TData = GetTimelineQuery, TError = unknown>(
  client: GraphQLClient,
  variables: GetTimelineQueryVariables,
  options?: Omit<
    UseQueryOptions<GetTimelineQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<GetTimelineQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetTimelineQuery, TError, TData>({
    queryKey: ['GetTimeline', variables],
    queryFn: fetcher<GetTimelineQuery, GetTimelineQueryVariables>(
      client,
      GetTimelineDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGetTimelineQuery.document = GetTimelineDocument;

useGetTimelineQuery.getKey = (variables: GetTimelineQueryVariables) => [
  'GetTimeline',
  variables,
];

export const useInfiniteGetTimelineQuery = <
  TData = InfiniteData<GetTimelineQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetTimelineQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetTimelineQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetTimelineQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetTimelineQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? ['GetTimeline.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GetTimelineQuery, GetTimelineQueryVariables>(
            client,
            GetTimelineDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

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

useGetTimelineQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetTimelineQueryVariables) =>
  (mutator: (cacheEntry: GetTimelineQuery) => GetTimelineQuery) => {
    const cacheKey = useGetTimelineQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetTimelineQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetTimelineQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetTimelineQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetTimelineQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetTimelineQuery>,
    ) => InfiniteData<GetTimelineQuery>,
  ) => {
    const cacheKey = useInfiniteGetTimelineQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetTimelineQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetTimelineQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
