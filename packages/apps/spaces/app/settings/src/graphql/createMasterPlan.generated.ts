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
export type CreateMasterPlanMutationVariables = Types.Exact<{
  input: Types.MasterPlanInput;
}>;

export type CreateMasterPlanMutation = {
  __typename?: 'Mutation';
  masterPlan_Create: { __typename?: 'MasterPlan'; id: string; name: string };
};

export const CreateMasterPlanDocument = `
    mutation createMasterPlan($input: MasterPlanInput!) {
  masterPlan_Create(input: $input) {
    id
    name
  }
}
    `;

export const useCreateMasterPlanMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateMasterPlanMutation,
    TError,
    CreateMasterPlanMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    CreateMasterPlanMutation,
    TError,
    CreateMasterPlanMutationVariables,
    TContext
  >({
    mutationKey: ['createMasterPlan'],
    mutationFn: (variables?: CreateMasterPlanMutationVariables) =>
      fetcher<CreateMasterPlanMutation, CreateMasterPlanMutationVariables>(
        client,
        CreateMasterPlanDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useCreateMasterPlanMutation.getKey = () => ['createMasterPlan'];

useCreateMasterPlanMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateMasterPlanMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateMasterPlanMutation, CreateMasterPlanMutationVariables>(
    client,
    CreateMasterPlanDocument,
    variables,
    headers,
  );
