// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

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
export type OrganizationOnboardingPlansQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
}>;

export type OrganizationOnboardingPlansQuery = {
  __typename?: 'Query';
  organizationPlansForOrganization: Array<{
    __typename?: 'OrganizationPlan';
    id: string;
    name: string;
    retired: boolean;
    masterPlanId: string;
    milestones: Array<{
      __typename?: 'OrganizationPlanMilestone';
      id: string;
      name: string;
      order: any;
      dueDate: any;
      optional: boolean;
      retired: boolean;
      items: Array<{
        __typename?: 'MilestoneItem';
        status: string;
        text: string;
        updatedAt: any;
      }>;
      statusDetails: {
        __typename?: 'StatusDetails';
        updatedAt: any;
        status: string;
        text: string;
      };
    }>;
    statusDetails: {
      __typename?: 'StatusDetails';
      updatedAt: any;
      status: string;
      text: string;
    };
  }>;
};

export const OrganizationOnboardingPlansDocument = `
    query organizationOnboardingPlans($organizationId: ID!) {
  organizationPlansForOrganization(organizationId: $organizationId) {
    id
    name
    retired
    masterPlanId
    milestones {
      id
      name
      order
      dueDate
      optional
      items {
        status
        text
        updatedAt
      }
      retired
      statusDetails {
        updatedAt
        status
        text
      }
    }
    statusDetails {
      updatedAt
      status
      text
    }
  }
}
    `;

export const useOrganizationOnboardingPlansQuery = <
  TData = OrganizationOnboardingPlansQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: OrganizationOnboardingPlansQueryVariables,
  options?: Omit<
    UseQueryOptions<OrganizationOnboardingPlansQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<
      OrganizationOnboardingPlansQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<OrganizationOnboardingPlansQuery, TError, TData>({
    queryKey: ['organizationOnboardingPlans', variables],
    queryFn: fetcher<
      OrganizationOnboardingPlansQuery,
      OrganizationOnboardingPlansQueryVariables
    >(client, OrganizationOnboardingPlansDocument, variables, headers),
    ...options,
  });
};

useOrganizationOnboardingPlansQuery.document =
  OrganizationOnboardingPlansDocument;

useOrganizationOnboardingPlansQuery.getKey = (
  variables: OrganizationOnboardingPlansQueryVariables,
) => ['organizationOnboardingPlans', variables];

export const useInfiniteOrganizationOnboardingPlansQuery = <
  TData = InfiniteData<OrganizationOnboardingPlansQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: OrganizationOnboardingPlansQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<OrganizationOnboardingPlansQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      OrganizationOnboardingPlansQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<OrganizationOnboardingPlansQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? [
          'organizationOnboardingPlans.infinite',
          variables,
        ],
        queryFn: (metaData) =>
          fetcher<
            OrganizationOnboardingPlansQuery,
            OrganizationOnboardingPlansQueryVariables
          >(
            client,
            OrganizationOnboardingPlansDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteOrganizationOnboardingPlansQuery.getKey = (
  variables: OrganizationOnboardingPlansQueryVariables,
) => ['organizationOnboardingPlans.infinite', variables];

useOrganizationOnboardingPlansQuery.fetcher = (
  client: GraphQLClient,
  variables: OrganizationOnboardingPlansQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    OrganizationOnboardingPlansQuery,
    OrganizationOnboardingPlansQueryVariables
  >(client, OrganizationOnboardingPlansDocument, variables, headers);

useOrganizationOnboardingPlansQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables?: OrganizationOnboardingPlansQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: OrganizationOnboardingPlansQuery,
    ) => OrganizationOnboardingPlansQuery,
  ) => {
    const cacheKey = useOrganizationOnboardingPlansQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<OrganizationOnboardingPlansQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<OrganizationOnboardingPlansQuery>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
useInfiniteOrganizationOnboardingPlansQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables?: OrganizationOnboardingPlansQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: InfiniteData<OrganizationOnboardingPlansQuery>,
    ) => InfiniteData<OrganizationOnboardingPlansQuery>,
  ) => {
    const cacheKey =
      useInfiniteOrganizationOnboardingPlansQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<OrganizationOnboardingPlansQuery>>(
        cacheKey,
      );
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<OrganizationOnboardingPlansQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
