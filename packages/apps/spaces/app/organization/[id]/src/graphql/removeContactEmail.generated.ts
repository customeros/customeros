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
export type RemoveContactEmailMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
  email: Types.Scalars['String'];
}>;

export type RemoveContactEmailMutation = {
  __typename?: 'Mutation';
  emailRemoveFromContact: { __typename?: 'Result'; result: boolean };
};

export const RemoveContactEmailDocument = `
    mutation removeContactEmail($contactId: ID!, $email: String!) {
  emailRemoveFromContact(contactId: $contactId, email: $email) {
    result
  }
}
    `;
export const useRemoveContactEmailMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveContactEmailMutation,
    TError,
    RemoveContactEmailMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    RemoveContactEmailMutation,
    TError,
    RemoveContactEmailMutationVariables,
    TContext
  >(
    ['removeContactEmail'],
    (variables?: RemoveContactEmailMutationVariables) =>
      fetcher<RemoveContactEmailMutation, RemoveContactEmailMutationVariables>(
        client,
        RemoveContactEmailDocument,
        variables,
        headers,
      )(),
    options,
  );
useRemoveContactEmailMutation.getKey = () => ['removeContactEmail'];

useRemoveContactEmailMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveContactEmailMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<RemoveContactEmailMutation, RemoveContactEmailMutationVariables>(
    client,
    RemoveContactEmailDocument,
    variables,
    headers,
  );
