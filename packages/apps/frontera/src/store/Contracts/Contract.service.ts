import type { Transport } from '@store/transport';

import { gql } from 'graphql-request';

import {
  Contract,
  Pagination,
  ContractInput,
  ContractUpdateInput,
} from '@graphql/types';

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

  async createContract(
    payload: CREATE_CONTRACT_PAYLOAD,
  ): Promise<CREATE_CONTRACT_RESPONSE> {
    return this.transport.graphql.request<
      CREATE_CONTRACT_RESPONSE,
      CREATE_CONTRACT_PAYLOAD
    >(CREATE_CONTRACT_MUTATION, payload);
  }

  async getContracts(
    payload: CONTRACTS_QUERY_PAYLOAD,
  ): Promise<CONTRACTS_QUERY_RESPONSE> {
    return this.transport.graphql.request<
      CONTRACTS_QUERY_RESPONSE,
      CONTRACTS_QUERY_PAYLOAD
    >(CONTRACTS_QUERY, payload);
  }
}

type CONTRACTS_QUERY_RESPONSE = {
  contracts: {
    totalPages: number;
    content: Contract[];
    totalElements: number;
    totalAvailable: number;
  };
};
type CONTRACTS_QUERY_PAYLOAD = {
  pagination: Pagination;
};

const CONTRACTS_QUERY = gql`
  query getContracts($pagination: Pagination!) {
    contracts(pagination: $pagination) {
      totalPages
      totalElements
      totalAvailable
      content {
        metadata {
          id
          created
          source
          lastUpdated
        }

        contractName
        serviceStarted
        contractSigned
        contractEnded
        contractStatus
        committedPeriodInMonths
        approved
        ltv

        contractUrl
        billingCycle
        billingEnabled
        currency
        invoiceEmail
        autoRenew

        billingDetails {
          nextInvoicing
          postalCode
          country
          locality
          addressLine1
          addressLine2
          invoiceNote
          organizationLegalName
          billingCycle
          payAutomatically
          billingCycleInMonths
          invoicingStarted
          region
          dueDays
          billingEmail
          billingEmailCC
          billingEmailBCC
          check
          payOnline
          canPayWithDirectDebit
          canPayWithBankTransfer
          canPayWithCard
        }
        upcomingInvoices {
          metadata {
            id
          }
        }
        opportunities {
          metadata {
            id
            created
            lastUpdated
            source
            sourceOfTruth
            appSource
          }
          name
          amount
          maxAmount
          internalType
          externalType
          internalStage
          externalStage
          estimatedClosedAt
          generalNotes
          nextSteps
          renewedAt
          renewalApproved
          renewalLikelihood
          renewalUpdatedByUserId
          renewalUpdatedByUserAt
          renewalAdjustedRate
          comments
          organization {
            metadata {
              id
              created
              lastUpdated
              sourceOfTruth
            }
          }
          createdBy {
            id
            firstName
            lastName
            name
          }
          owner {
            id
            firstName
            lastName
            name
          }
          externalLinks {
            externalUrl
            externalId
          }
          id
          createdAt
          updatedAt
          source
          appSource
        }
        contractLineItems {
          metadata {
            id
            created
            lastUpdated
            source
            appSource
            sourceOfTruth
          }
          paused
          description
          billingCycle
          price
          quantity
          comments
          serviceEnded
          parentId
          serviceStarted
          tax {
            salesTax
            vat
            taxRate
          }
        }
        attachments {
          id
          createdAt
          basePath
          cdnUrl
          fileName
          mimeType
          size
          source
          sourceOfTruth
          appSource
        }
      }
    }
  }
`;

type CREATE_CONTRACT_PAYLOAD = {
  input: ContractInput;
};
type CREATE_CONTRACT_RESPONSE = {
  contract_Create: {
    metadata: {
      id: string;
    };
  };
};
const CREATE_CONTRACT_MUTATION = gql`
  mutation createContract($input: ContractInput!) {
    contract_Create(input: $input) {
      metadata {
        id
      }
    }
  }
`;

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
