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
export type GetInvoicesCountQueryVariables = Types.Exact<{
  organizationId?: Types.InputMaybe<Types.Scalars['ID']['input']>;
}>;


export type GetInvoicesCountQuery = { __typename?: 'Query', invoices: { __typename?: 'InvoicesPage', totalElements: any } };



export const GetInvoicesCountDocument = `
    query getInvoicesCount($organizationId: ID) {
  invoices(
    pagination: {page: 0, limit: 0}
    organizationId: $organizationId
    where: {AND: [{filter: {property: "DRY_RUN", operation: EQ, value: false}}]}
  ) {
    totalElements
  }
}
    `;

export const useGetInvoicesCountQuery = <
      TData = GetInvoicesCountQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GetInvoicesCountQueryVariables,
      options?: Omit<UseQueryOptions<GetInvoicesCountQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetInvoicesCountQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetInvoicesCountQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['getInvoicesCount'] : ['getInvoicesCount', variables],
    queryFn: fetcher<GetInvoicesCountQuery, GetInvoicesCountQueryVariables>(client, GetInvoicesCountDocument, variables, headers),
    ...options
  }
    )};

useGetInvoicesCountQuery.document = GetInvoicesCountDocument;

useGetInvoicesCountQuery.getKey = (variables?: GetInvoicesCountQueryVariables) => variables === undefined ? ['getInvoicesCount'] : ['getInvoicesCount', variables];

export const useInfiniteGetInvoicesCountQuery = <
      TData = InfiniteData<GetInvoicesCountQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetInvoicesCountQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetInvoicesCountQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetInvoicesCountQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetInvoicesCountQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['getInvoicesCount.infinite'] : ['getInvoicesCount.infinite', variables],
      queryFn: (metaData) => fetcher<GetInvoicesCountQuery, GetInvoicesCountQueryVariables>(client, GetInvoicesCountDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetInvoicesCountQuery.getKey = (variables?: GetInvoicesCountQueryVariables) => variables === undefined ? ['getInvoicesCount.infinite'] : ['getInvoicesCount.infinite', variables];


useGetInvoicesCountQuery.fetcher = (client: GraphQLClient, variables?: GetInvoicesCountQueryVariables, headers?: RequestInit['headers']) => fetcher<GetInvoicesCountQuery, GetInvoicesCountQueryVariables>(client, GetInvoicesCountDocument, variables, headers);


useGetInvoicesCountQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoicesCountQueryVariables) =>
  (mutator: (cacheEntry: GetInvoicesCountQuery) => GetInvoicesCountQuery) => {
    const cacheKey = useGetInvoicesCountQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetInvoicesCountQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetInvoicesCountQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetInvoicesCountQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetInvoicesCountQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetInvoicesCountQuery>) => InfiniteData<GetInvoicesCountQuery>) => {
    const cacheKey = useInfiniteGetInvoicesCountQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetInvoicesCountQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetInvoicesCountQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }