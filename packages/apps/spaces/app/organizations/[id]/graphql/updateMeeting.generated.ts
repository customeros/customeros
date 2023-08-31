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
export type UpdateMeetingMutationVariables = Types.Exact<{
  meetingId: Types.Scalars['ID'];
  meeting: Types.MeetingUpdateInput;
}>;

export type UpdateMeetingMutation = {
  __typename?: 'Mutation';
  meeting_Update: { __typename?: 'Meeting'; id: string };
};

export const UpdateMeetingDocument = `
    mutation updateMeeting($meetingId: ID!, $meeting: MeetingUpdateInput!) {
  meeting_Update(meetingId: $meetingId, meeting: $meeting) {
    id
  }
}
    `;
export const useUpdateMeetingMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateMeetingMutation,
    TError,
    UpdateMeetingMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateMeetingMutation,
    TError,
    UpdateMeetingMutationVariables,
    TContext
  >(
    ['updateMeeting'],
    (variables?: UpdateMeetingMutationVariables) =>
      fetcher<UpdateMeetingMutation, UpdateMeetingMutationVariables>(
        client,
        UpdateMeetingDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateMeetingMutation.getKey = () => ['updateMeeting'];

useUpdateMeetingMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateMeetingMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateMeetingMutation, UpdateMeetingMutationVariables>(
    client,
    UpdateMeetingDocument,
    variables,
    headers,
  );
