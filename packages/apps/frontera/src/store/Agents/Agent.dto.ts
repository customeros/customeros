import set from 'lodash/set';
import omit from 'lodash/omit';
import merge from 'lodash/merge';
import { Entity } from '@store/record';
import { Tracer } from '@infra/tracer';
import { action, computed, observable } from 'mobx';
import { type AgentDatum } from '@infra/repositories/agent';

import { AgentType, CapabilityType, AgentListenerEvent } from '@graphql/types';

import { AgentStore } from './Agent.store';

type CapabilityConfig = {
  [key: string]: {
    value: unknown;
    error: string | null;
  };
};

export class Agent extends Entity<AgentDatum> {
  @observable accessor value: AgentDatum = Agent.default();

  constructor(store: AgentStore, data: AgentDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  @computed
  get id() {
    return this.value.id;
  }

  @action
  public setListenerConfig(
    listenerType: AgentListenerEvent,
    property: string,
    value: unknown,
  ) {
    const span = Tracer.span('Agent.setListenerConfig');
    const foundIndex = this.value.listeners.findIndex(
      (c) => c.type === listenerType,
    );

    if (foundIndex === -1) {
      console.error(
        'Agent.setListenerConfig: Listener not found. will not set',
      );
    }

    const config = Agent.parseConfig(this.value.listeners[foundIndex].config);

    if (!config) {
      console.error('Agent.setListenerConfig: Could not parse config');

      return;
    }

    set(config, `${property}.value`, value);

    this.draft();
    this.value.listeners[foundIndex].config = JSON.stringify(config);
    this.commit({ syncOnly: true });

    span.end();
  }

  @action
  public setCapabilityConfig(
    capabilityType: CapabilityType,
    property: string,
    value: unknown,
  ) {
    const span = Tracer.span('Agent.setCapabilityConfig');
    const foundIndex = this.value.capabilities.findIndex(
      (c) => c.type === capabilityType,
    );

    if (foundIndex === -1) {
      console.error(
        'Agent.setCapabilityConfig: Capability not found. will not set',
      );
    }

    const config = Agent.parseConfig(
      this.value.capabilities[foundIndex].config,
    );

    if (!config) {
      console.error('Agent.setCapabilityConfig: Could not parse config');

      return;
    }

    set(config, `${property}.value`, value);

    this.draft();
    this.value.capabilities[foundIndex].config = JSON.stringify(config);
    this.commit({ syncOnly: true });

    span.end();
  }

  public toggleStatus() {
    this.draft();
    this.value.isActive = !this.value.isActive;
    this.commit({ syncOnly: true });
  }

  public toPayload(): Omit<
    AgentDatum,
    'createdAt' | 'updatedAt' | 'isConfigured' | 'goalType'
  > {
    return omit(this.value, [
      'createdAt',
      'updatedAt',
      'error',
      'isConfigured',
      'goalType',
    ]);
  }

  public put(payload: AgentDatum) {
    const span = Tracer.span('Agent.put');

    this.draft();
    this.value = merge(this.value, payload);
    this.commit({ syncOnly: true });

    span.end();
  }

  static parseConfig(raw: string): CapabilityConfig | null {
    const span = Tracer.span('Agent.parseConfig', { raw });

    if (raw === '') {
      span.end();

      return {};
    }

    const parsed = JSON.parse(raw);

    span.end();

    return parsed;
  }

  default(payload?: Partial<AgentDatum>): AgentDatum {
    return merge(
      {
        id: crypto.randomUUID(),
        name: 'Unknown',
        tenant: '',
        listeners: [],
        capabilities: [],
        goal: '',
        goalType: '',
        type: AgentType.WebVisitIdentifier,
        icon: '',
        color: '',
        visible: true,
        isActive: true,
        flowId: '',
        error: null,
        isConfigured: false,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      payload,
    );
  }
}
