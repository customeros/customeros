import { Store } from '@store/_store';
import { type RootStore } from '@store/root';
import { type Transport } from '@infra/transport';
import { type AgentDatum } from '@infra/repositories/agent';

import { Agent } from './Agent.dto';

export class AgentStore extends Store<AgentDatum, Agent> {
  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Agents',
      getId: (data) => data.id,
      factory: Agent,
    });
  }
}
