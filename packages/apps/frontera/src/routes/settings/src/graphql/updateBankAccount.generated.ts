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
export type UpdateBankAccountMutationVariables = Types.Exact<{
  input: Types.BankAccountUpdateInput;
}>;


export type UpdateBankAccountMutation = { __typename?: 'Mutation', bankAccount_Update: { __typename?: 'BankAccount', currency?: Types.Currency | null, bankName?: string | null, bankTransferEnabled: boolean, iban?: string | null, bic?: string | null, sortCode?: string | null, accountNumber?: string | null, routingNumber?: string | null } };



export const UpdateBankAccountDocument = `
    mutation updateBankAccount($input: BankAccountUpdateInput!) {
  bankAccount_Update(input: $input) {
    currency
    bankName
    bankTransferEnabled
    iban
    bic
    sortCode
    accountNumber
    routingNumber
  }
}
    `;

export const useUpdateBankAccountMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<UpdateBankAccountMutation, TError, UpdateBankAccountMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<UpdateBankAccountMutation, TError, UpdateBankAccountMutationVariables, TContext>(
      {
    mutationKey: ['updateBankAccount'],
    mutationFn: (variables?: UpdateBankAccountMutationVariables) => fetcher<UpdateBankAccountMutation, UpdateBankAccountMutationVariables>(client, UpdateBankAccountDocument, variables, headers)(),
    ...options
  }
    )};

useUpdateBankAccountMutation.getKey = () => ['updateBankAccount'];


useUpdateBankAccountMutation.fetcher = (client: GraphQLClient, variables: UpdateBankAccountMutationVariables, headers?: RequestInit['headers']) => fetcher<UpdateBankAccountMutation, UpdateBankAccountMutationVariables>(client, UpdateBankAccountDocument, variables, headers);
