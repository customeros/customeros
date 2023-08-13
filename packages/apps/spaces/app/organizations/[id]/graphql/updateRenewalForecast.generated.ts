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
export type UpdateRenewalForecastMutationVariables = Types.Exact<{
  input: Types.RenewalForecastInput;
}>;

export type UpdateRenewalForecastMutation = {
  __typename?: 'Mutation';
  organization_UpdateRenewalForecast: {
    __typename?: 'Organization';
    id: string;
  };
};

export const UpdateRenewalForecastDocument = `
    mutation updateRenewalForecast($input: RenewalForecastInput!) {
  organization_UpdateRenewalForecast(input: $input) {
    id
  }
}
    `;
export const useUpdateRenewalForecastMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateRenewalForecastMutation,
    TError,
    UpdateRenewalForecastMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateRenewalForecastMutation,
    TError,
    UpdateRenewalForecastMutationVariables,
    TContext
  >(
    ['updateRenewalForecast'],
    (variables?: UpdateRenewalForecastMutationVariables) =>
      fetcher<
        UpdateRenewalForecastMutation,
        UpdateRenewalForecastMutationVariables
      >(client, UpdateRenewalForecastDocument, variables, headers)(),
    options,
  );
useUpdateRenewalForecastMutation.getKey = () => ['updateRenewalForecast'];

useUpdateRenewalForecastMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateRenewalForecastMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    UpdateRenewalForecastMutation,
    UpdateRenewalForecastMutationVariables
  >(client, UpdateRenewalForecastDocument, variables, headers);
