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
export type UpdateServiceMutationVariables = Types.Exact<{
  input: Types.ServiceLineItemUpdateInput;
}>;

export type UpdateServiceMutation = {
  __typename?: 'Mutation';
  serviceLineItemUpdate: { __typename?: 'ServiceLineItem'; id: string };
};

export const UpdateServiceDocument = `
    mutation updateService($input: ServiceLineItemUpdateInput!) {
  serviceLineItemUpdate(input: $input) {
    id
  }
}
    `;
export const useUpdateServiceMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateServiceMutation,
    TError,
    UpdateServiceMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateServiceMutation,
    TError,
    UpdateServiceMutationVariables,
    TContext
  >(
    ['updateService'],
    (variables?: UpdateServiceMutationVariables) =>
      fetcher<UpdateServiceMutation, UpdateServiceMutationVariables>(
        client,
        UpdateServiceDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateServiceMutation.getKey = () => ['updateService'];

useUpdateServiceMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateServiceMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateServiceMutation, UpdateServiceMutationVariables>(
    client,
    UpdateServiceDocument,
    variables,
    headers,
  );
