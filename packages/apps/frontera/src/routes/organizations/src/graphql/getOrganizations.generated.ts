// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

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
export type GetOrganizationsQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;

export type GetOrganizationsQuery = {
  __typename?: 'Query';
  dashboardView_Organizations?: {
    __typename?: 'OrganizationPage';
    totalElements: any;
    totalAvailable: any;
    content: Array<{
      __typename?: 'Organization';
      name: string;
      description?: string | null;
      industry?: string | null;
      website?: string | null;
      domains: Array<string>;
      isCustomer?: boolean | null;
      logo?: string | null;
      icon?: string | null;
      leadSource?: string | null;
      employees?: any | null;
      yearFounded?: any | null;
      metadata: { __typename?: 'Metadata'; id: string; created: any };
      parentCompanies: Array<{
        __typename?: 'LinkedOrganization';
        organization: {
          __typename?: 'Organization';
          name: string;
          metadata: { __typename?: 'Metadata'; id: string };
        };
      }>;
      owner?: {
        __typename?: 'User';
        id: string;
        firstName: string;
        lastName: string;
        name?: string | null;
      } | null;
      socialMedia: Array<{ __typename?: 'Social'; id: string; url: string }>;
      accountDetails?: {
        __typename?: 'OrgAccountDetails';
        renewalSummary?: {
          __typename?: 'RenewalSummary';
          arrForecast?: number | null;
          maxArrForecast?: number | null;
          renewalLikelihood?: Types.OpportunityRenewalLikelihood | null;
          nextRenewalDate?: any | null;
        } | null;
        onboarding?: {
          __typename?: 'OnboardingDetails';
          status: Types.OnboardingStatus;
          comments?: string | null;
          updatedAt?: any | null;
        } | null;
      } | null;
      locations: Array<{
        __typename?: 'Location';
        id: string;
        name?: string | null;
        country?: string | null;
        region?: string | null;
        locality?: string | null;
        zip?: string | null;
        street?: string | null;
        postalCode?: string | null;
        houseNumber?: string | null;
        rawAddress?: string | null;
      }>;
      lastTouchpoint?: {
        __typename?: 'LastTouchpoint';
        lastTouchPointTimelineEventId?: string | null;
        lastTouchPointAt?: any | null;
        lastTouchPointType?: Types.LastTouchpointType | null;
        lastTouchPointTimelineEvent?:
          | {
              __typename: 'Action';
              id: string;
              actionType: Types.ActionType;
              createdAt: any;
              source: Types.DataSource;
              createdBy?: {
                __typename?: 'User';
                id: string;
                firstName: string;
                lastName: string;
              } | null;
            }
          | { __typename: 'Analysis'; id: string }
          | {
              __typename: 'InteractionEvent';
              id: string;
              channel?: string | null;
              eventType?: string | null;
              externalLinks: Array<{
                __typename?: 'ExternalSystem';
                type: Types.ExternalSystemType;
              }>;
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
                      id: string;
                      email?: string | null;
                      rawEmail?: string | null;
                    };
                  }
                | {
                    __typename: 'JobRoleParticipant';
                    jobRoleParticipant: {
                      __typename?: 'JobRole';
                      contact?: {
                        __typename?: 'Contact';
                        id: string;
                        name?: string | null;
                        firstName?: string | null;
                        lastName?: string | null;
                      } | null;
                    };
                  }
                | { __typename: 'OrganizationParticipant' }
                | { __typename: 'PhoneNumberParticipant' }
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
            }
          | { __typename: 'InteractionSession' }
          | { __typename: 'Issue'; id: string; createdAt: any; updatedAt: any }
          | {
              __typename: 'LogEntry';
              id: string;
              createdBy?: {
                __typename?: 'User';
                lastName: string;
                firstName: string;
              } | null;
            }
          | {
              __typename: 'Meeting';
              id: string;
              name?: string | null;
              attendedBy: Array<
                | { __typename: 'ContactParticipant' }
                | { __typename: 'EmailParticipant' }
                | { __typename: 'OrganizationParticipant' }
                | { __typename: 'UserParticipant' }
              >;
            }
          | {
              __typename: 'Note';
              id: string;
              createdBy?: {
                __typename?: 'User';
                firstName: string;
                lastName: string;
              } | null;
            }
          | { __typename: 'Order' }
          | { __typename: 'PageView'; id: string }
          | null;
      } | null;
    }>;
  } | null;
};

