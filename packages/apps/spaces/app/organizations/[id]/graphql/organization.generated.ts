// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';

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
    contacts: {
      __typename?: 'ContactsPage';
      totalElements: any;
      content: Array<{
        __typename?: 'Contact';
        id: string;
        name?: string | null;
        firstName?: string | null;
        lastName?: string | null;
        prefix?: string | null;
        description?: string | null;
        tags?: Array<{ __typename?: 'Tag'; id: string; name: string }> | null;
        jobRoles: Array<{
          __typename?: 'JobRole';
          id: string;
          primary: boolean;
          jobTitle?: string | null;
        }>;
        phoneNumbers: Array<{
          __typename?: 'PhoneNumber';
          id: string;
          e164?: string | null;
          rawPhoneNumber?: string | null;
          label?: Types.PhoneNumberLabel | null;
          primary: boolean;
        }>;
        emails: Array<{
          __typename?: 'Email';
          id: string;
          email?: string | null;
        }>;
      }>;
    };
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
    contacts(pagination: {page: 0, limit: 100}) {
      content {
        id
        name
        firstName
        lastName
        prefix
        description
        tags {
          id
          name
        }
        jobRoles {
          id
          primary
          jobTitle
        }
        phoneNumbers {
          id
          e164
          rawPhoneNumber
          label
          primary
        }
        emails {
          id
          email
        }
      }
      totalElements
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
