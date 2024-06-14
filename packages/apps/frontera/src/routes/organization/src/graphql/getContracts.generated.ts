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
      contractName: string;
      serviceStarted?: any | null;
      contractSigned?: any | null;
      contractEnded?: any | null;
      contractStatus: Types.ContractStatus;
      committedPeriodInMonths?: any | null;
      approved: boolean;
      contractUrl?: string | null;
      billingCycle?: Types.ContractBillingCycle | null;
      billingEnabled: boolean;
      currency?: Types.Currency | null;
      invoiceEmail?: string | null;
      autoRenew: boolean;
      metadata: {
        __typename?: 'Metadata';
        id: string;
        created: any;
        source: Types.DataSource;
        lastUpdated: any;
      };
      billingDetails?: {
        __typename?: 'BillingDetails';
        nextInvoicing?: any | null;
        postalCode?: string | null;
        country?: string | null;
        locality?: string | null;
        addressLine1?: string | null;
        addressLine2?: string | null;
        invoiceNote?: string | null;
        organizationLegalName?: string | null;
        billingCycle?: Types.ContractBillingCycle | null;
        invoicingStarted?: any | null;
        region?: string | null;
        dueDays?: any | null;
        billingEmail?: string | null;
        billingEmailCC?: Array<string> | null;
        billingEmailBCC?: Array<string> | null;
      } | null;
      upcomingInvoices: Array<{
        __typename?: 'Invoice';
        invoicePeriodEnd: any;
        invoicePeriodStart: any;
        status?: Types.InvoiceStatus | null;
        issued: any;
        amountDue: number;
        due: any;
        currency: string;
        invoiceNumber: string;
        invoiceUrl: string;
        subtotal: number;
        taxDue: number;
        note?: string | null;
        metadata: { __typename?: 'Metadata'; id: string };
        invoiceLineItems: Array<{
          __typename?: 'InvoiceLine';
          quantity: any;
          subtotal: number;
          taxDue: number;
          total: number;
          price: number;
          description: string;
          metadata: { __typename?: 'Metadata'; id: string; created: any };
        }>;
        contract: {
          __typename?: 'Contract';
          billingDetails?: {
            __typename?: 'BillingDetails';
            canPayWithBankTransfer?: boolean | null;
          } | null;
        };
        customer: {
          __typename?: 'InvoiceCustomer';
          name?: string | null;
          email?: string | null;
          addressLine1?: string | null;
          addressLine2?: string | null;
          addressZip?: string | null;
          addressLocality?: string | null;
          addressCountry?: string | null;
          addressRegion?: string | null;
        };
        provider: {
          __typename?: 'InvoiceProvider';
          name?: string | null;
          addressLine1?: string | null;
          addressLine2?: string | null;
          addressZip?: string | null;
          addressLocality?: string | null;
          addressCountry?: string | null;
        };
      }>;
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
        renewalAdjustedRate: any;
        renewalUpdatedByUserId: string;
        renewedAt?: any | null;
        updatedAt?: any | null;
        owner?: {
          __typename?: 'User';
          id: string;
          firstName: string;
          lastName: string;
          name?: string | null;
        } | null;
      }> | null;
      contractLineItems?: Array<{
        __typename?: 'ServiceLineItem';
        description: string;
        billingCycle: Types.BilledType;
        price: number;
        quantity: any;
        comments: string;
        serviceEnded?: any | null;
        parentId: string;
        serviceStarted: any;
        metadata: {
          __typename?: 'Metadata';
          id: string;
          created: any;
          lastUpdated: any;
          source: Types.DataSource;
          appSource: string;
          sourceOfTruth: Types.DataSource;
        };
        tax: {
          __typename?: 'Tax';
          salesTax: boolean;
          vat: boolean;
          taxRate: number;
        };
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
      metadata {
        id
        created
        source
        lastUpdated
      }
      contractName
      serviceStarted
      contractSigned
      contractEnded
      contractStatus
      committedPeriodInMonths
      approved
      contractUrl
      billingCycle
      billingEnabled
      currency
      invoiceEmail
      autoRenew
      billingDetails {
        nextInvoicing
        postalCode
        country
        locality
        addressLine1
        addressLine2
        invoiceNote
        organizationLegalName
        billingCycle
        invoicingStarted
        region
        dueDays
        billingEmail
        billingEmailCC
        billingEmailBCC
      }
      upcomingInvoices {
        metadata {
          id
        }
        invoicePeriodEnd
        invoicePeriodStart
        status
        issued
        amountDue
        due
        currency
        invoiceLineItems {
          metadata {
            id
            created
          }
          quantity
          subtotal
          taxDue
          total
          price
          description
        }
        contract {
          billingDetails {
            canPayWithBankTransfer
          }
        }
        status
        invoiceNumber
        invoicePeriodStart
        invoicePeriodEnd
        invoiceUrl
        due
        issued
        subtotal
        taxDue
        currency
        note
        customer {
          name
          email
          addressLine1
          addressLine2
          addressZip
          addressLocality
          addressCountry
          addressRegion
        }
        provider {
          name
          addressLine1
          addressLine2
          addressZip
          addressLocality
          addressCountry
        }
      }
      opportunities {
        id
        comments
        internalStage
        internalType
        amount
        maxAmount
        name
        renewalLikelihood
        renewalAdjustedRate
        renewalUpdatedByUserId
        renewedAt
        updatedAt
        owner {
          id
          firstName
          lastName
          name
        }
      }
      contractLineItems {
        metadata {
          id
          created
          lastUpdated
          source
          appSource
          sourceOfTruth
        }
        description
        billingCycle
        price
        quantity
        comments
        serviceEnded
        parentId
        serviceStarted
        tax {
          salesTax
          vat
          taxRate
        }
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
