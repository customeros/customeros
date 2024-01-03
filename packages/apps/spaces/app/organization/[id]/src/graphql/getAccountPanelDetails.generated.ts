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
export type OrganizationAccountDetailsQueryVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type OrganizationAccountDetailsQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    id: string;
    name: string;
    note?: string | null;
  } | null;
};

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

useOrganizationAccountDetailsQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables: OrganizationAccountDetailsQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: OrganizationAccountDetailsQuery,
    ) => OrganizationAccountDetailsQuery,
  ) => {
    const cacheKey = useOrganizationAccountDetailsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<OrganizationAccountDetailsQuery>(cacheKey);
    if (previousEntry) {
      queryClient.setQueryData<OrganizationAccountDetailsQuery>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
useInfiniteOrganizationAccountDetailsQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables: OrganizationAccountDetailsQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: InfiniteData<OrganizationAccountDetailsQuery>,
    ) => InfiniteData<OrganizationAccountDetailsQuery>,
  ) => {
    const cacheKey =
      useInfiniteOrganizationAccountDetailsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<OrganizationAccountDetailsQuery>>(
        cacheKey,
      );
    if (previousEntry) {
      queryClient.setQueryData<InfiniteData<OrganizationAccountDetailsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
