import type { Transport } from '@infra/transport';

import {
  FlowParticipantAddMutation,
  FlowParticipantAddMutationVariables,
} from '@store/FlowParticipants/__service__/flowParticipantAdd.generated.ts';
import {
  FlowParticipantDeleteMutation,
  FlowParticipantDeleteMutationVariables,
} from '@store/FlowParticipants/__service__/flowParticipantDelete.generated.ts';
import {
  FlowParticipantAddBulkMutation,
  FlowParticipantAddBulkMutationVariables,
} from '@store/FlowParticipants/__service__/flowParticipantAddBulk.generated.ts';
import {
  FlowParticipantDeleteBulkMutation,
  FlowParticipantDeleteBulkMutationVariables,
} from '@store/FlowParticipants/__service__/flowParticipantsDeleteBulk.generated.ts';

import GetFlowParticipantDocument from './getFlowParticipant.graphql';
import CreateFlowParticipantDocument from './flowParticipantAdd.graphql';
import DeleteFlowParticipantDocument from './flowParticipantDelete.graphql';
import CreateFlowParticipantsDocument from './flowParticipantAddBulk.graphql';
import DeleteFlowParticipantsDocument from './flowParticipantsDeleteBulk.graphql';
import {
  GetFlowParticipantQuery,
  GetFlowParticipantQueryVariables,
} from './getFlowParticipant.generated';

class FlowParticipantsService {
  private static instance: FlowParticipantsService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): FlowParticipantsService {
    if (!FlowParticipantsService.instance) {
      FlowParticipantsService.instance = new FlowParticipantsService(transport);
    }

    return FlowParticipantsService.instance;
  }

  async getFlowParticipant(payload: GetFlowParticipantQueryVariables) {
    return this.transport.graphql.request<
      GetFlowParticipantQuery,
      GetFlowParticipantQueryVariables
    >(GetFlowParticipantDocument, payload);
  }

  async deleteFlowParticipant(payload: FlowParticipantDeleteMutationVariables) {
    return this.transport.graphql.request<
      FlowParticipantDeleteMutation,
      FlowParticipantDeleteMutationVariables
    >(DeleteFlowParticipantDocument, payload);
  }

  async deleteFlowParticipants(
    payload: FlowParticipantDeleteBulkMutationVariables,
  ) {
    return this.transport.graphql.request<
      FlowParticipantDeleteBulkMutation,
      FlowParticipantDeleteBulkMutationVariables
    >(DeleteFlowParticipantsDocument, payload);
  }

  async createFlowParticipant(payload: FlowParticipantAddMutationVariables) {
    return this.transport.graphql.request<
      FlowParticipantAddMutation,
      FlowParticipantAddMutationVariables
    >(CreateFlowParticipantDocument, payload);
  }

  async createFlowParticipants(
    payload: FlowParticipantAddBulkMutationVariables,
  ) {
    return this.transport.graphql.request<
      FlowParticipantAddBulkMutation,
      FlowParticipantAddBulkMutationVariables
    >(CreateFlowParticipantsDocument, payload);
  }
}

export { FlowParticipantsService };
