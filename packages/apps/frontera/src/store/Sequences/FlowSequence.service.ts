import type { Transport } from '@store/transport';

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
}

export { FlowSequenceService };
