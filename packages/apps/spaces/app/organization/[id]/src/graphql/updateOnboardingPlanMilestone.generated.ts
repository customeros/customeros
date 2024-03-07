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
export type UpdateOnboardingPlanMilestoneMutationVariables = Types.Exact<{
  input: Types.OrganizationPlanMilestoneUpdateInput;
}>;

export type UpdateOnboardingPlanMilestoneMutation = {
  __typename?: 'Mutation';
  organizationPlanMilestone_Update: {
    __typename?: 'OrganizationPlanMilestone';
    id: string;
  };
};

export const UpdateOnboardingPlanMilestoneDocument = `
    mutation updateOnboardingPlanMilestone($input: OrganizationPlanMilestoneUpdateInput!) {
  organizationPlanMilestone_Update(input: $input) {
    id
  }
}
    `;

export const useUpdateOnboardingPlanMilestoneMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateOnboardingPlanMilestoneMutation,
    TError,
    UpdateOnboardingPlanMilestoneMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    UpdateOnboardingPlanMilestoneMutation,
    TError,
    UpdateOnboardingPlanMilestoneMutationVariables,
    TContext
  >({
    mutationKey: ['updateOnboardingPlanMilestone'],
    mutationFn: (variables?: UpdateOnboardingPlanMilestoneMutationVariables) =>
      fetcher<
        UpdateOnboardingPlanMilestoneMutation,
        UpdateOnboardingPlanMilestoneMutationVariables
      >(client, UpdateOnboardingPlanMilestoneDocument, variables, headers)(),
    ...options,
  });
};

useUpdateOnboardingPlanMilestoneMutation.getKey = () => [
  'updateOnboardingPlanMilestone',
];

useUpdateOnboardingPlanMilestoneMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateOnboardingPlanMilestoneMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    UpdateOnboardingPlanMilestoneMutation,
    UpdateOnboardingPlanMilestoneMutationVariables
  >(client, UpdateOnboardingPlanMilestoneDocument, variables, headers);
