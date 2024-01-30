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
export type CreateOnboardingPlanMutationVariables = Types.Exact<{
  input: Types.OrganizationPlanInput;
}>;

export type CreateOnboardingPlanMutation = {
  __typename?: 'Mutation';
  organizationPlan_Create: { __typename?: 'OrganizationPlan'; id: string };
};

export const CreateOnboardingPlanDocument = `
    mutation createOnboardingPlan($input: OrganizationPlanInput!) {
  organizationPlan_Create(input: $input) {
    id
  }
}
    `;

export const useCreateOnboardingPlanMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateOnboardingPlanMutation,
    TError,
    CreateOnboardingPlanMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    CreateOnboardingPlanMutation,
    TError,
    CreateOnboardingPlanMutationVariables,
    TContext
  >({
    mutationKey: ['createOnboardingPlan'],
    mutationFn: (variables?: CreateOnboardingPlanMutationVariables) =>
      fetcher<
        CreateOnboardingPlanMutation,
        CreateOnboardingPlanMutationVariables
      >(client, CreateOnboardingPlanDocument, variables, headers)(),
    ...options,
  });
};

useCreateOnboardingPlanMutation.getKey = () => ['createOnboardingPlan'];

useCreateOnboardingPlanMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateOnboardingPlanMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateOnboardingPlanMutation, CreateOnboardingPlanMutationVariables>(
    client,
    CreateOnboardingPlanDocument,
    variables,
    headers,
  );
