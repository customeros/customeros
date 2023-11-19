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
export type CreateServiceMutationVariables = Types.Exact<{
  input: Types.ServiceLineItemInput;
}>;

export type CreateServiceMutation = {
  __typename?: 'Mutation';
  serviceLineItemCreate: { __typename?: 'ServiceLineItem'; id: string };
};

export const CreateServiceDocument = `
    mutation createService($input: ServiceLineItemInput!) {
  serviceLineItemCreate(input: $input) {
    id
  }
}
    `;
export const useCreateServiceMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateServiceMutation,
    TError,
    CreateServiceMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    CreateServiceMutation,
    TError,
    CreateServiceMutationVariables,
    TContext
  >(
    ['createService'],
    (variables?: CreateServiceMutationVariables) =>
      fetcher<CreateServiceMutation, CreateServiceMutationVariables>(
        client,
        CreateServiceDocument,
        variables,
        headers,
      )(),
    options,
  );
useCreateServiceMutation.getKey = () => ['createService'];

useCreateServiceMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateServiceMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateServiceMutation, CreateServiceMutationVariables>(
    client,
    CreateServiceDocument,
    variables,
    headers,
  );
