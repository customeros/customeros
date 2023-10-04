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
export type UpdateLogEntryMutationVariables = Types.Exact<{
  id: Types.Scalars['ID'];
  input: Types.LogEntryUpdateInput;
}>;

export type UpdateLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_Update: string;
};

export const UpdateLogEntryDocument = `
    mutation updateLogEntry($id: ID!, $input: LogEntryUpdateInput!) {
  logEntry_Update(id: $id, input: $input)
}
    `;
export const useUpdateLogEntryMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateLogEntryMutation,
    TError,
    UpdateLogEntryMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateLogEntryMutation,
    TError,
    UpdateLogEntryMutationVariables,
    TContext
  >(
    ['updateLogEntry'],
    (variables?: UpdateLogEntryMutationVariables) =>
      fetcher<UpdateLogEntryMutation, UpdateLogEntryMutationVariables>(
        client,
        UpdateLogEntryDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateLogEntryMutation.getKey = () => ['updateLogEntry'];

useUpdateLogEntryMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateLogEntryMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateLogEntryMutation, UpdateLogEntryMutationVariables>(
    client,
    UpdateLogEntryDocument,
    variables,
    headers,
  );
