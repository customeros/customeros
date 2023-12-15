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
export type UpdateOnboardingStatusMutationVariables = Types.Exact<{
  input: Types.OnboardingStatusInput;
}>;

export type UpdateOnboardingStatusMutation = {
  __typename?: 'Mutation';
  organization_UpdateOnboardingStatus: {
    __typename?: 'Organization';
    id: string;
    accountDetails?: {
      __typename?: 'OrgAccountDetails';
      onboarding?: {
        __typename?: 'OnboardingDetails';
        status: Types.OnboardingStatus;
        comments?: string | null;
      } | null;
    } | null;
  };
};

export const UpdateOnboardingStatusDocument = `
    mutation updateOnboardingStatus($input: OnboardingStatusInput!) {
  organization_UpdateOnboardingStatus(input: $input) {
    id
    accountDetails {
      onboarding {
        status
        comments
      }
    }
  }
}
    `;
export const useUpdateOnboardingStatusMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateOnboardingStatusMutation,
    TError,
    UpdateOnboardingStatusMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateOnboardingStatusMutation,
    TError,
    UpdateOnboardingStatusMutationVariables,
    TContext
  >(
    ['updateOnboardingStatus'],
    (variables?: UpdateOnboardingStatusMutationVariables) =>
      fetcher<
        UpdateOnboardingStatusMutation,
        UpdateOnboardingStatusMutationVariables
      >(client, UpdateOnboardingStatusDocument, variables, headers)(),
    options,
  );
useUpdateOnboardingStatusMutation.getKey = () => ['updateOnboardingStatus'];

useUpdateOnboardingStatusMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateOnboardingStatusMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    UpdateOnboardingStatusMutation,
    UpdateOnboardingStatusMutationVariables
  >(client, UpdateOnboardingStatusDocument, variables, headers);
