import set from 'lodash/set';
import omit from 'lodash/omit';
import merge from 'lodash/merge';
import { Entity } from '@store/record';
import { Tracer } from '@infra/tracer';
import { action, computed, observable } from 'mobx';
import { type AgentDatum } from '@infra/repositories/agent';

import { IconName } from '@ui/media/Icon';
import { AgentType, CapabilityType, AgentListenerEvent } from '@graphql/types';

import { AgentStore } from './Agent.store';

type CapabilityConfig = {
  [key: string]: {
    value: unknown;
    error: string | null;
  };
};

type ColorType =
  | 'grayModern'
  | 'error'
  | 'warning'
  | 'success'
  | 'grayWarm'
  | 'moss'
  | 'blueLight'
  | 'indigo'
  | 'violet'
  | 'pink';

export class Agent extends Entity<AgentDatum> {
  @observable accessor value: AgentDatum = Agent.default();
  static defaultNameByType: Record<AgentType, string> = {
    [AgentType.WebVisitIdentifier]: 'Web visitor identifier',
    [AgentType.SupportSpotter]: 'Support spotter',
    [AgentType.IcpQualifier]: 'Icp qualifier',
    [AgentType.CampaignManager]: 'Outbound campaign manager',
    [AgentType.CashflowGuardian]: 'Cashflow guardian',
    [AgentType.MeetingKeeper]: 'Meeting keeper',
  };

  constructor(store: AgentStore, data: AgentDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  @computed
  get id() {
    return this.value.id;
  }

  @computed
  get defaultName() {
    return Agent.defaultNameByType[this.value.type] || 'Unknown';
  }

  @action
  public setName(name: string) {
    this.draft();
    this.value.name = name;
    this.commit({ syncOnly: true });
  }

  @action
  public setIcon(iconName: IconName) {
    this.draft();
    this.value.icon = iconName;
    this.commit({ syncOnly: true });
  }

  @action
  public setColor(color: string) {
    this.draft();
    this.value.color = color;
    this.commit({ syncOnly: true });
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

      return;
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

  public toDuplicatePayload(
    name: string,
  ): Omit<
    AgentDatum,
    'createdAt' | 'updatedAt' | 'isConfigured' | 'goalType' | 'id'
  > {
    return omit({ ...this.value, name }, [
      'createdAt',
      'updatedAt',
      'error',
      'isConfigured',
      'goalType',
      'id',
    ]);
  }

  public put(payload: AgentDatum) {
    const span = Tracer.span('Agent.put');

    this.draft();
    this.value = merge(this.value, payload);
    this.commit({ syncOnly: true });

    span.end();
  }

  get colorMap(): [ring: string, bg: string, iconColor: string] {
    const options: Record<
      ColorType,
      [ring: string, bg: string, iconColor: string]
    > = {
      grayModern: [
        'group-hover:ring-grayModern-400',
        'group-hover:bg-grayModern-50',
        'group-hover:text-grayModern-500',
      ],
      error: [
        'group-hover:ring-error-400',
        'group-hover:bg-error-50',
        'group-hover:text-error-500',
      ],
      warning: [
        'group-hover:ring-warning-400',
        'group-hover:bg-warning-50',
        'group-hover:text-warning-500',
      ],
      success: [
        'group-hover:ring-success-400',
        'group-hover:bg-success-50',
        'group-hover:text-success-500',
      ],
      grayWarm: [
        'group-hover:ring-grayWarm-400',
        'group-hover:bg-grayWarm-50',
        'group-hover:text-grayWarm-500',
      ],
      moss: [
        'group-hover:ring-moss-400',
        'group-hover:bg-moss-50',
        'group-hover:text-moss-500',
      ],
      blueLight: [
        'group-hover:ring-blueLight-400',
        'group-hover:bg-blueLight-50',
        'group-hover:text-blueLight-500',
      ],
      indigo: [
        'group-hover:ring-indigo-400',
        'group-hover:bg-indigo-50',
        'group-hover:text-indigo-500',
      ],
      violet: [
        'group-hover:ring-violet-400',
        'group-hover:bg-violet-50',
        'group-hover:text-violet-500',
      ],
      pink: [
        'group-hover:ring-pink-400',
        'group-hover:bg-pink-50',
        'group-hover:text-pink-500',
      ],
    };

    return options?.[this.value.color as ColorType] || options.grayModern;
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

  static default(payload?: Partial<AgentDatum>): AgentDatum {
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
