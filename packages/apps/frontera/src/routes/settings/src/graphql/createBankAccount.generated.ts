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
export type CreateBankAccountMutationVariables = Types.Exact<{
  input: Types.BankAccountCreateInput;
}>;


export type CreateBankAccountMutation = { __typename?: 'Mutation', bankAccount_Create: { __typename?: 'BankAccount', currency?: Types.Currency | null, bankName?: string | null, bankTransferEnabled: boolean, iban?: string | null, bic?: string | null, sortCode?: string | null, accountNumber?: string | null, routingNumber?: string | null } };



export const CreateBankAccountDocument = `
    mutation createBankAccount($input: BankAccountCreateInput!) {
  bankAccount_Create(input: $input) {
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

export const useCreateBankAccountMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<CreateBankAccountMutation, TError, CreateBankAccountMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<CreateBankAccountMutation, TError, CreateBankAccountMutationVariables, TContext>(
      {
    mutationKey: ['createBankAccount'],
    mutationFn: (variables?: CreateBankAccountMutationVariables) => fetcher<CreateBankAccountMutation, CreateBankAccountMutationVariables>(client, CreateBankAccountDocument, variables, headers)(),
    ...options
  }
    )};

useCreateBankAccountMutation.getKey = () => ['createBankAccount'];


useCreateBankAccountMutation.fetcher = (client: GraphQLClient, variables: CreateBankAccountMutationVariables, headers?: RequestInit['headers']) => fetcher<CreateBankAccountMutation, CreateBankAccountMutationVariables>(client, CreateBankAccountDocument, variables, headers);
