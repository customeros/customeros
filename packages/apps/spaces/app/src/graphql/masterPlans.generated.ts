// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../types/__generated__/graphql.types';

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
export type MasterPlansQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type MasterPlansQuery = { __typename?: 'Query', masterPlans: Array<{ __typename?: 'MasterPlan', id: string, name: string, retired: boolean, milestones: Array<{ __typename?: 'MasterPlanMilestone', id: string, name: string, order: any, durationHours: any, optional: boolean, items: Array<string>, retired: boolean }>, retiredMilestones: Array<{ __typename?: 'MasterPlanMilestone', id: string, name: string, order: any, durationHours: any, optional: boolean, items: Array<string>, retired: boolean }> }> };



export const MasterPlansDocument = `
    query masterPlans {
  masterPlans {
    id
    name
    retired
    milestones {
      id
      name
      order
      durationHours
      optional
      items
      retired
    }
    retiredMilestones {
      id
      name
      order
      durationHours
      optional
      items
      retired
    }
  }
}
    `;

export const useMasterPlansQuery = <
      TData = MasterPlansQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: MasterPlansQueryVariables,
      options?: Omit<UseQueryOptions<MasterPlansQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<MasterPlansQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useQuery<MasterPlansQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['masterPlans'] : ['masterPlans', variables],
    queryFn: fetcher<MasterPlansQuery, MasterPlansQueryVariables>(client, MasterPlansDocument, variables, headers),
    ...options
  }
    )};

useMasterPlansQuery.document = MasterPlansDocument;

useMasterPlansQuery.getKey = (variables?: MasterPlansQueryVariables) => variables === undefined ? ['masterPlans'] : ['masterPlans', variables];

export const useInfiniteMasterPlansQuery = <
      TData = InfiniteData<MasterPlansQuery>,
      TError = unknown
    >(
      client: GraphQLClient,
      variables: MasterPlansQueryVariables,
      options: Omit<UseInfiniteQueryOptions<MasterPlansQuery, TError, TData>, 'queryKey'> & { queryKey?: UseInfiniteQueryOptions<MasterPlansQuery, TError, TData>['queryKey'] },
      headers?: RequestInit['headers']
    ) => {
    
    return useInfiniteQuery<MasterPlansQuery, TError, TData>(
      (() => {
    const { queryKey: optionsQueryKey, ...restOptions } = options;
    return {
      queryKey: optionsQueryKey ?? variables === undefined ? ['masterPlans.infinite'] : ['masterPlans.infinite', variables],
      queryFn: (metaData) => fetcher<MasterPlansQuery, MasterPlansQueryVariables>(client, MasterPlansDocument, {...variables, ...(metaData.pageParam ?? {})}, headers)(),
      ...restOptions
    }
  })()
    )};

useInfiniteMasterPlansQuery.getKey = (variables?: MasterPlansQueryVariables) => variables === undefined ? ['masterPlans.infinite'] : ['masterPlans.infinite', variables];


useMasterPlansQuery.fetcher = (client: GraphQLClient, variables?: MasterPlansQueryVariables, headers?: RequestInit['headers']) => fetcher<MasterPlansQuery, MasterPlansQueryVariables>(client, MasterPlansDocument, variables, headers);


useMasterPlansQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: MasterPlansQueryVariables) =>
  (mutator: (cacheEntry: MasterPlansQuery) => MasterPlansQuery) => {
    const cacheKey = useMasterPlansQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<MasterPlansQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<MasterPlansQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  }
useInfiniteMasterPlansQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: MasterPlansQueryVariables) =>
  (mutator: (cacheEntry: InfiniteData<MasterPlansQuery>) => InfiniteData<MasterPlansQuery>) => {
    const cacheKey = useInfiniteMasterPlansQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<MasterPlansQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<MasterPlansQuery>>(cacheKey, mutator);
    }
    return { previousEntries };
  }