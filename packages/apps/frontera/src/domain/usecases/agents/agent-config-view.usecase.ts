import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';

import {
  AgentType,
  Capability,
  AgentListener,
  CapabilityType,
  AgentListenerEvent,
} from '@graphql/types';

export class AgentConfigViewUsecase {
  private root = RootStore.getInstance();
  @observable private accessor _activeConfigId: string = '';

  constructor(private id: string, private defaultConfigId?: string | null) {
    if (this.defaultConfigId?.length) {
      this._activeConfigId = this.defaultConfigId;
    }

    this.setActiveConfig = this.setActiveConfig.bind(this);
  }

  @computed
  get agent() {
    return this.root.agents.getById(this.id);
  }

  @computed
  get activeConfig() {
    const agent = this.agent;

    if (!agent) return null;

    return (
      agent?.value.listeners.find((c) =>
        this._activeConfigId
          ? c.id === this._activeConfigId
          : c.type === AgentConfigViewUsecase.defaultConfigMap[agent.type],
      ) ??
      agent?.value.capabilities.find((c) =>
        this._activeConfigId
          ? c.id === this._activeConfigId
          : c.type === AgentConfigViewUsecase.defaultConfigMap[agent.type],
      )
    );
  }

  @action
  setActiveConfig(listenerOrCapability: Capability | AgentListener) {
    const span = Tracer.span('AgentViewUsecase.setActiveConfig', {
      previous: this.activeConfig?.type,
    });

    this._activeConfigId = listenerOrCapability.id;

    span.end({ current: this.activeConfig?.type });
  }

  private static defaultConfigMap: Record<
    AgentType,
    CapabilityType | AgentListenerEvent
  > = {
    [AgentType.WebVisitIdentifier]: AgentListenerEvent.NewWebSession,
    [AgentType.IcpQualifier]: CapabilityType.IcpQualify,
    [AgentType.SupportSpotter]: CapabilityType.DetectSupportWebvisit,
    [AgentType.MeetingKeeper]: CapabilityType.AddMeetingNotesToCompany,
    [AgentType.CampaignManager]: CapabilityType.ManageCampaignExecution,
  };
}
