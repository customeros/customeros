// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useMutation, UseMutationOptions } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type DeleteBankAccountMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type DeleteBankAccountMutation = { __typename?: 'Mutation', bankAccount_Delete: { __typename?: 'DeleteResponse', accepted: boolean, completed: boolean } };



export const DeleteBankAccountDocument = `
    mutation deleteBankAccount($id: ID!) {
  bankAccount_Delete(id: $id) {
    accepted
    completed
  }
}
    `;

export const useDeleteBankAccountMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<DeleteBankAccountMutation, TError, DeleteBankAccountMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<DeleteBankAccountMutation, TError, DeleteBankAccountMutationVariables, TContext>(
      {
    mutationKey: ['deleteBankAccount'],
    mutationFn: (variables?: DeleteBankAccountMutationVariables) => fetcher<DeleteBankAccountMutation, DeleteBankAccountMutationVariables>(client, DeleteBankAccountDocument, variables, headers)(),
    ...options
  }
    )};

useDeleteBankAccountMutation.getKey = () => ['deleteBankAccount'];


useDeleteBankAccountMutation.fetcher = (client: GraphQLClient, variables: DeleteBankAccountMutationVariables, headers?: RequestInit['headers']) => fetcher<DeleteBankAccountMutation, DeleteBankAccountMutationVariables>(client, DeleteBankAccountDocument, variables, headers);
