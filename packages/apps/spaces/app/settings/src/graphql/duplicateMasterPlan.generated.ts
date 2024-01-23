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
export type DuplicateMasterPlanMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type DuplicateMasterPlanMutation = {
  __typename?: 'Mutation';
  masterPlan_Duplicate: { __typename?: 'MasterPlan'; id: string };
};

export const DuplicateMasterPlanDocument = `
    mutation duplicateMasterPlan($id: ID!) {
  masterPlan_Duplicate(id: $id) {
    id
  }
}
    `;

export const useDuplicateMasterPlanMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    DuplicateMasterPlanMutation,
    TError,
    DuplicateMasterPlanMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    DuplicateMasterPlanMutation,
    TError,
    DuplicateMasterPlanMutationVariables,
    TContext
  >({
    mutationKey: ['duplicateMasterPlan'],
    mutationFn: (variables?: DuplicateMasterPlanMutationVariables) =>
      fetcher<
        DuplicateMasterPlanMutation,
        DuplicateMasterPlanMutationVariables
      >(client, DuplicateMasterPlanDocument, variables, headers)(),
    ...options,
  });
};

useDuplicateMasterPlanMutation.getKey = () => ['duplicateMasterPlan'];

useDuplicateMasterPlanMutation.fetcher = (
  client: GraphQLClient,
  variables: DuplicateMasterPlanMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<DuplicateMasterPlanMutation, DuplicateMasterPlanMutationVariables>(
    client,
    DuplicateMasterPlanDocument,
    variables,
    headers,
  );
