import type { Transport } from '@store/transport';

import {
  FlowMergeMutation,
  FlowMergeMutationVariables,
} from '@store/Flows/__service__/flowMerge.generated';
import {
  FlowChangeStatusMutation,
  FlowChangeStatusMutationVariables,
} from '@store/Flows/__service__/changeFlowStatus.generated';

import { Flow } from '@graphql/types';

import GetFlowDocument from './getFlow.graphql';
import GetFlowsDocument from './getFlows.graphql';
import MergeFlowDocument from './flowMerge.graphql';
import ChangeStatusDocument from './changeFlowStatus.graphql';

class FlowService {
  private static instance: FlowService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): FlowService {
    if (!FlowService.instance) {
      FlowService.instance = new FlowService(transport);
    }

    return FlowService.instance;
  }

  async getFlows() {
    return this.transport.graphql.request<{ flows: Flow[] }>(GetFlowsDocument);
  }

  async getFlow(id: string) {
    return this.transport.graphql.request<{ flow: Flow }, { id: string }>(
      GetFlowDocument,
      {
        id,
      },
    );
  }

  async mergeFlow(payload: FlowMergeMutationVariables) {
    return this.transport.graphql.request<
      FlowMergeMutation,
      FlowMergeMutationVariables
    >(MergeFlowDocument, payload);
  }

  async changeStatus(payload: FlowChangeStatusMutationVariables) {
    return this.transport.graphql.request<
      FlowChangeStatusMutation,
      FlowChangeStatusMutationVariables
    >(ChangeStatusDocument, payload);
  }
}

export { FlowService };
