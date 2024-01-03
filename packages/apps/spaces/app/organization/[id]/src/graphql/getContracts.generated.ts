// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

import type { InfiniteData } from '@tanstack/react-query';
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
export type GetContractsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type GetContractsQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    id: string;
    name: string;
    note?: string | null;
    accountDetails?: {
      __typename?: 'OrgAccountDetails';
      renewalSummary?: {
        __typename?: 'RenewalSummary';
        arrForecast?: number | null;
        maxArrForecast?: number | null;
        renewalLikelihood?: Types.OpportunityRenewalLikelihood | null;
      } | null;
    } | null;
    contracts?: Array<{
      __typename?: 'Contract';
      id: string;
      name: string;
      createdAt: any;
      serviceStartedAt?: any | null;
      signedAt?: any | null;
      endedAt?: any | null;
      renewalCycle: Types.ContractRenewalCycle;
      renewalPeriods?: any | null;
      status: Types.ContractStatus;
      contractUrl?: string | null;
      opportunities?: Array<{
        __typename?: 'Opportunity';
        id: string;
        comments: string;
        internalStage: Types.InternalStage;
        internalType: Types.InternalType;
        amount: number;
        maxAmount: number;
        name: string;
        renewalLikelihood: Types.OpportunityRenewalLikelihood;
        renewalUpdatedByUserId: string;
        renewalUpdatedByUserAt: any;
        renewedAt: any;
        owner?: {
          __typename?: 'User';
          id: string;
          firstName: string;
          lastName: string;
          name?: string | null;
        } | null;
      }> | null;
      serviceLineItems?: Array<{
        __typename?: 'ServiceLineItem';
        id: string;
        createdAt: any;
        updatedAt: any;
        name: string;
        billed: Types.BilledType;
        price: number;
        quantity: any;
        source: Types.DataSource;
        comments: string;
        sourceOfTruth: Types.DataSource;
        appSource: string;
        endedAt?: any | null;
        parentId: string;
      }> | null;
    }> | null;
  } | null;
};

export const GetContractsDocument = `
    query getContracts($id: ID!) {
  organization(id: $id) {
    id
    name
    note
    accountDetails {
      renewalSummary {
        arrForecast
        maxArrForecast
        renewalLikelihood
      }
    }
    contracts {
      id
      name
      createdAt
      serviceStartedAt
      signedAt
      endedAt
      renewalCycle
      renewalPeriods
      status
      contractUrl
      opportunities {
        id
        comments
        internalStage
        internalType
        amount
        maxAmount
        name
        renewalLikelihood
        renewalUpdatedByUserId
        renewalUpdatedByUserAt
        renewedAt
        owner {
          id
          firstName
          lastName
          name
        }
      }
      serviceLineItems {
        id
        createdAt
        updatedAt
        name
        billed
        price
        quantity
        source
        comments
        sourceOfTruth
        appSource
        endedAt
        parentId
      }
    }
  }
}
    `;
export const useGetContractsQuery = <
  TData = GetContractsQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetContractsQueryVariables,
  options?: UseQueryOptions<GetContractsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useQuery<GetContractsQuery, TError, TData>(
    ['getContracts', variables],
    fetcher<GetContractsQuery, GetContractsQueryVariables>(
      client,
      GetContractsDocument,
      variables,
      headers,
    ),
    options,
  );
useGetContractsQuery.document = GetContractsDocument;

useGetContractsQuery.getKey = (variables: GetContractsQueryVariables) => [
  'getContracts',
  variables,
];
export const useInfiniteGetContractsQuery = <
  TData = GetContractsQuery,
  TError = unknown,
>(
  pageParamKey: keyof GetContractsQueryVariables,
  client: GraphQLClient,
  variables: GetContractsQueryVariables,
  options?: UseInfiniteQueryOptions<GetContractsQuery, TError, TData>,
  headers?: RequestInit['headers'],
) =>
  useInfiniteQuery<GetContractsQuery, TError, TData>(
    ['getContracts.infinite', variables],
    (metaData) =>
      fetcher<GetContractsQuery, GetContractsQueryVariables>(
        client,
        GetContractsDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
        headers,
      )(),
    options,
  );

useInfiniteGetContractsQuery.getKey = (
  variables: GetContractsQueryVariables,
) => ['getContracts.infinite', variables];
useGetContractsQuery.fetcher = (
  client: GraphQLClient,
  variables: GetContractsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetContractsQuery, GetContractsQueryVariables>(
    client,
    GetContractsDocument,
    variables,
    headers,
  );

useGetContractsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetContractsQueryVariables) =>
  (mutator: (cacheEntry: GetContractsQuery) => GetContractsQuery) => {
    const cacheKey = useGetContractsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetContractsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetContractsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetContractsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables: GetContractsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetContractsQuery>,
    ) => InfiniteData<GetContractsQuery>,
  ) => {
    const cacheKey = useInfiniteGetContractsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetContractsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetContractsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
