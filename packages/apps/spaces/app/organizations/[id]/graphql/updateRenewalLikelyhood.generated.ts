// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

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
export type UpdateRenewalLikelihoodMutationVariables = Types.Exact<{
  input: Types.RenewalLikelihoodInput;
}>;

export type UpdateRenewalLikelihoodMutation = {
  __typename?: 'Mutation';
  organization_UpdateRenewalLikelihood: {
    __typename?: 'Organization';
    id: string;
  };
};

export const UpdateRenewalLikelihoodDocument = `
    mutation updateRenewalLikelihood($input: RenewalLikelihoodInput!) {
  organization_UpdateRenewalLikelihood(input: $input) {
    id
  }
}
    `;
export const useUpdateRenewalLikelihoodMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateRenewalLikelihoodMutation,
    TError,
    UpdateRenewalLikelihoodMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateRenewalLikelihoodMutation,
    TError,
    UpdateRenewalLikelihoodMutationVariables,
    TContext
  >(
    ['updateRenewalLikelihood'],
    (variables?: UpdateRenewalLikelihoodMutationVariables) =>
      fetcher<
        UpdateRenewalLikelihoodMutation,
        UpdateRenewalLikelihoodMutationVariables
      >(client, UpdateRenewalLikelihoodDocument, variables, headers)(),
    options,
  );
useUpdateRenewalLikelihoodMutation.getKey = () => ['updateRenewalLikelihood'];

useUpdateRenewalLikelihoodMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateRenewalLikelihoodMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    UpdateRenewalLikelihoodMutation,
    UpdateRenewalLikelihoodMutationVariables
  >(client, UpdateRenewalLikelihoodDocument, variables, headers);
