import { runInAction } from 'mobx';
import { Store } from '@store/_store';
import { type RootStore } from '@store/root';
import { type Transport } from '@infra/transport';
import { AgentRepository } from '@infra/repositories/agent';
import { type AgentDatum } from '@infra/repositories/agent';

import { unwrap } from '@utils/unwrap';

import { Agent } from './Agent.dto';

export class AgentStore extends Store<AgentDatum, Agent> {
  private agentRepository = new AgentRepository();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Agents',
      getId: (data) => data?.id,
      factory: Agent,
    });
  }

  public async bootstrap() {
    const [data, err] = await unwrap(this.agentRepository.getAgents());

    if (err) {
      console.error('Error bootstrapping agents:', err);

      return;
    }
    runInAction(() => {
      data?.agents.forEach((datum) => {
        this.value.set(datum.id, new Agent(this, datum));
      });
      this.isBootstrapped = true;
      this.isBootstrapping = false;
    });
  }
}
