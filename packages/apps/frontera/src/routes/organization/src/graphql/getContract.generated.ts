// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, useInfiniteQuery, UseQueryOptions, UseInfiniteQueryOptions, InfiniteData } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type GetContractQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type GetContractQuery = { __typename?: 'Query', contract: { __typename?: 'Contract', id: string, contractUrl?: string | null, billingEnabled: boolean, organizationLegalName?: string | null, committedPeriodInMonths?: any | null, currency?: Types.Currency | null, contractName: string, contractEnded?: any | null, serviceStarted?: any | null, autoRenew: boolean, approved: boolean, contractStatus: Types.ContractStatus, metadata: { __typename?: 'Metadata', id: string }, attachments?: Array<{ __typename?: 'Attachment', id: string, basePath: string, fileName: string }> | null, billingDetails?: { __typename?: 'BillingDetails', billingCycle?: Types.ContractBillingCycle | null, addressLine1?: string | null, addressLine2?: string | null, locality?: string | null, region?: string | null, invoicingStarted?: any | null, country?: string | null, postalCode?: string | null, billingEmail?: string | null, invoiceNote?: string | null, canPayWithCard?: boolean | null, canPayWithDirectDebit?: boolean | null, canPayWithBankTransfer?: boolean | null, nextInvoicing?: any | null, payAutomatically?: boolean | null, payOnline?: boolean | null, dueDays?: any | null, billingEmailCC?: Array<string> | null, check?: boolean | null, billingEmailBCC?: Array<string> | null } | null } };



export const GetContractDocument = `
    query getContract($id: ID!) {
  contract(id: $id) {
    metadata {
      id
    }
    id
    contractUrl
    billingEnabled
    organizationLegalName
    committedPeriodInMonths
    currency
    contractName
    contractEnded
    serviceStarted
    autoRenew
    approved
    contractStatus
    attachments {
      id
      basePath
      fileName
    }
    billingDetails {
      billingCycle
      addressLine1
      addressLine2
      locality
      region
      invoicingStarted
      country
      postalCode
      billingEmail
      invoiceNote
      canPayWithCard
      canPayWithDirectDebit
      canPayWithBankTransfer
      nextInvoicing
      payAutomatically
      payOnline
      invoicingStarted
      region
      dueDays
      billingEmail
      billingEmailCC
      check
      billingEmailBCC
    }
  }
}
    `;

export const useGetContractQuery = <
      TData = GetContractQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetContractQueryVariables,
      options?: Omit<UseQueryOptions<GetContractQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetContractQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetContractQuery, TError, TData>(
      {
    queryKey: ['getContract', variables],
    queryFn: fetcher<GetContractQuery, GetContractQueryVariables>(client, GetContractDocument, variables, headers),
    ...options
  }
    )};

useGetContractQuery.document = GetContractDocument;

useGetContractQuery.getKey = (variables: GetContractQueryVariables) => ['getContract', variables];

export const useInfiniteGetContractQuery = <
      TData = InfiniteData<GetContractQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetContractQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetContractQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetContractQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetContractQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['getContract.infinite', variables],
      queryFn: (metaData) => fetcher<GetContractQuery, GetContractQueryVariables>(client, GetContractDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetContractQuery.getKey = (variables: GetContractQueryVariables) => ['getContract.infinite', variables];


useGetContractQuery.fetcher = (client: GraphQLClient, variables: GetContractQueryVariables, headers?: RequestInit['headers']) => fetcher<GetContractQuery, GetContractQueryVariables>(client, GetContractDocument, variables, headers);


useGetContractQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetContractQueryVariables) =>
  (mutator: (cacheEntry: GetContractQuery) => GetContractQuery) => {
    const cacheKey = useGetContractQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetContractQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetContractQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetContractQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetContractQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetContractQuery>) => InfiniteData<GetContractQuery>) => {
    const cacheKey = useInfiniteGetContractQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetContractQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetContractQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }