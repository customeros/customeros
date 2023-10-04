// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../types/__generated__/graphql.types';

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
export type AddSocialMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  input: Types.SocialInput;
}>;

export type AddSocialMutation = {
  __typename?: 'Mutation';
  organization_AddSocial: { __typename?: 'Social'; id: string; url: string };
};

export const AddSocialDocument = `
    mutation addSocial($organizationId: ID!, $input: SocialInput!) {
  organization_AddSocial(organizationId: $organizationId, input: $input) {
    id
    url
  }
}
    `;
export const useAddSocialMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddSocialMutation,
    TError,
    AddSocialMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<AddSocialMutation, TError, AddSocialMutationVariables, TContext>(
    ['addSocial'],
    (variables?: AddSocialMutationVariables) =>
      fetcher<AddSocialMutation, AddSocialMutationVariables>(
        client,
        AddSocialDocument,
        variables,
        headers,
      )(),
    options,
  );
useAddSocialMutation.getKey = () => ['addSocial'];

useAddSocialMutation.fetcher = (
  client: GraphQLClient,
  variables: AddSocialMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<AddSocialMutation, AddSocialMutationVariables>(
    client,
    AddSocialDocument,
    variables,
    headers,
  );
