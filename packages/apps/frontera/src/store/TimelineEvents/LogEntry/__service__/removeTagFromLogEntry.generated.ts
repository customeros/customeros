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
export type RemoveTagFromLogEntryMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  input: Types.TagIdOrNameInput;
}>;

export type RemoveTagFromLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_RemoveTag: string;
};

export const RemoveTagFromLogEntryDocument = `
    mutation removeTagFromLogEntry($id: ID!, $input: TagIdOrNameInput!) {
  logEntry_RemoveTag(id: $id, input: $input)
}
    `;

export const useRemoveTagFromLogEntryMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveTagFromLogEntryMutation,
    TError,
    RemoveTagFromLogEntryMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    RemoveTagFromLogEntryMutation,
    TError,
    RemoveTagFromLogEntryMutationVariables,
    TContext
  >({
    mutationKey: ['removeTagFromLogEntry'],
    mutationFn: (variables?: RemoveTagFromLogEntryMutationVariables) =>
      fetcher<
        RemoveTagFromLogEntryMutation,
        RemoveTagFromLogEntryMutationVariables
      >(client, RemoveTagFromLogEntryDocument, variables, headers)(),
    ...options,
  });
};

useRemoveTagFromLogEntryMutation.getKey = () => ['removeTagFromLogEntry'];

useRemoveTagFromLogEntryMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveTagFromLogEntryMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    RemoveTagFromLogEntryMutation,
    RemoveTagFromLogEntryMutationVariables
  >(client, RemoveTagFromLogEntryDocument, variables, headers);
