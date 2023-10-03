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
export type UpdateContactPhoneNumberMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
  input: Types.PhoneNumberUpdateInput;
}>;

export type UpdateContactPhoneNumberMutation = {
  __typename?: 'Mutation';
  phoneNumberUpdateInContact: {
    __typename?: 'PhoneNumber';
    id: string;
    rawPhoneNumber?: string | null;
  };
};

export const UpdateContactPhoneNumberDocument = `
    mutation updateContactPhoneNumber($contactId: ID!, $input: PhoneNumberUpdateInput!) {
  phoneNumberUpdateInContact(contactId: $contactId, input: $input) {
    id
    rawPhoneNumber
  }
}
    `;
export const useUpdateContactPhoneNumberMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateContactPhoneNumberMutation,
    TError,
    UpdateContactPhoneNumberMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateContactPhoneNumberMutation,
    TError,
    UpdateContactPhoneNumberMutationVariables,
    TContext
  >(
    ['updateContactPhoneNumber'],
    (variables?: UpdateContactPhoneNumberMutationVariables) =>
      fetcher<
        UpdateContactPhoneNumberMutation,
        UpdateContactPhoneNumberMutationVariables
      >(client, UpdateContactPhoneNumberDocument, variables, headers)(),
    options,
  );
useUpdateContactPhoneNumberMutation.getKey = () => ['updateContactPhoneNumber'];

useUpdateContactPhoneNumberMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateContactPhoneNumberMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    UpdateContactPhoneNumberMutation,
    UpdateContactPhoneNumberMutationVariables
  >(client, UpdateContactPhoneNumberDocument, variables, headers);
