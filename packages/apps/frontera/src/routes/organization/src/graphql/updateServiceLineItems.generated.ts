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
export type UpdateServicesMutationVariables = Types.Exact<{
  input: Types.ServiceLineItemBulkUpdateInput;
}>;

export type UpdateServicesMutation = {
  __typename?: 'Mutation';
  serviceLineItem_BulkUpdate: Array<string>;
};

export const UpdateServicesDocument = `
    mutation updateServices($input: ServiceLineItemBulkUpdateInput!) {
  serviceLineItem_BulkUpdate(input: $input)
}
    `;

export const useUpdateServicesMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateServicesMutation,
    TError,
    UpdateServicesMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    UpdateServicesMutation,
    TError,
    UpdateServicesMutationVariables,
    TContext
  >({
    mutationKey: ['updateServices'],
    mutationFn: (variables?: UpdateServicesMutationVariables) =>
      fetcher<UpdateServicesMutation, UpdateServicesMutationVariables>(
        client,
        UpdateServicesDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useUpdateServicesMutation.getKey = () => ['updateServices'];

useUpdateServicesMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateServicesMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateServicesMutation, UpdateServicesMutationVariables>(
    client,
    UpdateServicesDocument,
    variables,
    headers,
  );
