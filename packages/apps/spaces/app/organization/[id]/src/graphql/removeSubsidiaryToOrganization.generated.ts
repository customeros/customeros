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
export type RemoveSubsidiaryToOrganizationMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID'];
  subsidiaryId: Types.Scalars['ID'];
}>;

export type RemoveSubsidiaryToOrganizationMutation = {
  __typename?: 'Mutation';
  organization_RemoveSubsidiary: {
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

export const RemoveSubsidiaryToOrganizationDocument = `
    mutation removeSubsidiaryToOrganization($organizationId: ID!, $subsidiaryId: ID!) {
  organization_RemoveSubsidiary(
    organizationId: $organizationId
    subsidiaryId: $subsidiaryId
  ) {
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
export const useRemoveSubsidiaryToOrganizationMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    RemoveSubsidiaryToOrganizationMutation,
    TError,
    RemoveSubsidiaryToOrganizationMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) =>
  useMutation<
    RemoveSubsidiaryToOrganizationMutation,
    TError,
    RemoveSubsidiaryToOrganizationMutationVariables,
    TContext
  >(
    ['removeSubsidiaryToOrganization'],
    (variables?: RemoveSubsidiaryToOrganizationMutationVariables) =>
      fetcher<
        RemoveSubsidiaryToOrganizationMutation,
        RemoveSubsidiaryToOrganizationMutationVariables
      >(client, RemoveSubsidiaryToOrganizationDocument, variables, headers)(),
    options,
  );
useRemoveSubsidiaryToOrganizationMutation.getKey = () => [
  'removeSubsidiaryToOrganization',
];

useRemoveSubsidiaryToOrganizationMutation.fetcher = (
  client: GraphQLClient,
  variables: RemoveSubsidiaryToOrganizationMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    RemoveSubsidiaryToOrganizationMutation,
    RemoveSubsidiaryToOrganizationMutationVariables
  >(client, RemoveSubsidiaryToOrganizationDocument, variables, headers);
