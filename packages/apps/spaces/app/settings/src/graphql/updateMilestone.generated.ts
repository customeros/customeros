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
export type UpdateMilestoneMutationVariables = Types.Exact<{
  input: Types.MasterPlanMilestoneUpdateInput;
}>;


export type UpdateMilestoneMutation = { __typename?: 'Mutation', masterPlanMilestone_Update: { __typename?: 'MasterPlanMilestone', id: string, name: string, order: any, durationHours: any, optional: boolean, items: Array<string>, retired: boolean } };



export const UpdateMilestoneDocument = `
    mutation updateMilestone($input: MasterPlanMilestoneUpdateInput!) {
  masterPlanMilestone_Update(input: $input) {
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

export const useUpdateMilestoneMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<UpdateMilestoneMutation, TError, UpdateMilestoneMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<UpdateMilestoneMutation, TError, UpdateMilestoneMutationVariables, TContext>(
      {
    mutationKey: ['updateMilestone'],
    mutationFn: (variables?: UpdateMilestoneMutationVariables) => fetcher<UpdateMilestoneMutation, UpdateMilestoneMutationVariables>(client, UpdateMilestoneDocument, variables, headers)(),
    ...options
  }
    )};

useUpdateMilestoneMutation.getKey = () => ['updateMilestone'];


useUpdateMilestoneMutation.fetcher = (client: GraphQLClient, variables: UpdateMilestoneMutationVariables, headers?: RequestInit['headers']) => fetcher<UpdateMilestoneMutation, UpdateMilestoneMutationVariables>(client, UpdateMilestoneDocument, variables, headers);
