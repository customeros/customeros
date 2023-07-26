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
export type RemoveContactPhoneNumberMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID'];
  id: Types.Scalars['ID'];
}>;

export type RemoveContactPhoneNumberMutation = {
  __typename?: 'Mutation';
  phoneNumberRemoveFromContactById: { __typename?: 'Result'; result: boolean };
};

export const RemoveContactPhoneNumberDocument = `
    mutation removeContactPhoneNumber($contactId: ID!, $id: ID!) {
  phoneNumberRemoveFromContactById(contactId: $contactId, id: $id) {
    result
  }
}
    `;
export const useRemoveContactPhoneNumberMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveContactPhoneNumberMutation,
    TError,
    RemoveContactPhoneNumberMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    RemoveContactPhoneNumberMutation,
    TError,
    RemoveContactPhoneNumberMutationVariables,
    TContext
  >(
    ['removeContactPhoneNumber'],
    (variables?: RemoveContactPhoneNumberMutationVariables) =>
      fetcher<
        RemoveContactPhoneNumberMutation,
        RemoveContactPhoneNumberMutationVariables
      >(client, RemoveContactPhoneNumberDocument, variables, headers)(),
    options,
  );
useRemoveContactPhoneNumberMutation.getKey = () => ['removeContactPhoneNumber'];

useRemoveContactPhoneNumberMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveContactPhoneNumberMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    RemoveContactPhoneNumberMutation,
    RemoveContactPhoneNumberMutationVariables
  >(client, RemoveContactPhoneNumberDocument, variables, headers);
