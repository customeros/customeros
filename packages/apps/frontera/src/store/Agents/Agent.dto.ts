import { Entity } from '@store/record';
import { computed, observable } from 'mobx';
import { type AgentDatum } from '@infra/repositories/agent';

import { AgentType } from '@graphql/types';

import { type AgentStore } from './Agent.store';

export class Agent extends Entity<AgentDatum> {
  @observable accessor value: AgentDatum = Agent.default();

  constructor(store: AgentStore, data: AgentDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store, data as any);
  }

  @computed
  get id() {
    return this.value.id;
  }

  static default(): AgentDatum {
    return {
      id: crypto.randomUUID(),
      name: 'Unknown',
      tenant: '',
      capabilities: [],
      goal: '',
      type: AgentType.WebVisitIdentifier,
      icon: '',
      color: '',
      visible: true,
      isActive: true,
      flowId: '',
      error: null,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
  }
}
