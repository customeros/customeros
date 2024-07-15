import type { Transport } from '@store/transport';

import { gql } from 'graphql-request';

import { ServiceLineItemUpdateInput } from '@graphql/types';

class ContractLineItemService {
  private static instance: ContractLineItemService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): ContractLineItemService {
    if (!ContractLineItemService.instance) {
      ContractLineItemService.instance = new ContractLineItemService(transport);
    }

    return ContractLineItemService.instance;
  }

  async updateContractLineItem(
    payload: CONTRACT_LINE_ITEM_UPDATE_PAYLOAD,
  ): Promise<CONTRACT_LINE_ITEM_UPDATE_RESPONSE> {
    return this.transport.graphql.request<
      CONTRACT_LINE_ITEM_UPDATE_RESPONSE,
      CONTRACT_LINE_ITEM_UPDATE_PAYLOAD
    >(CONTRACT_LINE_ITEM_UPDATE_MUTATION, payload);
  }
}

type CONTRACT_LINE_ITEM_UPDATE_RESPONSE = {
  metadata: {
    id: string;
  };
};
type CONTRACT_LINE_ITEM_UPDATE_PAYLOAD = {
  input: ServiceLineItemUpdateInput;
};

const CONTRACT_LINE_ITEM_UPDATE_MUTATION = gql`
  mutation contractLineItemUpdate($input: ServiceLineItemUpdateInput!) {
    contractLineItem_Update(input: $input) {
      metadata {
        id
      }
    }
  }
`;

export { ContractLineItemService };
