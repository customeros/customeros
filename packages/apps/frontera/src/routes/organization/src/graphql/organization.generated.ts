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
export type OrganizationQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type OrganizationQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    id: string;
    name: string;
    description?: string | null;
    domains: Array<string>;
    website?: string | null;
    industry?: string | null;
    subIndustry?: string | null;
    industryGroup?: string | null;
    targetAudience?: string | null;
    valueProposition?: string | null;
    lastFundingRound?: Types.FundingRound | null;
    lastFundingAmount?: string | null;
    isPublic?: boolean | null;
    market?: Types.Market | null;
    employees?: any | null;
    referenceId?: string | null;
    customerOsId: string;
    hide: boolean;
    slackChannelId?: string | null;
    stage?: Types.OrganizationStage | null;
    relationship?: Types.OrganizationRelationship | null;
    socials: Array<{ __typename?: 'Social'; id: string; url: string }>;
    subsidiaryOf: Array<{
      __typename?: 'LinkedOrganization';
      organization: { __typename?: 'Organization'; id: string; name: string };
    }>;
    subsidiaries: Array<{
      __typename?: 'LinkedOrganization';
      organization: { __typename?: 'Organization'; id: string; name: string };
    }>;
    owner?: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
    } | null;
    accountDetails?: {
      __typename?: 'OrgAccountDetails';
      onboarding?: {
        __typename?: 'OnboardingDetails';
        status: Types.OnboardingStatus;
        comments?: string | null;
        updatedAt?: any | null;
      } | null;
    } | null;
  } | null;
};

export const OrganizationDocument = `
    query Organization($id: ID!) {
  organization(id: $id) {
    id
    name
    description
    domains
    website
    industry
    subIndustry
    industryGroup
    targetAudience
    valueProposition
    lastFundingRound
    lastFundingAmount
    isPublic
    market
    employees
    referenceId
    customerOsId
    hide
    slackChannelId
    stage
    relationship
    socials {
      id
      url
    }
    subsidiaryOf {
      organization {
        id
        name
      }
    }
    subsidiaries {
      organization {
        id
        name
      }
    }
    owner {
      id
      firstName
      lastName
    }
    accountDetails {
      onboarding {
        status
        comments
        updatedAt
      }
    }
  }
}
    `;

export const useOrganizationQuery = <
  TData = OrganizationQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: OrganizationQueryVariables,
  options?: Omit<
    UseQueryOptions<OrganizationQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<OrganizationQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<OrganizationQuery, TError, TData>({
    queryKey: ['Organization', variables],
    queryFn: fetcher<OrganizationQuery, OrganizationQueryVariables>(
      client,
      OrganizationDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useOrganizationQuery.document = OrganizationDocument;

useOrganizationQuery.getKey = (variables: OrganizationQueryVariables) => [
  'Organization',
  variables,
];

export const useInfiniteOrganizationQuery = <
  TData = InfiniteData<OrganizationQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: OrganizationQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<OrganizationQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      OrganizationQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<OrganizationQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? ['Organization.infinite', variables],
        queryFn: (metaData) =>
          fetcher<OrganizationQuery, OrganizationQueryVariables>(
            client,
            OrganizationDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteOrganizationQuery.getKey = (
  variables: OrganizationQueryVariables,
) => ['Organization.infinite', variables];

useOrganizationQuery.fetcher = (
  client: GraphQLClient,
  variables: OrganizationQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<OrganizationQuery, OrganizationQueryVariables>(
    client,
    OrganizationDocument,
    variables,
    headers,
  );

useOrganizationQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: OrganizationQueryVariables) =>
  (mutator: (cacheEntry: OrganizationQuery) => OrganizationQuery) => {
    const cacheKey = useOrganizationQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<OrganizationQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<OrganizationQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteOrganizationQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: OrganizationQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<OrganizationQuery>,
    ) => InfiniteData<OrganizationQuery>,
  ) => {
    const cacheKey = useInfiniteOrganizationQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<OrganizationQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<OrganizationQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
