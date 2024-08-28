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
export type RemoveContractAttachmentMutationVariables = Types.Exact<{
  contractId: Types.Scalars['ID']['input'];
  attachmentId: Types.Scalars['ID']['input'];
}>;


export type RemoveContractAttachmentMutation = { __typename?: 'Mutation', contract_RemoveAttachment: { __typename?: 'Contract', attachments?: Array<{ __typename?: 'Attachment', id: string, basePath: string, fileName: string }> | null } };



export const RemoveContractAttachmentDocument = `
    mutation removeContractAttachment($contractId: ID!, $attachmentId: ID!) {
  contract_RemoveAttachment(contractId: $contractId, attachmentId: $attachmentId) {
    attachments {
      id
      basePath
      fileName
    }
  }
}
    `;

export const useRemoveContractAttachmentMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<RemoveContractAttachmentMutation, TError, RemoveContractAttachmentMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) => {
    
    return useMutation<RemoveContractAttachmentMutation, TError, RemoveContractAttachmentMutationVariables, TContext>(
      {
    mutationKey: ['removeContractAttachment'],
    mutationFn: (variables?: RemoveContractAttachmentMutationVariables) => fetcher<RemoveContractAttachmentMutation, RemoveContractAttachmentMutationVariables>(client, RemoveContractAttachmentDocument, variables, headers)(),
    ...options
  }
    )};

useRemoveContractAttachmentMutation.getKey = () => ['removeContractAttachment'];


useRemoveContractAttachmentMutation.fetcher = (client: GraphQLClient, variables: RemoveContractAttachmentMutationVariables, headers?: RequestInit['headers']) => fetcher<RemoveContractAttachmentMutation, RemoveContractAttachmentMutationVariables>(client, RemoveContractAttachmentDocument, variables, headers);
