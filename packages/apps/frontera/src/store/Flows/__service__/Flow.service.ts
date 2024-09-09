import type { Transport } from '@store/transport.ts';

import { Flow } from '@graphql/types';

import GetFlowsDocument from './getFlows.graphql';

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
}

export { FlowService };
