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
export type CloseServiceLineItemMutationVariables = Types.Exact<{
  input: Types.ServiceLineItemCloseInput;
}>;

export type CloseServiceLineItemMutation = {
  __typename?: 'Mutation';
  serviceLineItem_Close: string;
};

export const CloseServiceLineItemDocument = `
    mutation CloseServiceLineItem($input: ServiceLineItemCloseInput!) {
  serviceLineItem_Close(input: $input)
}
    `;
export const useCloseServiceLineItemMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    CloseServiceLineItemMutation,
    TError,
    CloseServiceLineItemMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    CloseServiceLineItemMutation,
    TError,
    CloseServiceLineItemMutationVariables,
    TContext
  >(
    ['CloseServiceLineItem'],
    (variables?: CloseServiceLineItemMutationVariables) =>
      fetcher<
        CloseServiceLineItemMutation,
        CloseServiceLineItemMutationVariables
      >(client, CloseServiceLineItemDocument, variables, headers)(),
    options,
  );
useCloseServiceLineItemMutation.getKey = () => ['CloseServiceLineItem'];

useCloseServiceLineItemMutation.fetcher = (
  client: GraphQLClient,
  variables: CloseServiceLineItemMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<CloseServiceLineItemMutation, CloseServiceLineItemMutationVariables>(
    client,
    CloseServiceLineItemDocument,
    variables,
    headers,
  );
