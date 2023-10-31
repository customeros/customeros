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
export type AddSubsidiaryToOrganizationMutationVariables = Types.Exact<{
  input: Types.LinkOrganizationsInput;
}>;

export type AddSubsidiaryToOrganizationMutation = {
  __typename?: 'Mutation';
  organization_AddSubsidiary: {
    __typename?: 'Organization';
    id: string;
    subsidiaries: Array<{
      __typename?: 'LinkedOrganization';
      organization: {
        __typename?: 'Organization';
        id: string;
        name: string;
        locations: Array<{
          __typename?: 'Location';
          id: string;
          address?: string | null;
        }>;
      };
    }>;
  };
};

export const AddSubsidiaryToOrganizationDocument = `
    mutation addSubsidiaryToOrganization($input: LinkOrganizationsInput!) {
  organization_AddSubsidiary(input: $input) {
    id
    subsidiaries {
      organization {
        id
        name
        locations {
          id
          address
        }
      }
    }
  }
}
    `;
export const useAddSubsidiaryToOrganizationMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddSubsidiaryToOrganizationMutation,
    TError,
    AddSubsidiaryToOrganizationMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    AddSubsidiaryToOrganizationMutation,
    TError,
    AddSubsidiaryToOrganizationMutationVariables,
    TContext
  >(
    ['addSubsidiaryToOrganization'],
    (variables?: AddSubsidiaryToOrganizationMutationVariables) =>
      fetcher<
        AddSubsidiaryToOrganizationMutation,
        AddSubsidiaryToOrganizationMutationVariables
      >(client, AddSubsidiaryToOrganizationDocument, variables, headers)(),
    options,
  );
useAddSubsidiaryToOrganizationMutation.getKey = () => [
  'addSubsidiaryToOrganization',
];

useAddSubsidiaryToOrganizationMutation.fetcher = (
  client: GraphQLClient,
  variables: AddSubsidiaryToOrganizationMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    AddSubsidiaryToOrganizationMutation,
    AddSubsidiaryToOrganizationMutationVariables
  >(client, AddSubsidiaryToOrganizationDocument, variables, headers);
