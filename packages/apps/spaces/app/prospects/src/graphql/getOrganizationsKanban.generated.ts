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
      relationship?: Types.OrganizationRelationship | null;
      stage?: Types.OrganizationStage | null;
      website?: string | null;
      lastFundingRound?: Types.FundingRound | null;
      employees?: any | null;
      contactCount: any;
      logo?: string | null;
      timelineEventsTotalCount: any;
      metadata: {
        __typename?: 'Metadata';
        id: string;
        created: any;
        lastUpdated: any;
      };
      owner?: {
        __typename?: 'User';
        id: string;
        firstName: string;
        lastName: string;
        name?: string | null;
        profilePhotoUrl?: string | null;
      } | null;
      lastTouchpoint?: {
        __typename?: 'LastTouchpoint';
        lastTouchPointAt?: any | null;
      } | null;
      accountDetails?: {
        __typename?: 'OrgAccountDetails';
        renewalSummary?: {
          __typename?: 'RenewalSummary';
          maxArrForecast?: number | null;
        } | null;
      } | null;
    }>;
  } | null;
};

export const GetOrganizationsKanbanDocument = `
    query getOrganizationsKanban($pagination: Pagination!, $sort: SortBy) {
  dashboardView_Organizations(
    pagination: $pagination
    where: {AND: [{filter: {property: "RELATIONSHIP", value: "PROSPECT"}}, {filter: {property: "STAGE", value: ["TARGET", "INTERESTED", "ENGAGED", "CLOSED_WON"], operation: IN}}]}
    sort: $sort
  ) {
    content {
      name
      metadata {
        id
        created
        lastUpdated
      }
      relationship
      stage
      website
      lastFundingRound
      owner {
        id
        firstName
        lastName
        name
        profilePhotoUrl
      }
      employees
      contactCount
      logo
      timelineEventsTotalCount(timelineEventTypes: [MEETING, INTERACTION_EVENT])
      lastTouchpoint {
        lastTouchPointAt
      }
      accountDetails {
        renewalSummary {
          maxArrForecast
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
