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
export type UpdateBillingDetailsMutationVariables = Types.Exact<{
  input: Types.BillingDetailsInput;
}>;

export type UpdateBillingDetailsMutation = {
  __typename?: 'Mutation';
  organization_UpdateBillingDetailsAsync: string;
};

export const UpdateBillingDetailsDocument = `
    mutation updateBillingDetails($input: BillingDetailsInput!) {
  organization_UpdateBillingDetailsAsync(input: $input)
}
    `;
export const useUpdateBillingDetailsMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateBillingDetailsMutation,
    TError,
    UpdateBillingDetailsMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateBillingDetailsMutation,
    TError,
    UpdateBillingDetailsMutationVariables,
    TContext
  >(
    ['updateBillingDetails'],
    (variables?: UpdateBillingDetailsMutationVariables) =>
      fetcher<
        UpdateBillingDetailsMutation,
        UpdateBillingDetailsMutationVariables
      >(client, UpdateBillingDetailsDocument, variables, headers)(),
    options,
  );
useUpdateBillingDetailsMutation.getKey = () => ['updateBillingDetails'];

useUpdateBillingDetailsMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateBillingDetailsMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateBillingDetailsMutation, UpdateBillingDetailsMutationVariables>(
    client,
    UpdateBillingDetailsDocument,
    variables,
    headers,
  );
