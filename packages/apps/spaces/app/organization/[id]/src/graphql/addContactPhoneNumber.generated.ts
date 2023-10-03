// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../types/__generated__/graphql.types';

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
export type AddContactPhoneNumberMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
  input: Types.PhoneNumberInput;
}>;

export type AddContactPhoneNumberMutation = {
  __typename?: 'Mutation';
  phoneNumberMergeToContact: {
    __typename?: 'PhoneNumber';
    id: string;
    rawPhoneNumber?: string | null;
  };
};

export const AddContactPhoneNumberDocument = `
    mutation addContactPhoneNumber($contactId: ID!, $input: PhoneNumberInput!) {
  phoneNumberMergeToContact(contactId: $contactId, input: $input) {
    id
    rawPhoneNumber
  }
}
    `;
export const useAddContactPhoneNumberMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddContactPhoneNumberMutation,
    TError,
    AddContactPhoneNumberMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddContactPhoneNumberMutation,
    TError,
    AddContactPhoneNumberMutationVariables,
    TContext
  >(
    ['addContactPhoneNumber'],
    (variables?: AddContactPhoneNumberMutationVariables) =>
      fetcher<
        AddContactPhoneNumberMutation,
        AddContactPhoneNumberMutationVariables
      >(client, AddContactPhoneNumberDocument, variables, headers)(),
    options,
  );
useAddContactPhoneNumberMutation.getKey = () => ['addContactPhoneNumber'];

useAddContactPhoneNumberMutation.fetcher = (
  client: GraphQLClient,
  variables: AddContactPhoneNumberMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    AddContactPhoneNumberMutation,
    AddContactPhoneNumberMutationVariables
  >(client, AddContactPhoneNumberDocument, variables, headers);
