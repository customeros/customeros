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
export type UpdateTenantSettingsMutationVariables = Types.Exact<{
  input: Types.TenantSettingsInput;
}>;


export type UpdateTenantSettingsMutation = { __typename?: 'Mutation', tenant_UpdateSettings: { __typename?: 'TenantSettings', logoUrl: string, billingEnabled: boolean, baseCurrency?: Types.Currency | null } };



export const UpdateTenantSettingsDocument = `
    mutation UpdateTenantSettings($input: TenantSettingsInput!) {
  tenant_UpdateSettings(input: $input) {
    logoUrl
    billingEnabled
    baseCurrency
  }
}
    `;

export const useUpdateTenantSettingsMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<UpdateTenantSettingsMutation, TError, UpdateTenantSettingsMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<UpdateTenantSettingsMutation, TError, UpdateTenantSettingsMutationVariables, TContext>(
      {
    mutationKey: ['UpdateTenantSettings'],
    mutationFn: (variables?: UpdateTenantSettingsMutationVariables) => fetcher<UpdateTenantSettingsMutation, UpdateTenantSettingsMutationVariables>(client, UpdateTenantSettingsDocument, variables, headers)(),
    ...options
  }
    )};

useUpdateTenantSettingsMutation.getKey = () => ['UpdateTenantSettings'];


useUpdateTenantSettingsMutation.fetcher = (client: GraphQLClient, variables: UpdateTenantSettingsMutationVariables, headers?: RequestInit['headers']) => fetcher<UpdateTenantSettingsMutation, UpdateTenantSettingsMutationVariables>(client, UpdateTenantSettingsDocument, variables, headers);
