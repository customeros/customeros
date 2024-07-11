import type { Transport } from '@store/transport';

import { gql } from 'graphql-request';

import type { ContractUpdateInput } from '@graphql/types';

class ContractService {
  private static instance: ContractService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): ContractService {
    if (!ContractService.instance) {
      ContractService.instance = new ContractService(transport);
    }

    return ContractService.instance;
  }

  async updateContract(
    payload: CONTRACT_UPDATE_PAYLOAD,
  ): Promise<{ id: string }> {
    return this.transport.graphql.request<
      { id: string },
      CONTRACT_UPDATE_PAYLOAD
    >(UPDATE_CONTRACT_MUTATION, payload);
  }

  async renewContract(
    payload: RENEW_CONTRACT_PAYLOAD,
  ): Promise<{ id: string }> {
    return this.transport.graphql.request<
      { id: string },
      RENEW_CONTRACT_PAYLOAD
    >(RENEW_CONTRACT_MUTATION, payload);
  }

  async addContractAttachment(
    payload: ADD_CONTRACT_ATTACHMENT_PAYLOAD,
  ): Promise<ADD_CONTRACT_ATTACHMENT_RESPONSE> {
    return this.transport.graphql.request<
      ADD_CONTRACT_ATTACHMENT_RESPONSE,
      ADD_CONTRACT_ATTACHMENT_PAYLOAD
    >(ADD_CONTRACT_ATTACHMENT_MUTATION, payload);
  }

  async removeContractAttachment(
    payload: REMOVE_CONTRACT_ATTACHMENT_PAYLOAD,
  ): Promise<REMOVE_CONTRACT_ATTACHMENT_RESPONSE> {
    return this.transport.graphql.request<
      REMOVE_CONTRACT_ATTACHMENT_RESPONSE,
      REMOVE_CONTRACT_ATTACHMENT_PAYLOAD
    >(REMOVE_CONTRACT_ATTACHMENT_MUTATION, payload);
  }
}

type CONTRACT_UPDATE_PAYLOAD = { input: ContractUpdateInput };
const UPDATE_CONTRACT_MUTATION = gql`
  mutation updateContract($input: ContractUpdateInput!) {
    contract_Update(input: $input) {
      id
    }
  }
`;

type RENEW_CONTRACT_PAYLOAD = { input: ContractUpdateInput };
const RENEW_CONTRACT_MUTATION = gql`
  mutation renewContract($input: ContractRenewalInput!) {
    contract_Renew(input: $input) {
      id
    }
  }
`;

type ADD_CONTRACT_ATTACHMENT_PAYLOAD = {
  contractId: string;
  attachmentId: string;
};
type ADD_CONTRACT_ATTACHMENT_RESPONSE = {
  attachments?: Array<{
    id: string;
    basePath: string;
    fileName: string;
    __typename?: 'Attachment';
  }> | null;
};
const ADD_CONTRACT_ATTACHMENT_MUTATION = gql`
  mutation addContractAttachment($contractId: ID!, $attachmentId: ID!) {
    contract_AddAttachment(
      contractId: $contractId
      attachmentId: $attachmentId
    ) {
      attachments {
        id
        basePath
        fileName
      }
    }
  }
`;
type REMOVE_CONTRACT_ATTACHMENT_PAYLOAD = {
  contractId: string;
  attachmentId: string;
};
type REMOVE_CONTRACT_ATTACHMENT_RESPONSE = {
  attachments?: Array<{
    id: string;
    basePath: string;
    fileName: string;
    __typename?: 'Attachment';
  }> | null;
};
const REMOVE_CONTRACT_ATTACHMENT_MUTATION = gql`
  mutation removeContractAttachment($contractId: ID!, $attachmentId: ID!) {
    contract_RemoveAttachment(
      contractId: $contractId
      attachmentId: $attachmentId
    ) {
      attachments {
        id
        basePath
        fileName
      }
    }
  }
`;

export { ContractService };
