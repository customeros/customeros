import type { Transport } from '@store/transport';

import {
  FlowMergeMutation,
  FlowMergeMutationVariables,
} from '@store/Flows/__service__/flowMerge.generated';
import {
  FlowContactAddMutation,
  FlowContactAddMutationVariables,
} from '@store/Flows/__service__/flowContactAdd.generated';
import {
  FlowChangeStatusMutation,
  FlowChangeStatusMutationVariables,
} from '@store/Flows/__service__/changeFlowStatus.generated';
import {
  FlowContactAddBulkMutation,
  FlowContactAddBulkMutationVariables,
} from '@store/Flows/__service__/flowContactAddBulk.generated.ts';

import { Flow } from '@graphql/types';

import GetFlowDocument from './getFlow.graphql';
import GetFlowsDocument from './getFlows.graphql';
import MergeFlowDocument from './flowMerge.graphql';
import AddContactDocument from './flowContactAdd.graphql';
import ChangeStatusDocument from './changeFlowStatus.graphql';
import AddContactBulkDocument from './flowContactAddBulk.graphql';

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

  async addContact(payload: FlowContactAddMutationVariables) {
    return this.transport.graphql.request<
      FlowContactAddMutation,
      FlowContactAddMutationVariables
    >(AddContactDocument, payload);
  }

  async addContactBulk(payload: FlowContactAddBulkMutationVariables) {
    return this.transport.graphql.request<
      FlowContactAddBulkMutation,
      FlowContactAddBulkMutationVariables
    >(AddContactBulkDocument, payload);
  }

  async changeStatus(payload: FlowChangeStatusMutationVariables) {
    return this.transport.graphql.request<
      FlowChangeStatusMutation,
      FlowChangeStatusMutationVariables
    >(ChangeStatusDocument, payload);
  }
}

export { FlowService };
