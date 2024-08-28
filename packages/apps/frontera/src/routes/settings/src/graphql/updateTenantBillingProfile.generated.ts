// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useMutation, UseMutationOptions } from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(client: GraphQLClient, query: string, variables?: TVariables, requestHeaders?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request({
    document: query,
    variables,
    requestHeaders
  });
}
export type TenantUpdateBillingProfileMutationVariables = Types.Exact<{
  input: Types.TenantBillingProfileUpdateInput;
}>;


export type TenantUpdateBillingProfileMutation = { __typename?: 'Mutation', tenant_UpdateBillingProfile: { __typename?: 'TenantBillingProfile', id: string } };



export const TenantUpdateBillingProfileDocument = `
    mutation TenantUpdateBillingProfile($input: TenantBillingProfileUpdateInput!) {
  tenant_UpdateBillingProfile(input: $input) {
    id
  }
}
    `;

export const useTenantUpdateBillingProfileMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<TenantUpdateBillingProfileMutation, TError, TenantUpdateBillingProfileMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<TenantUpdateBillingProfileMutation, TError, TenantUpdateBillingProfileMutationVariables, TContext>(
      {
    mutationKey: ['TenantUpdateBillingProfile'],
    mutationFn: (variables?: TenantUpdateBillingProfileMutationVariables) => fetcher<TenantUpdateBillingProfileMutation, TenantUpdateBillingProfileMutationVariables>(client, TenantUpdateBillingProfileDocument, variables, headers)(),
    ...options
  }
    )};

useTenantUpdateBillingProfileMutation.getKey = () => ['TenantUpdateBillingProfile'];


useTenantUpdateBillingProfileMutation.fetcher = (client: GraphQLClient, variables: TenantUpdateBillingProfileMutationVariables, headers?: RequestInit['headers']) => fetcher<TenantUpdateBillingProfileMutation, TenantUpdateBillingProfileMutationVariables>(client, TenantUpdateBillingProfileDocument, variables, headers);
