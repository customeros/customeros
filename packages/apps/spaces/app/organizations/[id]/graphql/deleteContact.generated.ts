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
export type DeleteContactMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
}>;

export type DeleteContactMutation = {
  __typename?: 'Mutation';
  contact_HardDelete: { __typename?: 'Result'; result: boolean };
};

export const DeleteContactDocument = `
    mutation deleteContact($contactId: ID!) {
  contact_HardDelete(contactId: $contactId) {
    result
  }
}
    `;
export const useDeleteContactMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    DeleteContactMutation,
    TError,
    DeleteContactMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    DeleteContactMutation,
    TError,
    DeleteContactMutationVariables,
    TContext
  >(
    ['deleteContact'],
    (variables?: DeleteContactMutationVariables) =>
      fetcher<DeleteContactMutation, DeleteContactMutationVariables>(
        client,
        DeleteContactDocument,
        variables,
        headers,
      )(),
    options,
  );
useDeleteContactMutation.getKey = () => ['deleteContact'];

useDeleteContactMutation.fetcher = (
  client: GraphQLClient,
  variables: DeleteContactMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<DeleteContactMutation, DeleteContactMutationVariables>(
    client,
    DeleteContactDocument,
    variables,
    headers,
  );
