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
export type UpdateMilestonesMutationVariables = Types.Exact<{
  input: Array<Types.MasterPlanMilestoneUpdateInput> | Types.MasterPlanMilestoneUpdateInput;
}>;


export type UpdateMilestonesMutation = { __typename?: 'Mutation', masterPlanMilestone_BulkUpdate: Array<{ __typename?: 'MasterPlanMilestone', id: string, name: string, order: any, durationHours: any, optional: boolean, items: Array<string>, retired: boolean }> };



export const UpdateMilestonesDocument = `
    mutation updateMilestones($input: [MasterPlanMilestoneUpdateInput!]!) {
  masterPlanMilestone_BulkUpdate(input: $input) {
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

export const useUpdateMilestonesMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<UpdateMilestonesMutation, TError, UpdateMilestonesMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<UpdateMilestonesMutation, TError, UpdateMilestonesMutationVariables, TContext>(
      {
    mutationKey: ['updateMilestones'],
    mutationFn: (variables?: UpdateMilestonesMutationVariables) => fetcher<UpdateMilestonesMutation, UpdateMilestonesMutationVariables>(client, UpdateMilestonesDocument, variables, headers)(),
    ...options
  }
    )};

useUpdateMilestonesMutation.getKey = () => ['updateMilestones'];


useUpdateMilestonesMutation.fetcher = (client: GraphQLClient, variables: UpdateMilestonesMutationVariables, headers?: RequestInit['headers']) => fetcher<UpdateMilestonesMutation, UpdateMilestonesMutationVariables>(client, UpdateMilestonesDocument, variables, headers);
