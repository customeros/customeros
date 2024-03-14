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
export type CreateReminderMutationVariables = Types.Exact<{
  input: Types.ReminderInput;
}>;

export type CreateReminderMutation = {
  __typename?: 'Mutation';
  reminder_Create: {
    __typename?: 'Reminder';
    content: string;
    dueDate: any;
    dismissed: boolean;
    metadata: { __typename?: 'Metadata'; id: string };
    owner: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
      name?: string | null;
    };
  };
};

export const CreateReminderDocument = `
    mutation createReminder($input: ReminderInput!) {
  reminder_Create(input: $input) {
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

export const useCreateReminderMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateReminderMutation,
    TError,
    CreateReminderMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    CreateReminderMutation,
    TError,
    CreateReminderMutationVariables,
    TContext
  >({
    mutationKey: ['createReminder'],
    mutationFn: (variables?: CreateReminderMutationVariables) =>
      fetcher<CreateReminderMutation, CreateReminderMutationVariables>(
        client,
        CreateReminderDocument,
        variables,
        headers,
      )(),
    ...options,
  });
};

useCreateReminderMutation.getKey = () => ['createReminder'];

useCreateReminderMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateReminderMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateReminderMutation, CreateReminderMutationVariables>(
    client,
    CreateReminderDocument,
    variables,
    headers,
  );
