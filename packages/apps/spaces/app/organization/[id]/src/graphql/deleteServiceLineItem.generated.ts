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
export type DeleteServiceLineItemMutationVariables = Types.Exact<{
  id: Types.Scalars['ID'];
}>;

export type DeleteServiceLineItemMutation = {
  __typename?: 'Mutation';
  serviceLineItem_Delete: {
    __typename?: 'DeleteResponse';
    accepted: boolean;
    completed: boolean;
  };
};

export const DeleteServiceLineItemDocument = `
    mutation deleteServiceLineItem($id: ID!) {
  serviceLineItem_Delete(id: $id) {
    accepted
    completed
  }
}
    `;
export const useDeleteServiceLineItemMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    DeleteServiceLineItemMutation,
    TError,
    DeleteServiceLineItemMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    DeleteServiceLineItemMutation,
    TError,
    DeleteServiceLineItemMutationVariables,
    TContext
  >(
    ['deleteServiceLineItem'],
    (variables?: DeleteServiceLineItemMutationVariables) =>
      fetcher<
        DeleteServiceLineItemMutation,
        DeleteServiceLineItemMutationVariables
      >(client, DeleteServiceLineItemDocument, variables, headers)(),
    options,
  );
useDeleteServiceLineItemMutation.getKey = () => ['deleteServiceLineItem'];

useDeleteServiceLineItemMutation.fetcher = (
  client: GraphQLClient,
  variables: DeleteServiceLineItemMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    DeleteServiceLineItemMutation,
    DeleteServiceLineItemMutationVariables
  >(client, DeleteServiceLineItemDocument, variables, headers);
