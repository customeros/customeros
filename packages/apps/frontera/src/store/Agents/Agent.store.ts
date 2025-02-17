import { Store } from '@store/_store';
import { Tracer } from '@infra/tracer';
import { action, runInAction } from 'mobx';
import { type RootStore } from '@store/root';
import { type Transport } from '@infra/transport';
import { type AgentDatum, AgentRepository } from '@infra/repositories/agent';

import { unwrap } from '@utils/unwrap';
import { AgentType } from '@graphql/types';

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

  getFirstAgentByType(type: AgentType) {
    return this.toArray().find((agent) => agent.value.type === type);
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

  @action
  public addOne(agent: AgentDatum) {
    const span = Tracer.span('AgentStore.addOne', {
      payload: agent,
    });

    this.value.set(agent.id, new Agent(this, agent));
    this.version++;

    span.end();
  }
}
