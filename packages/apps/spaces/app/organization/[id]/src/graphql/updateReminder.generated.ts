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
export type UpdateReminderMutationVariables = Types.Exact<{
  input: Types.ReminderUpdateInput;
}>;

export type UpdateReminderMutation = {
  __typename?: 'Mutation';
  reminder_Update: {
    __typename?: 'Reminder';
    content?: string | null;
    dueDate?: any | null;
    dismissed?: boolean | null;
    metadata: { __typename?: 'Metadata'; id: string };
    owner?: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
      name?: string | null;
    } | null;
  };
};

export const UpdateReminderDocument = `
    mutation updateReminder($input: ReminderUpdateInput!) {
  reminder_Update(input: $input) {
    metadata {
      id
    }
    content
    owner {
      id
      firstName
      lastName
      name
    }
    dueDate
    dismissed
  }
}
    `;

export const useUpdateReminderMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateReminderMutation,
    TError,
    UpdateReminderMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    UpdateReminderMutation,
    TError,
    UpdateReminderMutationVariables,
    TContext
  >({
    mutationKey: ['updateReminder'],
    mutationFn: (variables?: UpdateReminderMutationVariables) =>
      fetcher<UpdateReminderMutation, UpdateReminderMutationVariables>(
        client,
        UpdateReminderDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useUpdateReminderMutation.getKey = () => ['updateReminder'];

useUpdateReminderMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateReminderMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateReminderMutation, UpdateReminderMutationVariables>(
    client,
    UpdateReminderDocument,
    variables,
    headers,
  );
