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
export type AddOnboardingPlanMilestoneMutationVariables = Types.Exact<{
  input: Types.OrganizationPlanMilestoneInput;
}>;

export type AddOnboardingPlanMilestoneMutation = {
  __typename?: 'Mutation';
  organizationPlanMilestone_Create: {
    __typename?: 'OrganizationPlanMilestone';
    id: string;
  };
};

export const AddOnboardingPlanMilestoneDocument = `
    mutation addOnboardingPlanMilestone($input: OrganizationPlanMilestoneInput!) {
  organizationPlanMilestone_Create(input: $input) {
    id
  }
}
    `;

export const useAddOnboardingPlanMilestoneMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddOnboardingPlanMilestoneMutation,
    TError,
    AddOnboardingPlanMilestoneMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    AddOnboardingPlanMilestoneMutation,
    TError,
    AddOnboardingPlanMilestoneMutationVariables,
    TContext
  >({
    mutationKey: ['addOnboardingPlanMilestone'],
    mutationFn: (variables?: AddOnboardingPlanMilestoneMutationVariables) =>
      fetcher<
        AddOnboardingPlanMilestoneMutation,
        AddOnboardingPlanMilestoneMutationVariables
      >(client, AddOnboardingPlanMilestoneDocument, variables, headers)(),
    ...options,
  });
};

useAddOnboardingPlanMilestoneMutation.getKey = () => [
  'addOnboardingPlanMilestone',
];

useAddOnboardingPlanMilestoneMutation.fetcher = (
  client: GraphQLClient,
  variables: AddOnboardingPlanMilestoneMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    AddOnboardingPlanMilestoneMutation,
    AddOnboardingPlanMilestoneMutationVariables
  >(client, AddOnboardingPlanMilestoneDocument, variables, headers);
