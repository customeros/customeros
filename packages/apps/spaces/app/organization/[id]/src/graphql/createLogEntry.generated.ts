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
export type CreateLogEntryMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  logEntry: Types.LogEntryInput;
}>;

export type CreateLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_CreateForOrganization: string;
};

export const CreateLogEntryDocument = `
    mutation createLogEntry($organizationId: ID!, $logEntry: LogEntryInput!) {
  logEntry_CreateForOrganization(
    organizationId: $organizationId
    input: $logEntry
  )
}
    `;
export const useCreateLogEntryMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateLogEntryMutation,
    TError,
    CreateLogEntryMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    CreateLogEntryMutation,
    TError,
    CreateLogEntryMutationVariables,
    TContext
  >(
    ['createLogEntry'],
    (variables?: CreateLogEntryMutationVariables) =>
      fetcher<CreateLogEntryMutation, CreateLogEntryMutationVariables>(
        client,
        CreateLogEntryDocument,
        variables,
        headers,
      )(),
    options,
  );
useCreateLogEntryMutation.getKey = () => ['createLogEntry'];

useCreateLogEntryMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateLogEntryMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateLogEntryMutation, CreateLogEntryMutationVariables>(
    client,
    CreateLogEntryDocument,
    variables,
    headers,
  );
