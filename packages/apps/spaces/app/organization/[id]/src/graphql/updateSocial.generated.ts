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
export type UpdateSocialMutationVariables = Types.Exact<{
  input: Types.SocialUpdateInput;
}>;

export type UpdateSocialMutation = {
  __typename?: 'Mutation';
  social_Update: { __typename?: 'Social'; id: string; url: string };
};

export const UpdateSocialDocument = `
    mutation updateSocial($input: SocialUpdateInput!) {
  social_Update(input: $input) {
    id
    url
  }
}
    `;
export const useUpdateSocialMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    UpdateSocialMutation,
    TError,
    UpdateSocialMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    UpdateSocialMutation,
    TError,
    UpdateSocialMutationVariables,
    TContext
  >(
    ['updateSocial'],
    (variables?: UpdateSocialMutationVariables) =>
      fetcher<UpdateSocialMutation, UpdateSocialMutationVariables>(
        client,
        UpdateSocialDocument,
        variables,
        headers,
      )(),
    options,
  );
useUpdateSocialMutation.getKey = () => ['updateSocial'];

useUpdateSocialMutation.fetcher = (
  client: GraphQLClient,
  variables: UpdateSocialMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<UpdateSocialMutation, UpdateSocialMutationVariables>(
    client,
    UpdateSocialDocument,
    variables,
    headers,
  );
