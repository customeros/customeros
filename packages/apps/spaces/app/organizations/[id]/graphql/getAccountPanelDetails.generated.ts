// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

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
export type OrganizationAccountDetailsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type OrganizationAccountDetailsQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    id: string;
    name: string;
    accountDetails?: {
      __typename?: 'OrgAccountDetails';
      renewalForecast?: {
        __typename?: 'RenewalForecast';
        amount?: number | null;
        potentialAmount?: number | null;
        comment?: string | null;
        updatedAt?: any | null;
        updatedBy?: {
          __typename?: 'User';
          id: string;
          firstName: string;
          lastName: string;
          emails?: Array<{
            __typename?: 'Email';
            email?: string | null;
          }> | null;
        } | null;
      } | null;
      renewalLikelihood?: {
        __typename?: 'RenewalLikelihood';
        probability?: Types.RenewalLikelihoodProbability | null;
        previousProbability?: Types.RenewalLikelihoodProbability | null;
        comment?: string | null;
        updatedAt?: any | null;
        updatedBy?: {
          __typename?: 'User';
          id: string;
          firstName: string;
          lastName: string;
          emails?: Array<{
            __typename?: 'Email';
            email?: string | null;
          }> | null;
        } | null;
      } | null;
      billingDetails?: {
        __typename?: 'BillingDetails';
        renewalCycle?: Types.RenewalCycle | null;
        frequency?: Types.RenewalCycle | null;
        amount?: number | null;
        renewalCycleStart?: any | null;
        renewalCycleNext?: any | null;
      } | null;
    } | null;
  } | null;
};

export const OrganizationAccountDetailsDocument = `
    query OrganizationAccountDetails($id: ID!) {
  organization(id: $id) {
    id
    name
    accountDetails {
      renewalForecast {
        amount
        potentialAmount
        comment
        updatedAt
        updatedBy {
          id
          firstName
          lastName
          emails {
            email
          }
        }
      }
      renewalLikelihood {
        probability
        previousProbability
        comment
        updatedBy {
          id
          firstName
          lastName
          emails {
            email
          }
        }
        updatedAt
      }
      billingDetails {
        renewalCycle
        frequency
        amount
        renewalCycleStart
        renewalCycleNext
      }
    }
  }
}
    `;
export const useOrganizationAccountDetailsQuery = <
  TData = OrganizationAccountDetailsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: OrganizationAccountDetailsQueryVariables,
  options?: UseQueryOptions<OrganizationAccountDetailsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<OrganizationAccountDetailsQuery, TError, TData>(
    ['OrganizationAccountDetails', variables],
    fetcher<
      OrganizationAccountDetailsQuery,
      OrganizationAccountDetailsQueryVariables
    >(client, OrganizationAccountDetailsDocument, variables, headers),
    options,
  );
useOrganizationAccountDetailsQuery.document =
  OrganizationAccountDetailsDocument;

useOrganizationAccountDetailsQuery.getKey = (
  variables: OrganizationAccountDetailsQueryVariables,
) => ['OrganizationAccountDetails', variables];
export const useInfiniteOrganizationAccountDetailsQuery = <
  TData = OrganizationAccountDetailsQuery,
  TError = unknown,
>(
  pageParamKey: keyof OrganizationAccountDetailsQueryVariables,
  client: GraphQLClient,
  variables: OrganizationAccountDetailsQueryVariables,
  options?: UseInfiniteQueryOptions<
    OrganizationAccountDetailsQuery,
    TError,
    TData
  >,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<OrganizationAccountDetailsQuery, TError, TData>(
    ['OrganizationAccountDetails.infinite', variables],
    (metaData) =>
      fetcher<
        OrganizationAccountDetailsQuery,
        OrganizationAccountDetailsQueryVariables
      >(
        client,
        OrganizationAccountDetailsDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteOrganizationAccountDetailsQuery.getKey = (
  variables: OrganizationAccountDetailsQueryVariables,
) => ['OrganizationAccountDetails.infinite', variables];
useOrganizationAccountDetailsQuery.fetcher = (
  client: GraphQLClient,
  variables: OrganizationAccountDetailsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    OrganizationAccountDetailsQuery,
    OrganizationAccountDetailsQueryVariables
  >(client, OrganizationAccountDetailsDocument, variables, headers);
