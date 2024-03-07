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
export type UpdateMasterPlanMutationVariables = Types.Exact<{
  input: Types.MasterPlanUpdateInput;
}>;


export type UpdateMasterPlanMutation = { __typename?: 'Mutation', masterPlan_Update: { __typename?: 'MasterPlan', id: string } };



export const UpdateMasterPlanDocument = `
    mutation updateMasterPlan($input: MasterPlanUpdateInput!) {
  masterPlan_Update(input: $input) {
    id
  }
}
    `;

export const useUpdateMasterPlanMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<UpdateMasterPlanMutation, TError, UpdateMasterPlanMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<UpdateMasterPlanMutation, TError, UpdateMasterPlanMutationVariables, TContext>(
      {
    mutationKey: ['updateMasterPlan'],
    mutationFn: (variables?: UpdateMasterPlanMutationVariables) => fetcher<UpdateMasterPlanMutation, UpdateMasterPlanMutationVariables>(client, UpdateMasterPlanDocument, variables, headers)(),
    ...options
  }
    )};

useUpdateMasterPlanMutation.getKey = () => ['updateMasterPlan'];


useUpdateMasterPlanMutation.fetcher = (client: GraphQLClient, variables: UpdateMasterPlanMutationVariables, headers?: RequestInit['headers']) => fetcher<UpdateMasterPlanMutation, UpdateMasterPlanMutationVariables>(client, UpdateMasterPlanDocument, variables, headers);
