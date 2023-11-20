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
export type UpdateContractMutationVariables = Types.Exact<{
  input: Types.ContractUpdateInput;
}>;

export type UpdateContractMutation = {
  __typename?: 'Mutation';
  contract_Update: { __typename?: 'Contract'; id: string };
};

export const UpdateContractDocument = `
    mutation updateContract($input: ContractUpdateInput!) {
  contract_Update(input: $input) {
    id
  }
}
    `;
export const useUpdateContractMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateContractMutation,
    TError,
    UpdateContractMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateContractMutation,
    TError,
    UpdateContractMutationVariables,
    TContext
  >(
    ['updateContract'],
    (variables?: UpdateContractMutationVariables) =>
      fetcher<UpdateContractMutation, UpdateContractMutationVariables>(
        client,
        UpdateContractDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateContractMutation.getKey = () => ['updateContract'];

useUpdateContractMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateContractMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateContractMutation, UpdateContractMutationVariables>(
    client,
    UpdateContractDocument,
    variables,
    headers,
  );
