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
export type AddContractAttachmentMutationVariables = Types.Exact<{
  contractId: Types.Scalars['ID']['input'];
  attachmentId: Types.Scalars['ID']['input'];
}>;

export type AddContractAttachmentMutation = {
  __typename?: 'Mutation';
  contract_AddAttachment: {
    __typename?: 'Contract';
    attachments?: Array<{
      __typename?: 'Attachment';
      id: string;
      basePath: string;
      fileName: string;
    }> | null;
  };
};

export const AddContractAttachmentDocument = `
    mutation addContractAttachment($contractId: ID!, $attachmentId: ID!) {
  contract_AddAttachment(contractId: $contractId, attachmentId: $attachmentId) {
    attachments {
      id
      basePath
      fileName
    }
  }
}
    `;

export const useAddContractAttachmentMutation = <
  TError = unknown,
  TContext = unknown,
>(
  client: GraphQLClient,
  options?: UseMutationOptions<
    AddContractAttachmentMutation,
    TError,
    AddContractAttachmentMutationVariables,
    TContext
  >,
  headers?: RequestInit['headers'],
) => {
  return useMutation<
    AddContractAttachmentMutation,
    TError,
    AddContractAttachmentMutationVariables,
    TContext
  >({
    mutationKey: ['addContractAttachment'],
    mutationFn: (variables?: AddContractAttachmentMutationVariables) =>
      fetcher<
        AddContractAttachmentMutation,
        AddContractAttachmentMutationVariables
      >(client, AddContractAttachmentDocument, variables, headers)(),
    ...options,
  });
};

useAddContractAttachmentMutation.getKey = () => ['addContractAttachment'];

useAddContractAttachmentMutation.fetcher = (
  client: GraphQLClient,
  variables: AddContractAttachmentMutationVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<
    AddContractAttachmentMutation,
    AddContractAttachmentMutationVariables
  >(client, AddContractAttachmentDocument, variables, headers);
