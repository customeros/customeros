// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

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
export type CreateOrganizationMutationVariables = Types.Exact<{
  input: Types.OrganizationInput;
}>;

export type CreateOrganizationMutation = {
  __typename?: 'Mutation';
  organization_Create: {
    __typename?: 'Organization';
    id: string;
    name: string;
  };
};

export const CreateOrganizationDocument = `
    mutation createOrganization($input: OrganizationInput!) {
  organization_Create(input: $input) {
    id
    name
  }
}
    `;
export const useCreateOrganizationMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CreateOrganizationMutation,
    TError,
    CreateOrganizationMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    CreateOrganizationMutation,
    TError,
    CreateOrganizationMutationVariables,
    TContext
  >(
    ['createOrganization'],
    (variables?: CreateOrganizationMutationVariables) =>
      fetcher<CreateOrganizationMutation, CreateOrganizationMutationVariables>(
        client,
        CreateOrganizationDocument,
        variables,
        headers,
      )(),
    options,
  );
useCreateOrganizationMutation.getKey = () => ['createOrganization'];

useCreateOrganizationMutation.fetcher = (
  client: GraphQLClient,
  variables: CreateOrganizationMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CreateOrganizationMutation, CreateOrganizationMutationVariables>(
    client,
    CreateOrganizationDocument,
    variables,
    headers,
  );
