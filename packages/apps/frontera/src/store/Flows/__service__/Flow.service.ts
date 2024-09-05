import type { Transport } from '@store/transport.ts';

import { GetFlowsQuery } from '@store/Flows/__service__/getFlows.generated.ts';

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
    return this.transport.graphql.request<GetFlowsQuery>(GetFlowsDocument);
  }
}

export { FlowService };
