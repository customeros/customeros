// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../types/__generated__/graphql.types';

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
export type OrganizationQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
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
    socials: Array<{ __typename?: 'Social'; id: string; url: string }>;
    relationshipStages: Array<{
      __typename?: 'OrganizationRelationshipStage';
      relationship: Types.OrganizationRelationship;
      stage?: string | null;
    }>;
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
    socials {
      id
      url
    }
    relationshipStages {
      relationship
      stage
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
  options?: UseQueryOptions<OrganizationQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<OrganizationQuery, TError, TData>(
    ['Organization', variables],
    fetcher<OrganizationQuery, OrganizationQueryVariables>(
      client,
      OrganizationDocument,
      variables,
      headers,
    ),
    options,
  );
useOrganizationQuery.document = OrganizationDocument;

useOrganizationQuery.getKey = (variables: OrganizationQueryVariables) => [
  'Organization',
  variables,
];
export const useInfiniteOrganizationQuery = <
  TData = OrganizationQuery,
  TError = unknown,
>(
  pageParamKey: keyof OrganizationQueryVariables,
  client: GraphQLClient,
  variables: OrganizationQueryVariables,
  options?: UseInfiniteQueryOptions<OrganizationQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<OrganizationQuery, TError, TData>(
    ['Organization.infinite', variables],
    (metaData) =>
      fetcher<OrganizationQuery, OrganizationQueryVariables>(
        client,
        OrganizationDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

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
