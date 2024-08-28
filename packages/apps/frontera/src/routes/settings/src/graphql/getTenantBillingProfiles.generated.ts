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
export type TenantBillingProfilesQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type TenantBillingProfilesQuery = { __typename?: 'Query', tenantBillingProfiles: Array<{ __typename?: 'TenantBillingProfile', id: string, addressLine1: string, addressLine2: string, addressLine3: string, locality: string, zip: string, legalName: string, country: string, vatNumber: string, sendInvoicesFrom: string, sendInvoicesBcc: string, canPayWithPigeon: boolean, canPayWithBankTransfer: boolean, check: boolean, region: string }> };



export const TenantBillingProfilesDocument = `
    query TenantBillingProfiles {
  tenantBillingProfiles {
    id
    addressLine1
    addressLine2
    addressLine3
    locality
    zip
    legalName
    country
    vatNumber
    sendInvoicesFrom
    sendInvoicesBcc
    canPayWithPigeon
    canPayWithBankTransfer
    check
    region
  }
}
    `;

export const useTenantBillingProfilesQuery = <
      TData = TenantBillingProfilesQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: TenantBillingProfilesQueryVariables,
      options?: Omit<UseQueryOptions<TenantBillingProfilesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<TenantBillingProfilesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<TenantBillingProfilesQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['TenantBillingProfiles'] : ['TenantBillingProfiles', variables],
    queryFn: fetcher<TenantBillingProfilesQuery, TenantBillingProfilesQueryVariables>(client, TenantBillingProfilesDocument, variables, headers),
    ...options
  }
    )};

useTenantBillingProfilesQuery.document = TenantBillingProfilesDocument;

useTenantBillingProfilesQuery.getKey = (variables?: TenantBillingProfilesQueryVariables) => variables === undefined ? ['TenantBillingProfiles'] : ['TenantBillingProfiles', variables];

export const useInfiniteTenantBillingProfilesQuery = <
      TData = InfiniteData<TenantBillingProfilesQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: TenantBillingProfilesQueryVariables,
      options: Omit<UseInfiniteQueryOptions<TenantBillingProfilesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<TenantBillingProfilesQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<TenantBillingProfilesQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['TenantBillingProfiles.infinite'] : ['TenantBillingProfiles.infinite', variables],
      queryFn: (metaData) => fetcher<TenantBillingProfilesQuery, TenantBillingProfilesQueryVariables>(client, TenantBillingProfilesDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteTenantBillingProfilesQuery.getKey = (variables?: TenantBillingProfilesQueryVariables) => variables === undefined ? ['TenantBillingProfiles.infinite'] : ['TenantBillingProfiles.infinite', variables];


useTenantBillingProfilesQuery.fetcher = (client: GraphQLClient, variables?: TenantBillingProfilesQueryVariables, headers?: RequestInit['headers']) => fetcher<TenantBillingProfilesQuery, TenantBillingProfilesQueryVariables>(client, TenantBillingProfilesDocument, variables, headers);


useTenantBillingProfilesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantBillingProfilesQueryVariables) =>
  (mutator: (cacheEntry: TenantBillingProfilesQuery) => TenantBillingProfilesQuery) => {
    const cacheKey = useTenantBillingProfilesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TenantBillingProfilesQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TenantBillingProfilesQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteTenantBillingProfilesQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantBillingProfilesQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<TenantBillingProfilesQuery>) => InfiniteData<TenantBillingProfilesQuery>) => {
    const cacheKey = useInfiniteTenantBillingProfilesQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TenantBillingProfilesQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TenantBillingProfilesQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }