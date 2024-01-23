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
export type CreateMilestoneMutationVariables = Types.Exact<{
  input: Types.MasterPlanMilestoneInput;
}>;

export type CreateMilestoneMutation = {
  __typename?: 'Mutation';
  masterPlanMilestone_Create: {
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

export const CreateMilestoneDocument = `
    mutation createMilestone($input: MasterPlanMilestoneInput!) {
  masterPlanMilestone_Create(input: $input) {
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

export const useCreateMilestoneMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateMilestoneMutation,
    TError,
    CreateMilestoneMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    CreateMilestoneMutation,
    TError,
    CreateMilestoneMutationVariables,
    TContext
  >({
    mutationKey: ['createMilestone'],
    mutationFn: (variables?: CreateMilestoneMutationVariables) =>
      fetcher<CreateMilestoneMutation, CreateMilestoneMutationVariables>(
        client,
        CreateMilestoneDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useCreateMilestoneMutation.getKey = () => ['createMilestone'];

useCreateMilestoneMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateMilestoneMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateMilestoneMutation, CreateMilestoneMutationVariables>(
    client,
    CreateMilestoneDocument,
    variables,
    headers,
  );
