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
export type TenantBillingProfileQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type TenantBillingProfileQuery = { __typename?: 'Query', tenantBillingProfile: { __typename?: 'TenantBillingProfile', id: string, addressLine1: string, addressLine2: string, addressLine3: string, locality: string, zip: string, legalName: string, country: string, vatNumber: string, sendInvoicesFrom: string, canPayWithBankTransfer: boolean, canPayWithPigeon: boolean, check: boolean, region: string } };



export const TenantBillingProfileDocument = `
    query TenantBillingProfile($id: ID!) {
  tenantBillingProfile(id: $id) {
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
    canPayWithBankTransfer
    canPayWithPigeon
    check
    region
  }
}
    `;

export const useTenantBillingProfileQuery = <
      TData = TenantBillingProfileQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: TenantBillingProfileQueryVariables,
      options?: Omit<UseQueryOptions<TenantBillingProfileQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<TenantBillingProfileQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<TenantBillingProfileQuery, TError, TData>(
      {
    queryKey: ['TenantBillingProfile', variables],
    queryFn: fetcher<TenantBillingProfileQuery, TenantBillingProfileQueryVariables>(client, TenantBillingProfileDocument, variables, headers),
    ...options
  }
    )};

useTenantBillingProfileQuery.document = TenantBillingProfileDocument;

useTenantBillingProfileQuery.getKey = (variables: TenantBillingProfileQueryVariables) => ['TenantBillingProfile', variables];

export const useInfiniteTenantBillingProfileQuery = <
      TData = InfiniteData<TenantBillingProfileQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: TenantBillingProfileQueryVariables,
      options: Omit<UseInfiniteQueryOptions<TenantBillingProfileQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<TenantBillingProfileQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<TenantBillingProfileQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['TenantBillingProfile.infinite', variables],
      queryFn: (metaData) => fetcher<TenantBillingProfileQuery, TenantBillingProfileQueryVariables>(client, TenantBillingProfileDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteTenantBillingProfileQuery.getKey = (variables: TenantBillingProfileQueryVariables) => ['TenantBillingProfile.infinite', variables];


useTenantBillingProfileQuery.fetcher = (client: GraphQLClient, variables: TenantBillingProfileQueryVariables, headers?: RequestInit['headers']) => fetcher<TenantBillingProfileQuery, TenantBillingProfileQueryVariables>(client, TenantBillingProfileDocument, variables, headers);


useTenantBillingProfileQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantBillingProfileQueryVariables) =>
  (mutator: (cacheEntry: TenantBillingProfileQuery) => TenantBillingProfileQuery) => {
    const cacheKey = useTenantBillingProfileQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<TenantBillingProfileQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<TenantBillingProfileQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteTenantBillingProfileQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: TenantBillingProfileQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<TenantBillingProfileQuery>) => InfiniteData<TenantBillingProfileQuery>) => {
    const cacheKey = useInfiniteTenantBillingProfileQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<TenantBillingProfileQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<TenantBillingProfileQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }