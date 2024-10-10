import type { Transport } from '@store/transport';

import { EmailVariableEntity } from '@graphql/types';

import GetFlowEmailVariablesDocument from './getFlowEmailVariables.graphql';

class FlowEmailVariablesService {
  private static instance: FlowEmailVariablesService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): FlowEmailVariablesService {
    if (!FlowEmailVariablesService.instance) {
      FlowEmailVariablesService.instance = new FlowEmailVariablesService(
        transport,
      );
    }

    return FlowEmailVariablesService.instance;
  }

  async getFlowEmailVariables() {
    return this.transport.graphql.request<{
      flow_emailVariables: EmailVariableEntity[];
    }>(GetFlowEmailVariablesDocument);
  }
}

export { FlowEmailVariablesService };
