// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useMutation, UseMutationOptions } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type CreateDefaultMasterPlanMutationVariables = Types.Exact<{
  input: Types.MasterPlanInput;
}>;


export type CreateDefaultMasterPlanMutation = { __typename?: 'Mutation', masterPlan_CreateDefault: { __typename?: 'MasterPlan', id: string, name: string } };



export const CreateDefaultMasterPlanDocument = `
    mutation createDefaultMasterPlan($input: MasterPlanInput!) {
  masterPlan_CreateDefault(input: $input) {
    id
    name
  }
}
    `;

export const useCreateDefaultMasterPlanMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<CreateDefaultMasterPlanMutation, TError, CreateDefaultMasterPlanMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<CreateDefaultMasterPlanMutation, TError, CreateDefaultMasterPlanMutationVariables, TContext>(
      {
    mutationKey: ['createDefaultMasterPlan'],
    mutationFn: (variables?: CreateDefaultMasterPlanMutationVariables) => fetcher<CreateDefaultMasterPlanMutation, CreateDefaultMasterPlanMutationVariables>(client, CreateDefaultMasterPlanDocument, variables, headers)(),
    ...options
  }
    )};

useCreateDefaultMasterPlanMutation.getKey = () => ['createDefaultMasterPlan'];


useCreateDefaultMasterPlanMutation.fetcher = (client: GraphQLClient, variables: CreateDefaultMasterPlanMutationVariables, headers?: RequestInit['headers']) => fetcher<CreateDefaultMasterPlanMutation, CreateDefaultMasterPlanMutationVariables>(client, CreateDefaultMasterPlanDocument, variables, headers);
