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
export type GetContractsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
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
      billingCycle?: Types.ContractBillingCycle | null;
      invoicingStartDate?: any | null;
      currency?: Types.Currency | null;
      organizationLegalName?: string | null;
      addressLine1?: string | null;
      addressLine2?: string | null;
      locality?: string | null;
      country?: string | null;
      zip?: string | null;
      invoiceEmail?: string | null;
      invoiceNote?: string | null;
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
        serviceStarted: any;
        vatRate: number;
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
      billingCycle
      invoicingStartDate
      currency
      organizationLegalName
      addressLine1
      addressLine2
      locality
      country
      zip
      invoiceEmail
      invoiceNote
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
        serviceStarted
        vatRate
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
  options?: Omit<
    UseQueryOptions<GetContractsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<GetContractsQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetContractsQuery, TError, TData>({
    queryKey: ['getContracts', variables],
    queryFn: fetcher<GetContractsQuery, GetContractsQueryVariables>(
      client,
      GetContractsDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGetContractsQuery.document = GetContractsDocument;

useGetContractsQuery.getKey = (variables: GetContractsQueryVariables) => [
  'getContracts',
  variables,
];

export const useInfiniteGetContractsQuery = <
  TData = InfiniteData<GetContractsQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetContractsQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetContractsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetContractsQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetContractsQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? ['getContracts.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GetContractsQuery, GetContractsQueryVariables>(
            client,
            GetContractsDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

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
  (queryClient: QueryClient, variables?: GetContractsQueryVariables) =>
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
  (queryClient: QueryClient, variables?: GetContractsQueryVariables) =>
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
