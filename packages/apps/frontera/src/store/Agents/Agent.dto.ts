import set from 'lodash/set';
import merge from 'lodash/merge';
import { Entity } from '@store/record';
import { action, computed, observable } from 'mobx';
import { type AgentDatum } from '@infra/repositories/agent';

import { unwrap, UnwrapResult } from '@utils/unwrap';
import { AgentType, CapabilityType } from '@graphql/types';

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
    super(store, data as any);
  }

  @computed
  get id() {
    return this.value.id;
  }

  @action
  public setCapabilityConfig(
    capabilityType: CapabilityType,
    property: string,
    value: unknown,
  ) {
    const foundIndex = this.value.capabilities.findIndex(
      (c) => c.type === capabilityType,
    );

    if (foundIndex === -1) {
      console.error('Agent.setCapability: Capability not found. will not set');
    }

    const config = Agent.parseCapabilityConfig(
      this.value.capabilities[foundIndex].config,
    );

    if (!config) {
      console.error('Agent.setCapability: Could not parse config');

      return;
    }

    set(config, property, value);

    this.draft();
    this.value.capabilities[foundIndex].config = JSON.stringify(config);
    this.commit({ syncOnly: true });
  }

  static parseCapabilityConfig(raw: string): CapabilityConfig | null {
    let res: UnwrapResult<CapabilityConfig> = [null, null];
    const [data, err] = res;

    (async () => (res = await unwrap(JSON.parse(raw))))();

    if (err) {
      console.error(
        'Agent.parseCapabilityConfig: Error parsing raw config',
        err,
      );

      return null;
    }

    return data;
  }

  static default(payload?: Partial<AgentDatum>): AgentDatum {
    return merge(
      {
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
      },
      payload,
    );
  }
}
