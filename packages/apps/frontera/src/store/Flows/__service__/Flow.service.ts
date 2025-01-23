import type { Transport } from '@infra/transport';

import {
  FlowChangeNameMutation,
  FlowChangeNameMutationVariables,
} from '@store/Flows/__service__/flowChangeName.generated.ts';
import {
  FlowArchiveBulkMutation,
  FlowArchiveBulkMutationVariables,
} from '@store/Flows/__service__/flowArchiveBulk.generated.ts';

import { Flow } from '@graphql/types';

import FlowOnDocument from './flowOn.graphql';
import GetFlowDocument from './getFlow.graphql';
import FlowOffDocument from './flowOff.graphql';
import GetFlowsDocument from './getFlows.graphql';
import MergeFlowDocument from './flowMerge.graphql';
import ArchiveFlowDocument from './flowArchive.graphql';
import FlowChangeNameDocument from './flowChangeName.graphql';
import ArchiveFlowBulkDocument from './flowArchiveBulk.graphql';
import SendTestEmailDocument from './flowEmailActionTest.graphql';
import { FlowOnMutation, FlowOnMutationVariables } from './flowOn.generated.ts';
import {
  FlowOffMutation,
  FlowOffMutationVariables,
} from './flowOff.generated.ts';
import {
  FlowMergeMutation,
  FlowMergeMutationVariables,
} from './flowMerge.generated';
import {
  FlowArchiveMutation,
  FlowArchiveMutationVariables,
} from './flowArchive.generated.ts';
import {
  FlowEmailActionTestMutation,
  FlowEmailActionTestMutationVariables,
} from './flowEmailActionTest.generated.ts';
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

  async changeFlowName(payload: FlowChangeNameMutationVariables) {
    return this.transport.graphql.request<
      FlowChangeNameMutation,
      FlowChangeNameMutationVariables
    >(FlowChangeNameDocument, payload);
  }

  async archiveFlow(payload: FlowArchiveMutationVariables) {
    return this.transport.graphql.request<
      FlowArchiveMutation,
      FlowArchiveMutationVariables
    >(ArchiveFlowDocument, payload);
  }

  async archiveFlowBulk(payload: FlowArchiveBulkMutationVariables) {
    return this.transport.graphql.request<
      FlowArchiveBulkMutation,
      FlowArchiveBulkMutationVariables
    >(ArchiveFlowBulkDocument, payload);
  }

  async startFlow(payload: FlowOnMutationVariables) {
    return this.transport.graphql.request<
      FlowOnMutation,
      FlowOnMutationVariables
    >(FlowOnDocument, payload);
  }

  async stopFlow(payload: FlowOffMutationVariables) {
    return this.transport.graphql.request<
      FlowOffMutation,
      FlowOffMutationVariables
    >(FlowOffDocument, payload);
  }

  async sendTestEmail(payload: FlowEmailActionTestMutationVariables) {
    return this.transport.graphql.request<
      FlowEmailActionTestMutation,
      FlowEmailActionTestMutationVariables
    >(SendTestEmailDocument, payload);
  }
}

export { FlowService };
