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
export type UpdateContactRoleMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
  input: Types.JobRoleUpdateInput;
}>;

export type UpdateContactRoleMutation = {
  __typename?: 'Mutation';
  jobRole_Update: { __typename?: 'JobRole'; id: string };
};

export const UpdateContactRoleDocument = `
    mutation updateContactRole($contactId: ID!, $input: JobRoleUpdateInput!) {
  jobRole_Update(contactId: $contactId, input: $input) {
    id
  }
}
    `;
export const useUpdateContactRoleMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateContactRoleMutation,
    TError,
    UpdateContactRoleMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateContactRoleMutation,
    TError,
    UpdateContactRoleMutationVariables,
    TContext
  >(
    ['updateContactRole'],
    (variables?: UpdateContactRoleMutationVariables) =>
      fetcher<UpdateContactRoleMutation, UpdateContactRoleMutationVariables>(
        client,
        UpdateContactRoleDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateContactRoleMutation.getKey = () => ['updateContactRole'];

useUpdateContactRoleMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateContactRoleMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateContactRoleMutation, UpdateContactRoleMutationVariables>(
    client,
    UpdateContactRoleDocument,
    variables,
    headers,
  );
