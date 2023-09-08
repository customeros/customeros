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
export type AddOrganizationNoteMutationVariables = Types.Exact<{
  organzationId: Types.Scalars['ID'];
  input: Types.NoteInput;
}>;

export type AddOrganizationNoteMutation = {
  __typename?: 'Mutation';
  note_CreateForOrganization: { __typename?: 'Note'; id: string };
};

export const AddOrganizationNoteDocument = `
    mutation addOrganizationNote($organzationId: ID!, $input: NoteInput!) {
  note_CreateForOrganization(organizationId: $organzationId, input: $input) {
    id
  }
}
    `;
export const useAddOrganizationNoteMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddOrganizationNoteMutation,
    TError,
    AddOrganizationNoteMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddOrganizationNoteMutation,
    TError,
    AddOrganizationNoteMutationVariables,
    TContext
  >(
    ['addOrganizationNote'],
    (variables?: AddOrganizationNoteMutationVariables) =>
      fetcher<
        AddOrganizationNoteMutation,
        AddOrganizationNoteMutationVariables
      >(client, AddOrganizationNoteDocument, variables, headers)(),
    options,
  );
useAddOrganizationNoteMutation.getKey = () => ['addOrganizationNote'];

useAddOrganizationNoteMutation.fetcher = (
  client: GraphQLClient,
  variables: AddOrganizationNoteMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<AddOrganizationNoteMutation, AddOrganizationNoteMutationVariables>(
    client,
    AddOrganizationNoteDocument,
    variables,
    headers,
  );
