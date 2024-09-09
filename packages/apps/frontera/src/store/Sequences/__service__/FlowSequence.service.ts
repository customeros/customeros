import type { Transport } from '@store/transport';

import {
  CreateSequenceMutation,
  CreateSequenceMutationVariables,
} from '@store/Sequences/__service__/createSequence.generated';
import {
  UpdateSequenceMutation,
  UpdateSequenceMutationVariables,
} from '@store/Sequences/__service__/updateSequence.generated';
import {
  FlowSequenceLinkContactMutation,
  FlowSequenceLinkContactMutationVariables,
} from '@store/Sequences/__service__/flowSequenceLinkContact.generated';
import {
  ChangeFlowSequenceStatusMutation,
  ChangeFlowSequenceStatusMutationVariables,
} from '@store/Sequences/__service__/changeFlowSequenceStatus.generated';
import {
  FlowSequenceUnlinkContactMutation,
  FlowSequenceUnlinkContactMutationVariables,
} from '@store/Sequences/__service__/flowSequenceUnlinkContact.generated';

import { FlowSequence } from '@graphql/types';

import UpdateSequencesDocument from './updateSequence.graphql';
import GetFlowSequencesDocument from './getFlowSequences.graphql';
import LinkContactDocument from './flowSequenceLinkContact.graphql';
import CreateSequenceMutationDocument from './createSequence.graphql';
import UnlinkContactDocument from './flowSequenceUnlinkContact.graphql';
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
    return this.transport.graphql.request<{ sequences: FlowSequence[] }>(
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

  async updateSequence(payload: UpdateSequenceMutationVariables) {
    return this.transport.graphql.request<
      UpdateSequenceMutation,
      UpdateSequenceMutationVariables
    >(UpdateSequencesDocument, payload);
  }

  async linkContact(payload: FlowSequenceLinkContactMutationVariables) {
    return this.transport.graphql.request<
      FlowSequenceLinkContactMutation,
      FlowSequenceLinkContactMutationVariables
    >(LinkContactDocument, payload);
  }

  async unlinkContact(payload: FlowSequenceUnlinkContactMutationVariables) {
    return this.transport.graphql.request<
      FlowSequenceUnlinkContactMutation,
      FlowSequenceUnlinkContactMutationVariables
    >(UnlinkContactDocument, payload);
  }
}

export { FlowSequenceService };
