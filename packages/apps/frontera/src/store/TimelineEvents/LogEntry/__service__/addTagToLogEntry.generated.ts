// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

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
export type AddTagToLogEntryMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  input: Types.TagIdOrNameInput;
}>;

export type AddTagToLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_AddTag: string;
};

export const AddTagToLogEntryDocument = `
    mutation addTagToLogEntry($id: ID!, $input: TagIdOrNameInput!) {
  logEntry_AddTag(id: $id, input: $input)
}
    `;

export const useAddTagToLogEntryMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddTagToLogEntryMutation,
    TError,
    AddTagToLogEntryMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    AddTagToLogEntryMutation,
    TError,
    AddTagToLogEntryMutationVariables,
    TContext
  >({
    mutationKey: ['addTagToLogEntry'],
    mutationFn: (variables?: AddTagToLogEntryMutationVariables) =>
      fetcher<AddTagToLogEntryMutation, AddTagToLogEntryMutationVariables>(
        client,
        AddTagToLogEntryDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useAddTagToLogEntryMutation.getKey = () => ['addTagToLogEntry'];

useAddTagToLogEntryMutation.fetcher = (
  client: GraphQLClient,
  variables: AddTagToLogEntryMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<AddTagToLogEntryMutation, AddTagToLogEntryMutationVariables>(
    client,
    AddTagToLogEntryDocument,
    variables,
    headers,
  );
