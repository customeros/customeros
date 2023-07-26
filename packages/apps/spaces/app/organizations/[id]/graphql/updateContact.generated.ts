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
export type UpdateContactMutationVariables = Types.Exact<{
  input: Types.ContactUpdateInput;
}>;

export type UpdateContactMutation = {
  __typename?: 'Mutation';
  contact_Update: { __typename?: 'Contact'; id: string };
};

export const UpdateContactDocument = `
    mutation updateContact($input: ContactUpdateInput!) {
  contact_Update(input: $input) {
    id
  }
}
    `;
export const useUpdateContactMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateContactMutation,
    TError,
    UpdateContactMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateContactMutation,
    TError,
    UpdateContactMutationVariables,
    TContext
  >(
    ['updateContact'],
    (variables?: UpdateContactMutationVariables) =>
      fetcher<UpdateContactMutation, UpdateContactMutationVariables>(
        client,
        UpdateContactDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateContactMutation.getKey = () => ['updateContact'];

useUpdateContactMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateContactMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateContactMutation, UpdateContactMutationVariables>(
    client,
    UpdateContactDocument,
    variables,
    headers,
  );
