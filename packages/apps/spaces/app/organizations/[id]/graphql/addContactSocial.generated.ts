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
export type AddContactSocialMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
  input: Types.SocialInput;
}>;

export type AddContactSocialMutation = {
  __typename?: 'Mutation';
  contact_AddSocial: { __typename?: 'Social'; id: string; url: string };
};

export const AddContactSocialDocument = `
    mutation addContactSocial($contactId: ID!, $input: SocialInput!) {
  contact_AddSocial(contactId: $contactId, input: $input) {
    id
    url
  }
}
    `;
export const useAddContactSocialMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddContactSocialMutation,
    TError,
    AddContactSocialMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddContactSocialMutation,
    TError,
    AddContactSocialMutationVariables,
    TContext
  >(
    ['addContactSocial'],
    (variables?: AddContactSocialMutationVariables) =>
      fetcher<AddContactSocialMutation, AddContactSocialMutationVariables>(
        client,
        AddContactSocialDocument,
        variables,
        headers,
      )(),
    options,
  );
useAddContactSocialMutation.getKey = () => ['addContactSocial'];

useAddContactSocialMutation.fetcher = (
  client: GraphQLClient,
  variables: AddContactSocialMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<AddContactSocialMutation, AddContactSocialMutationVariables>(
    client,
    AddContactSocialDocument,
    variables,
    headers,
  );