export const GetOrganizationsDocument = `
    query getOrganizations($pagination: Pagination!, $where: Filter, $sort: SortBy) {
  dashboardView_Organizations(pagination: $pagination, where: $where, sort: $sort) {
    content {
      name
      metadata {
        id
        created
      }
      parentCompanies {
        organization {
          metadata {
            id
          }
          name
        }
      }
      owner {
        id
        firstName
        lastName
        name
      }
      description
      industry
      website
      domains
      isCustomer
      logo
      icon
      leadSource
      socialMedia {
        id
        url
      }
      employees
      yearFounded
      accountDetails {
        renewalSummary {
          arrForecast
          maxArrForecast
          renewalLikelihood
          nextRenewalDate
        }
        onboarding {
          status
          comments
          updatedAt
        }
      }
      locations {
        id
        name
        country
        region
        locality
        zip
        street
        postalCode
        houseNumber
        rawAddress
      }
      lastTouchpoint {
        lastTouchPointTimelineEventId
        lastTouchPointAt
        lastTouchPointType
        lastTouchPointTimelineEvent {
          __typename
          ... on PageView {
            id
          }
          ... on Issue {
            id
            createdAt
            updatedAt
          }
          ... on LogEntry {
            id
            createdBy {
              lastName
              firstName
            }
          }
          ... on Note {
            id
            createdBy {
              firstName
              lastName
            }
          }
          ... on InteractionEvent {
            id
            channel
            eventType
            externalLinks {
              type
            }
            sentBy {
              __typename
              ... on EmailParticipant {
                type
                emailParticipant {
                  id
                  email
                  rawEmail
                }
              }
              ... on ContactParticipant {
                contactParticipant {
                  id
                  name
                  firstName
                  lastName
                }
              }
              ... on JobRoleParticipant {
                jobRoleParticipant {
                  contact {
                    id
                    name
                    firstName
                    lastName
                  }
                }
              }
              ... on UserParticipant {
                userParticipant {
                  id
                  firstName
                  lastName
                }
              }
            }
          }
          ... on Analysis {
            id
          }
          ... on Meeting {
            id
            name
            attendedBy {
              __typename
            }
          }
          ... on Action {
            id
            actionType
            createdAt
            source
            actionType
            createdBy {
              id
              firstName
              lastName
            }
          }
        }
      }
    }
    totalElements
    totalAvailable
  }
}
    `;

export const useGetOrganizationsQuery = <
  TData = GetOrganizationsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationsQueryVariables,
  options?: Omit<
    UseQueryOptions<GetOrganizationsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<
      GetOrganizationsQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetOrganizationsQuery, TError, TData>({
    queryKey: ['getOrganizations', variables],
    queryFn: fetcher<GetOrganizationsQuery, GetOrganizationsQueryVariables>(
      client,
      GetOrganizationsDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGetOrganizationsQuery.document = GetOrganizationsDocument;

useGetOrganizationsQuery.getKey = (
  variables: GetOrganizationsQueryVariables,
) => ['getOrganizations', variables];

export const useInfiniteGetOrganizationsQuery = <
  TData = InfiniteData<GetOrganizationsQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationsQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetOrganizationsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetOrganizationsQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetOrganizationsQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? ['getOrganizations.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GetOrganizationsQuery, GetOrganizationsQueryVariables>(
            client,
            GetOrganizationsDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetOrganizationsQuery.getKey = (
  variables: GetOrganizationsQueryVariables,
) => ['getOrganizations.infinite', variables];

useGetOrganizationsQuery.fetcher = (
  client: GraphQLClient,
  variables: GetOrganizationsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetOrganizationsQuery, GetOrganizationsQueryVariables>(
    client,
    GetOrganizationsDocument,
    variables,
    headers,
  );

useGetOrganizationsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetOrganizationsQueryVariables) =>
  (mutator: (cacheEntry: GetOrganizationsQuery) => GetOrganizationsQuery) => {
    const cacheKey = useGetOrganizationsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetOrganizationsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetOrganizationsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetOrganizationsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetOrganizationsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetOrganizationsQuery>,
    ) => InfiniteData<GetOrganizationsQuery>,
  ) => {
    const cacheKey = useInfiniteGetOrganizationsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
