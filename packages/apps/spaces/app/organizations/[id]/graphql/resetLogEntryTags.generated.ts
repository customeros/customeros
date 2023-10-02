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
export type ResetLogEntryTagsMutationVariables = Types.Exact<{
  id: Types.Scalars['ID'];
  input?: Types.InputMaybe<
    Array<Types.TagIdOrNameInput> | Types.TagIdOrNameInput
  >;
}>;

export type ResetLogEntryTagsMutation = {
  __typename?: 'Mutation';
  logEntry_ResetTags: string;
};

export const ResetLogEntryTagsDocument = `
    mutation resetLogEntryTags($id: ID!, $input: [TagIdOrNameInput!]) {
  logEntry_ResetTags(id: $id, input: $input)
}
    `;
export const useResetLogEntryTagsMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    ResetLogEntryTagsMutation,
    TError,
    ResetLogEntryTagsMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    ResetLogEntryTagsMutation,
    TError,
    ResetLogEntryTagsMutationVariables,
    TContext
  >(
    ['resetLogEntryTags'],
    (variables?: ResetLogEntryTagsMutationVariables) =>
      fetcher<ResetLogEntryTagsMutation, ResetLogEntryTagsMutationVariables>(
        client,
        ResetLogEntryTagsDocument,
        variables,
        headers,
      )(),
    options,
  );
useResetLogEntryTagsMutation.getKey = () => ['resetLogEntryTags'];

useResetLogEntryTagsMutation.fetcher = (
  client: GraphQLClient,
  variables: ResetLogEntryTagsMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<ResetLogEntryTagsMutation, ResetLogEntryTagsMutationVariables>(
    client,
    ResetLogEntryTagsDocument,
    variables,
    headers,
  );
