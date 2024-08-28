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
export type CreateBillingProfileMutationVariables = Types.Exact<{
  input: Types.TenantBillingProfileInput;
}>;


export type CreateBillingProfileMutation = { __typename?: 'Mutation', tenant_AddBillingProfile: { __typename?: 'TenantBillingProfile', id: string } };



export const CreateBillingProfileDocument = `
    mutation CreateBillingProfile($input: TenantBillingProfileInput!) {
  tenant_AddBillingProfile(input: $input) {
    id
  }
}
    `;

export const useCreateBillingProfileMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<CreateBillingProfileMutation, TError, CreateBillingProfileMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<CreateBillingProfileMutation, TError, CreateBillingProfileMutationVariables, TContext>(
      {
    mutationKey: ['CreateBillingProfile'],
    mutationFn: (variables?: CreateBillingProfileMutationVariables) => fetcher<CreateBillingProfileMutation, CreateBillingProfileMutationVariables>(client, CreateBillingProfileDocument, variables, headers)(),
    ...options
  }
    )};

useCreateBillingProfileMutation.getKey = () => ['CreateBillingProfile'];


useCreateBillingProfileMutation.fetcher = (client: GraphQLClient, variables: CreateBillingProfileMutationVariables, headers?: RequestInit['headers']) => fetcher<CreateBillingProfileMutation, CreateBillingProfileMutationVariables>(client, CreateBillingProfileDocument, variables, headers);
