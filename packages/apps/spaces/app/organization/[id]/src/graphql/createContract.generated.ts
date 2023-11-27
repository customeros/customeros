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
export type CreateContractMutationVariables = Types.Exact<{
  input: Types.ContractInput;
}>;

export type CreateContractMutation = {
  __typename?: 'Mutation';
  contract_Create: { __typename?: 'Contract'; id: string };
};

export const CreateContractDocument = `
    mutation createContract($input: ContractInput!) {
  contract_Create(input: $input) {
    id
  }
}
    `;
export const useCreateContractMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateContractMutation,
    TError,
    CreateContractMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    CreateContractMutation,
    TError,
    CreateContractMutationVariables,
    TContext
  >(
    ['createContract'],
    (variables?: CreateContractMutationVariables) =>
      fetcher<CreateContractMutation, CreateContractMutationVariables>(
        client,
        CreateContractDocument,
        variables,
        headers,
      )(),
    options,
  );
useCreateContractMutation.getKey = () => ['createContract'];

useCreateContractMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateContractMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateContractMutation, CreateContractMutationVariables>(
    client,
    CreateContractDocument,
    variables,
    headers,
  );
