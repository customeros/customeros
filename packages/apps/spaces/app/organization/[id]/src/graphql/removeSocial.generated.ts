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
export type RemoveSocialMutationVariables = Types.Exact<{
  socialId: Types.Scalars['ID'];
}>;

export type RemoveSocialMutation = {
  __typename?: 'Mutation';
  social_Remove: { __typename?: 'Result'; result: boolean };
};

export const RemoveSocialDocument = `
    mutation removeSocial($socialId: ID!) {
  social_Remove(socialId: $socialId) {
    result
  }
}
    `;
export const useRemoveSocialMutation = <TError = unknown, TContext = unknown>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveSocialMutation,
    TError,
    RemoveSocialMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    RemoveSocialMutation,
    TError,
    RemoveSocialMutationVariables,
    TContext
  >(
    ['removeSocial'],
    (variables?: RemoveSocialMutationVariables) =>
      fetcher<RemoveSocialMutation, RemoveSocialMutationVariables>(
        client,
        RemoveSocialDocument,
        variables,
        headers,
      )(),
    options,
  );
useRemoveSocialMutation.getKey = () => ['removeSocial'];

useRemoveSocialMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveSocialMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<RemoveSocialMutation, RemoveSocialMutationVariables>(
    client,
    RemoveSocialDocument,
    variables,
    headers,
  );
