// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { InteractionEventParticipantFragmentFragmentDoc, MeetingParticipantFragmentFragmentDoc } from './participantsFragment.generated';
import { useQuery, useInfiniteQuery, UseQueryOptions, UseInfiniteQueryOptions, InfiniteData } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type GetIssuesQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  from: Types.Scalars['Time']['input'];
  size: Types.Scalars['Int']['input'];
}>;


export type GetIssuesQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', name: string, timelineEventsTotalCount: any, timelineEvents: Array<{ __typename?: 'Action' } | { __typename?: 'Analysis' } | { __typename?: 'InteractionEvent' } | { __typename?: 'InteractionSession' } | { __typename?: 'Issue', id: string, subject?: string | null, status: string, appSource: string, source: Types.DataSource, updatedAt: any, createdAt: any, externalLinks: Array<{ __typename?: 'ExternalSystem', externalId?: string | null, externalUrl?: string | null }>, submittedBy?: { __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, profilePhotoUrl?: string | null } } | { __typename: 'OrganizationParticipant', organizationParticipant: { __typename?: 'Organization', id: string, name: string } } | { __typename: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, name?: string | null, firstName: string, lastName: string, profilePhotoUrl?: string | null } } | null, reportedBy?: { __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, profilePhotoUrl?: string | null } } | { __typename: 'OrganizationParticipant', organizationParticipant: { __typename?: 'Organization', id: string, name: string } } | { __typename: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, name?: string | null, firstName: string, lastName: string, profilePhotoUrl?: string | null } } | null } | { __typename?: 'LogEntry' } | { __typename?: 'Meeting' } | { __typename?: 'Note' } | { __typename?: 'Order' } | { __typename?: 'PageView' }> } | null };



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
        source
        updatedAt
        externalLinks {
          externalId
          externalUrl
        }
        createdAt
        submittedBy {
          ...InteractionEventParticipantFragment
        }
        reportedBy {
          ...InteractionEventParticipantFragment
        }
      }
    }
  }
}
    ${InteractionEventParticipantFragmentFragmentDoc}`;

export const useGetIssuesQuery = <
      TData = GetIssuesQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetIssuesQueryVariables,
      options?: Omit<UseQueryOptions<GetIssuesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetIssuesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetIssuesQuery, TError, TData>(
      {
    queryKey: ['GetIssues', variables],
    queryFn: fetcher<GetIssuesQuery, GetIssuesQueryVariables>(client, GetIssuesDocument, variables, headers),
    ...options
  }
    )};

useGetIssuesQuery.document = GetIssuesDocument;

useGetIssuesQuery.getKey = (variables: GetIssuesQueryVariables) => ['GetIssues', variables];

export const useInfiniteGetIssuesQuery = <
      TData = InfiniteData<GetIssuesQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetIssuesQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetIssuesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetIssuesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetIssuesQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['GetIssues.infinite', variables],
      queryFn: (metaData) => fetcher<GetIssuesQuery, GetIssuesQueryVariables>(client, GetIssuesDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetIssuesQuery.getKey = (variables: GetIssuesQueryVariables) => ['GetIssues.infinite', variables];


useGetIssuesQuery.fetcher = (client: GraphQLClient, variables: GetIssuesQueryVariables, headers?: RequestInit['headers']) => fetcher<GetIssuesQuery, GetIssuesQueryVariables>(client, GetIssuesDocument, variables, headers);


useGetIssuesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetIssuesQueryVariables) =>
  (mutator: (cacheEntry: GetIssuesQuery) => GetIssuesQuery) => {
    const cacheKey = useGetIssuesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetIssuesQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetIssuesQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetIssuesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetIssuesQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetIssuesQuery>) => InfiniteData<GetIssuesQuery>) => {
    const cacheKey = useInfiniteGetIssuesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetIssuesQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetIssuesQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }