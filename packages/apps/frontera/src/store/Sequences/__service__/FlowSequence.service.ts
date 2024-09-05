import type { Transport } from '@store/transport.ts';

import { GetFlowSequencesQuery } from '@store/Sequences/__service__/getFlowSequences.generated.ts';
import {
  CreateSequenceMutation,
  CreateSequenceMutationVariables,
} from '@store/Sequences/__service__/createSequence.generated.ts';
import {
  ChangeFlowSequenceStatusMutation,
  ChangeFlowSequenceStatusMutationVariables,
} from '@store/Sequences/__service__/changeFlowSequenceStatus.generated.ts';

import GetFlowSequencesDocument from './getFlowSequences.graphql';
import CreateSequenceMutationDocument from './createSequence.graphql';
import ChangeFlowSequenceStatusDocument from './changeFlowSequenceStatus.graphql';

class FlowSequenceService {
  private static instance: FlowSequenceService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): FlowSequenceService {
    if (!FlowSequenceService.instance) {
      FlowSequenceService.instance = new FlowSequenceService(transport);
    }

    return FlowSequenceService.instance;
  }

  async getSequences() {
    return this.transport.graphql.request<GetFlowSequencesQuery>(
      GetFlowSequencesDocument,
    );
  }

  async createSequence(payload: CreateSequenceMutationVariables) {
    return this.transport.graphql.request<
      CreateSequenceMutation,
      CreateSequenceMutationVariables
    >(CreateSequenceMutationDocument, payload);
  }

  async updateSequenceStatus(
    payload: ChangeFlowSequenceStatusMutationVariables,
  ) {
    return this.transport.graphql.request<
      ChangeFlowSequenceStatusMutation,
      ChangeFlowSequenceStatusMutationVariables
    >(ChangeFlowSequenceStatusDocument, payload);
  }
}

export { FlowSequenceService };
