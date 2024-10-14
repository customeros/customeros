import type { Transport } from '@store/transport';

import {
  FlowSenderMergeMutation,
  FlowSenderMergeMutationVariables,
} from '@store/FlowSenders/__service__/flowSenderMerge.generated.ts';
import {
  FlowSenderDeleteMutation,
  FlowSenderDeleteMutationVariables,
} from '@store/FlowSenders/__service__/flowSenderDelete.generated';

import CreateSenderDocument from './flowSenderMerge.graphql';
import DeleteSenderDocument from './flowSenderDelete.graphql';

class FlowSendersService {
  private static instance: FlowSendersService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): FlowSendersService {
    if (!FlowSendersService.instance) {
      FlowSendersService.instance = new FlowSendersService(transport);
    }

    return FlowSendersService.instance;
  }

  async deleteFlowSender(payload: FlowSenderDeleteMutationVariables) {
    return this.transport.graphql.request<
      FlowSenderDeleteMutation,
      FlowSenderDeleteMutationVariables
    >(DeleteSenderDocument, payload);
  }

  async createFlowSender(payload: FlowSenderMergeMutationVariables) {
    return this.transport.graphql.request<
      FlowSenderMergeMutation,
      FlowSenderMergeMutationVariables
    >(CreateSenderDocument, payload);
  }
}

export { FlowSendersService };
