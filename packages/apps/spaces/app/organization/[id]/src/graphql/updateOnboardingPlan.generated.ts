// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

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
export type UpdateOnboardingPlanMutationVariables = Types.Exact<{
  input: Types.OrganizationPlanUpdateInput;
}>;

export type UpdateOnboardingPlanMutation = {
  __typename?: 'Mutation';
  organizationPlan_Update: { __typename?: 'OrganizationPlan'; id: string };
};

export const UpdateOnboardingPlanDocument = `
    mutation updateOnboardingPlan($input: OrganizationPlanUpdateInput!) {
  organizationPlan_Update(input: $input) {
    id
  }
}
    `;

export const useUpdateOnboardingPlanMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateOnboardingPlanMutation,
    TError,
    UpdateOnboardingPlanMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    UpdateOnboardingPlanMutation,
    TError,
    UpdateOnboardingPlanMutationVariables,
    TContext
  >({
    mutationKey: ['updateOnboardingPlan'],
    mutationFn: (variables?: UpdateOnboardingPlanMutationVariables) =>
      fetcher<
        UpdateOnboardingPlanMutation,
        UpdateOnboardingPlanMutationVariables
      >(client, UpdateOnboardingPlanDocument, variables, headers)(),
    ...options,
  });
};

useUpdateOnboardingPlanMutation.getKey = () => ['updateOnboardingPlan'];

useUpdateOnboardingPlanMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateOnboardingPlanMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateOnboardingPlanMutation, UpdateOnboardingPlanMutationVariables>(
    client,
    UpdateOnboardingPlanDocument,
    variables,
    headers,
  );
