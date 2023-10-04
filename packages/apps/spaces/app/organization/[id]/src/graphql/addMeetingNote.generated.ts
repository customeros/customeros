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
export type AddMeetingNoteMutationVariables = Types.Exact<{
  meetingId: Types.Scalars['ID'];
  note: Types.NoteInput;
}>;

export type AddMeetingNoteMutation = {
  __typename?: 'Mutation';
  meeting_AddNote: { __typename?: 'Meeting'; id: string };
};

export const AddMeetingNoteDocument = `
    mutation addMeetingNote($meetingId: ID!, $note: NoteInput!) {
  meeting_AddNote(meetingId: $meetingId, note: $note) {
    id
  }
}
    `;
export const useAddMeetingNoteMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddMeetingNoteMutation,
    TError,
    AddMeetingNoteMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddMeetingNoteMutation,
    TError,
    AddMeetingNoteMutationVariables,
    TContext
  >(
    ['addMeetingNote'],
    (variables?: AddMeetingNoteMutationVariables) =>
      fetcher<AddMeetingNoteMutation, AddMeetingNoteMutationVariables>(
        client,
        AddMeetingNoteDocument,
        variables,
        headers,
      )(),
    options,
  );
useAddMeetingNoteMutation.getKey = () => ['addMeetingNote'];

useAddMeetingNoteMutation.fetcher = (
  client: GraphQLClient,
  variables: AddMeetingNoteMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<AddMeetingNoteMutation, AddMeetingNoteMutationVariables>(
    client,
    AddMeetingNoteDocument,
    variables,
    headers,
  );
