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
export type OrganizationAccountDetailsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type OrganizationAccountDetailsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, name: string, note?: string | null } | null };



export const OrganizationAccountDetailsDocument = `
    query OrganizationAccountDetails($id: ID!) {
  organization(id: $id) {
    id
    name
    note
  }
}
    `;

export const useOrganizationAccountDetailsQuery = <
      TData = OrganizationAccountDetailsQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: OrganizationAccountDetailsQueryVariables,
      options?: Omit<UseQueryOptions<OrganizationAccountDetailsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<OrganizationAccountDetailsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<OrganizationAccountDetailsQuery, TError, TData>(
      {
    queryKey: ['OrganizationAccountDetails', variables],
    queryFn: fetcher<OrganizationAccountDetailsQuery, OrganizationAccountDetailsQueryVariables>(client, OrganizationAccountDetailsDocument, variables, headers),
    ...options
  }
    )};

useOrganizationAccountDetailsQuery.document = OrganizationAccountDetailsDocument;

useOrganizationAccountDetailsQuery.getKey = (variables: OrganizationAccountDetailsQueryVariables) => ['OrganizationAccountDetails', variables];

export const useInfiniteOrganizationAccountDetailsQuery = <
      TData = InfiniteData<OrganizationAccountDetailsQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: OrganizationAccountDetailsQueryVariables,
      options: Omit<UseInfiniteQueryOptions<OrganizationAccountDetailsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<OrganizationAccountDetailsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<OrganizationAccountDetailsQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['OrganizationAccountDetails.infinite', variables],
      queryFn: (metaData) => fetcher<OrganizationAccountDetailsQuery, OrganizationAccountDetailsQueryVariables>(client, OrganizationAccountDetailsDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteOrganizationAccountDetailsQuery.getKey = (variables: OrganizationAccountDetailsQueryVariables) => ['OrganizationAccountDetails.infinite', variables];


useOrganizationAccountDetailsQuery.fetcher = (client: GraphQLClient, variables: OrganizationAccountDetailsQueryVariables, headers?: RequestInit['headers']) => fetcher<OrganizationAccountDetailsQuery, OrganizationAccountDetailsQueryVariables>(client, OrganizationAccountDetailsDocument, variables, headers);


useOrganizationAccountDetailsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: OrganizationAccountDetailsQueryVariables) =>
  (mutator: (cacheEntry: OrganizationAccountDetailsQuery) => OrganizationAccountDetailsQuery) => {
    const cacheKey = useOrganizationAccountDetailsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<OrganizationAccountDetailsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<OrganizationAccountDetailsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteOrganizationAccountDetailsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: OrganizationAccountDetailsQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<OrganizationAccountDetailsQuery>) => InfiniteData<OrganizationAccountDetailsQuery>) => {
    const cacheKey = useInfiniteOrganizationAccountDetailsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<OrganizationAccountDetailsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<OrganizationAccountDetailsQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }