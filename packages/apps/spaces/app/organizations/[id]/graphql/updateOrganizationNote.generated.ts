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
export type UpdateOrganizationNoteMutationVariables = Types.Exact<{
  input: Types.NoteUpdateInput;
}>;

export type UpdateOrganizationNoteMutation = {
  __typename?: 'Mutation';
  note_Update: { __typename?: 'Note'; id: string };
};

export const UpdateOrganizationNoteDocument = `
    mutation updateOrganizationNote($input: NoteUpdateInput!) {
  note_Update(input: $input) {
    id
  }
}
    `;
export const useUpdateOrganizationNoteMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateOrganizationNoteMutation,
    TError,
    UpdateOrganizationNoteMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateOrganizationNoteMutation,
    TError,
    UpdateOrganizationNoteMutationVariables,
    TContext
  >(
    ['updateOrganizationNote'],
    (variables?: UpdateOrganizationNoteMutationVariables) =>
      fetcher<
        UpdateOrganizationNoteMutation,
        UpdateOrganizationNoteMutationVariables
      >(client, UpdateOrganizationNoteDocument, variables, headers)(),
    options,
  );
useUpdateOrganizationNoteMutation.getKey = () => ['updateOrganizationNote'];

useUpdateOrganizationNoteMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateOrganizationNoteMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    UpdateOrganizationNoteMutation,
    UpdateOrganizationNoteMutationVariables
  >(client, UpdateOrganizationNoteDocument, variables, headers);
