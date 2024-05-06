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
export type GetOrganizationsKanbanQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;

export type GetOrganizationsKanbanQuery = {
  __typename?: 'Query';
  dashboardView_Organizations?: {
    __typename?: 'OrganizationPage';
    totalElements: any;
    totalAvailable: any;
    content: Array<{
      __typename?: 'Organization';
      name: string;
      description?: string | null;
      isCustomer?: boolean | null;
      logo?: string | null;
      metadata: {
        __typename?: 'Metadata';
        id: string;
        created: any;
        lastUpdated: any;
      };
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
        profilePhotoUrl?: string | null;
      } | null;
      accountDetails?: {
        __typename?: 'OrgAccountDetails';
        renewalSummary?: {
          __typename?: 'RenewalSummary';
          arrForecast?: number | null;
          maxArrForecast?: number | null;
          renewalLikelihood?: Types.OpportunityRenewalLikelihood | null;
          nextRenewalDate?: any | null;
        } | null;
      } | null;
      contracts?: Array<{
        __typename?: 'Contract';
        id: string;
        contractStatus: Types.ContractStatus;
        contractLineItems?: Array<{
          __typename?: 'ServiceLineItem';
          metadata: { __typename?: 'Metadata'; id: string };
        }> | null;
        opportunities?: Array<{
          __typename?: 'Opportunity';
          id: string;
          amount: number;
          maxAmount: number;
        }> | null;
      }> | null;
      lastTouchpoint?: {
        __typename?: 'LastTouchpoint';
        lastTouchPointTimelineEventId?: string | null;
        lastTouchPointAt?: any | null;
        lastTouchPointType?: Types.LastTouchpointType | null;
        lastTouchPointTimelineEvent?:
          | { __typename?: 'Action' }
          | { __typename?: 'Analysis' }
          | {
              __typename?: 'InteractionEvent';
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
          | { __typename?: 'InteractionSession' }
          | { __typename?: 'Issue' }
          | {
              __typename?: 'LogEntry';
              id: string;
              createdBy?: {
                __typename?: 'User';
                lastName: string;
                firstName: string;
              } | null;
              tags: Array<{ __typename?: 'Tag'; id: string; name: string }>;
            }
          | {
              __typename?: 'Meeting';
              id: string;
              name?: string | null;
              attendedBy: Array<
                | { __typename: 'ContactParticipant' }
                | { __typename: 'EmailParticipant' }
                | { __typename: 'OrganizationParticipant' }
                | { __typename: 'UserParticipant' }
              >;
            }
          | { __typename?: 'Note' }
          | { __typename?: 'Order' }
          | { __typename?: 'PageView' }
          | null;
      } | null;
    }>;
  } | null;
};

export const GetOrganizationsKanbanDocument = `
    query getOrganizationsKanban($pagination: Pagination!, $sort: SortBy) {
  dashboardView_Organizations(
    pagination: $pagination
    where: {filter: {property: "IS_CUSTOMER", value: false, operation: IN}}
    sort: $sort
  ) {
    content {
      name
      metadata {
        id
        created
        lastUpdated
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
        profilePhotoUrl
      }
      description
      isCustomer
      logo
      accountDetails {
        renewalSummary {
          arrForecast
          maxArrForecast
          renewalLikelihood
          nextRenewalDate
        }
      }
      contracts {
        id
        contractStatus
        contractLineItems {
          metadata {
            id
          }
        }
        opportunities {
          id
          amount
          maxAmount
        }
      }
      lastTouchpoint {
        lastTouchPointTimelineEventId
        lastTouchPointAt
        lastTouchPointType
        lastTouchPointTimelineEvent {
          ... on LogEntry {
            id
            createdBy {
              lastName
              firstName
            }
            tags {
              id
              name
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
          ... on Meeting {
            id
            name
            attendedBy {
              __typename
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

export const useGetOrganizationsKanbanQuery = <
  TData = GetOrganizationsKanbanQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationsKanbanQueryVariables,
  options?: Omit<
    UseQueryOptions<GetOrganizationsKanbanQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<
      GetOrganizationsKanbanQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetOrganizationsKanbanQuery, TError, TData>({
    queryKey: ['getOrganizationsKanban', variables],
    queryFn: fetcher<
      GetOrganizationsKanbanQuery,
      GetOrganizationsKanbanQueryVariables
    >(client, GetOrganizationsKanbanDocument, variables, headers),
    ...options,
  });
};

useGetOrganizationsKanbanQuery.document = GetOrganizationsKanbanDocument;

useGetOrganizationsKanbanQuery.getKey = (
  variables: GetOrganizationsKanbanQueryVariables,
) => ['getOrganizationsKanban', variables];

export const useInfiniteGetOrganizationsKanbanQuery = <
  TData = InfiniteData<GetOrganizationsKanbanQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetOrganizationsKanbanQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetOrganizationsKanbanQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetOrganizationsKanbanQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetOrganizationsKanbanQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? [
          'getOrganizationsKanban.infinite',
          variables,
        ],
        queryFn: (metaData) =>
          fetcher<
            GetOrganizationsKanbanQuery,
            GetOrganizationsKanbanQueryVariables
          >(
            client,
            GetOrganizationsKanbanDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetOrganizationsKanbanQuery.getKey = (
  variables: GetOrganizationsKanbanQueryVariables,
) => ['getOrganizationsKanban.infinite', variables];

useGetOrganizationsKanbanQuery.fetcher = (
  client: GraphQLClient,
  variables: GetOrganizationsKanbanQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetOrganizationsKanbanQuery, GetOrganizationsKanbanQueryVariables>(
    client,
    GetOrganizationsKanbanDocument,
    variables,
    headers,
  );

useGetOrganizationsKanbanQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables?: GetOrganizationsKanbanQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: GetOrganizationsKanbanQuery,
    ) => GetOrganizationsKanbanQuery,
  ) => {
    const cacheKey = useGetOrganizationsKanbanQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetOrganizationsKanbanQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetOrganizationsKanbanQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetOrganizationsKanbanQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables?: GetOrganizationsKanbanQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetOrganizationsKanbanQuery>,
    ) => InfiniteData<GetOrganizationsKanbanQuery>,
  ) => {
    const cacheKey = useInfiniteGetOrganizationsKanbanQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
        cacheKey,
      );
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
