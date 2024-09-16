import type { Transport } from '@store/transport';

import {
  FlowContactDeleteMutation,
  FlowContactDeleteMutationVariables,
} from '@store/FlowContacts/__service__/flowContactDelete.generated';

import DeleteContactDocument from './flowContactDelete.graphql';

class FlowContactsService {
  private static instance: FlowContactsService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): FlowContactsService {
    if (!FlowContactsService.instance) {
      FlowContactsService.instance = new FlowContactsService(transport);
    }

    return FlowContactsService.instance;
  }

  async deleteFlowContact(payload: FlowContactDeleteMutationVariables) {
    return this.transport.graphql.request<
      FlowContactDeleteMutation,
      FlowContactDeleteMutationVariables
    >(DeleteContactDocument, payload);
  }
}

export { FlowContactsService };
