// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useMutation, UseMutationOptions } from '@tanstack/react-query';

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
export type DuplicateMilestoneMutationVariables = Types.Exact<{
  masterPlanId: Types.Scalars['ID']['input'];
  id: Types.Scalars['ID']['input'];
}>;

export type DuplicateMilestoneMutation = {
  __typename?: 'Mutation';
  masterPlanMilestone_Duplicate: {
    __typename?: 'MasterPlanMilestone';
    id: string;
    name: string;
    order: any;
    durationHours: any;
    optional: boolean;
    items: Array<string>;
    retired: boolean;
  };
};

export const DuplicateMilestoneDocument = `
    mutation duplicateMilestone($masterPlanId: ID!, $id: ID!) {
  masterPlanMilestone_Duplicate(masterPlanId: $masterPlanId, id: $id) {
    id
    name
    order
    durationHours
    optional
    items
    retired
  }
}
    `;

export const useDuplicateMilestoneMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    DuplicateMilestoneMutation,
    TError,
    DuplicateMilestoneMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    DuplicateMilestoneMutation,
    TError,
    DuplicateMilestoneMutationVariables,
    TContext
  >({
    mutationKey: ['duplicateMilestone'],
    mutationFn: (variables?: DuplicateMilestoneMutationVariables) =>
      fetcher<DuplicateMilestoneMutation, DuplicateMilestoneMutationVariables>(
        client,
        DuplicateMilestoneDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useDuplicateMilestoneMutation.getKey = () => ['duplicateMilestone'];

useDuplicateMilestoneMutation.fetcher = (
  client: GraphQLClient,
  variables: DuplicateMilestoneMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<DuplicateMilestoneMutation, DuplicateMilestoneMutationVariables>(
    client,
    DuplicateMilestoneDocument,
    variables,
    headers,
  );
