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
export type AddOrganizationToContactMutationVariables = Types.Exact<{
  input: Types.ContactOrganizationInput;
}>;

export type AddOrganizationToContactMutation = {
  __typename?: 'Mutation';
  contact_AddOrganizationById: { __typename?: 'Contact'; id: string };
};

export const AddOrganizationToContactDocument = `
    mutation addOrganizationToContact($input: ContactOrganizationInput!) {
  contact_AddOrganizationById(input: $input) {
    id
  }
}
    `;
export const useAddOrganizationToContactMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddOrganizationToContactMutation,
    TError,
    AddOrganizationToContactMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddOrganizationToContactMutation,
    TError,
    AddOrganizationToContactMutationVariables,
    TContext
  >(
    ['addOrganizationToContact'],
    (variables?: AddOrganizationToContactMutationVariables) =>
      fetcher<
        AddOrganizationToContactMutation,
        AddOrganizationToContactMutationVariables
      >(client, AddOrganizationToContactDocument, variables, headers)(),
    options,
  );
useAddOrganizationToContactMutation.getKey = () => ['addOrganizationToContact'];

useAddOrganizationToContactMutation.fetcher = (
  client: GraphQLClient,
  variables: AddOrganizationToContactMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    AddOrganizationToContactMutation,
    AddOrganizationToContactMutationVariables
  >(client, AddOrganizationToContactDocument, variables, headers);
