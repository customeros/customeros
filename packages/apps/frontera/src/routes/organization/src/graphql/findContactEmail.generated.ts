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
export type FindContactEmailMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  organizationId: Types.Scalars['ID']['input'];
}>;


export type FindContactEmailMutation = { __typename?: 'Mutation', contact_FindEmail: string };



export const FindContactEmailDocument = `
    mutation findContactEmail($contactId: ID!, $organizationId: ID!) {
  contact_FindEmail(contactId: $contactId, organizationId: $organizationId)
}
    `;

export const useFindContactEmailMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<FindContactEmailMutation, TError, FindContactEmailMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<FindContactEmailMutation, TError, FindContactEmailMutationVariables, TContext>(
      {
    mutationKey: ['findContactEmail'],
    mutationFn: (variables?: FindContactEmailMutationVariables) => fetcher<FindContactEmailMutation, FindContactEmailMutationVariables>(client, FindContactEmailDocument, variables, headers)(),
    ...options
  }
    )};

useFindContactEmailMutation.getKey = () => ['findContactEmail'];


useFindContactEmailMutation.fetcher = (client: GraphQLClient, variables: FindContactEmailMutationVariables, headers?: RequestInit['headers']) => fetcher<FindContactEmailMutation, FindContactEmailMutationVariables>(client, FindContactEmailDocument, variables, headers);
