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
export type GetOrganizationOptionsQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;


export type GetOrganizationOptionsQuery = { __typename?: 'Query', dashboardView_Organizations?: { __typename?: 'OrganizationPage', content: Array<{ __typename?: 'Organization', value: string, label: string, subsidiaries: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', id: string } }>, subsidiaryOf: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', id: string } }> }> } | null };



export const GetOrganizationOptionsDocument = `
    query getOrganizationOptions($pagination: Pagination!, $where: Filter, $sort: SortBy) {
  dashboardView_Organizations(pagination: $pagination, where: $where, sort: $sort) {
    content {
      value: id
      label: name
      subsidiaries {
        organization {
          id
        }
      }
      subsidiaryOf {
        organization {
          id
        }
      }
    }
  }
}
    `;

export const useGetOrganizationOptionsQuery = <
      TData = GetOrganizationOptionsQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetOrganizationOptionsQueryVariables,
      options?: Omit<UseQueryOptions<GetOrganizationOptionsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetOrganizationOptionsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<GetOrganizationOptionsQuery, TError, TData>(
      {
    queryKey: ['getOrganizationOptions', variables],
    queryFn: fetcher<GetOrganizationOptionsQuery, GetOrganizationOptionsQueryVariables>(client, GetOrganizationOptionsDocument, variables, headers),
    ...options
  }
    )};

useGetOrganizationOptionsQuery.document = GetOrganizationOptionsDocument;

useGetOrganizationOptionsQuery.getKey = (variables: GetOrganizationOptionsQueryVariables) => ['getOrganizationOptions', variables];

export const useInfiniteGetOrganizationOptionsQuery = <
      TData = InfiniteData<GetOrganizationOptionsQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: GetOrganizationOptionsQueryVariables,
      options: Omit<UseInfiniteQueryOptions<GetOrganizationOptionsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<GetOrganizationOptionsQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<GetOrganizationOptionsQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? ['getOrganizationOptions.infinite', variables],
      queryFn: (metaData) => fetcher<GetOrganizationOptionsQuery, GetOrganizationOptionsQueryVariables>(client, GetOrganizationOptionsDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteGetOrganizationOptionsQuery.getKey = (variables: GetOrganizationOptionsQueryVariables) => ['getOrganizationOptions.infinite', variables];


useGetOrganizationOptionsQuery.fetcher = (client: GraphQLClient, variables: GetOrganizationOptionsQueryVariables, headers?: RequestInit['headers']) => fetcher<GetOrganizationOptionsQuery, GetOrganizationOptionsQueryVariables>(client, GetOrganizationOptionsDocument, variables, headers);


useGetOrganizationOptionsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetOrganizationOptionsQueryVariables) =>
  (mutator: (cacheEntry: GetOrganizationOptionsQuery) => GetOrganizationOptionsQuery) => {
    const cacheKey = useGetOrganizationOptionsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetOrganizationOptionsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetOrganizationOptionsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteGetOrganizationOptionsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetOrganizationOptionsQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<GetOrganizationOptionsQuery>) => InfiniteData<GetOrganizationOptionsQuery>) => {
    const cacheKey = useInfiniteGetOrganizationOptionsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetOrganizationOptionsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetOrganizationOptionsQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }