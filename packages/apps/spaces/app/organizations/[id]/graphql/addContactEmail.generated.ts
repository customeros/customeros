// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../types/__generated__/graphql.types';

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
export type AddContactEmailMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
  input: Types.EmailInput;
}>;

export type AddContactEmailMutation = {
  __typename?: 'Mutation';
  emailMergeToContact: {
    __typename?: 'Email';
    id: string;
    email?: string | null;
  };
};

export const AddContactEmailDocument = `
    mutation addContactEmail($contactId: ID!, $input: EmailInput!) {
  emailMergeToContact(contactId: $contactId, input: $input) {
    id
    email
  }
}
    `;
export const useAddContactEmailMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddContactEmailMutation,
    TError,
    AddContactEmailMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddContactEmailMutation,
    TError,
    AddContactEmailMutationVariables,
    TContext
  >(
    ['addContactEmail'],
    (variables?: AddContactEmailMutationVariables) =>
      fetcher<AddContactEmailMutation, AddContactEmailMutationVariables>(
        client,
        AddContactEmailDocument,
        variables,
        headers,
      )(),
    options,
  );
useAddContactEmailMutation.getKey = () => ['addContactEmail'];

useAddContactEmailMutation.fetcher = (
  client: GraphQLClient,
  variables: AddContactEmailMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<AddContactEmailMutation, AddContactEmailMutationVariables>(
    client,
    AddContactEmailDocument,
    variables,
    headers,
  );
