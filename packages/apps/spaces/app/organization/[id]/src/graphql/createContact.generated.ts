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
export type CreateContactMutationVariables = Types.Exact<{
  input: Types.ContactInput;
}>;

export type CreateContactMutation = {
  __typename?: 'Mutation';
  contact_Create: { __typename?: 'Contact'; id: string };
};

export const CreateContactDocument = `
    mutation createContact($input: ContactInput!) {
  contact_Create(input: $input) {
    id
  }
}
    `;
export const useCreateContactMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateContactMutation,
    TError,
    CreateContactMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    CreateContactMutation,
    TError,
    CreateContactMutationVariables,
    TContext
  >(
    ['createContact'],
    (variables?: CreateContactMutationVariables) =>
      fetcher<CreateContactMutation, CreateContactMutationVariables>(
        client,
        CreateContactDocument,
        variables,
        headers,
      )(),
    options,
  );
useCreateContactMutation.getKey = () => ['createContact'];

useCreateContactMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateContactMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateContactMutation, CreateContactMutationVariables>(
    client,
    CreateContactDocument,
    variables,
    headers,
  );
