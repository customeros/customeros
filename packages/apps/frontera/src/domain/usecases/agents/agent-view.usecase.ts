import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import {
  AgentType,
  Capability,
  AgentListener,
  CapabilityType,
  AgentListenerEvent,
} from '@graphql/types';

export class AgentViewUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();
  @observable private accessor _activeConfigId: string = '';

  constructor(private id: string, private defaultConfigId?: string | null) {
    if (this.defaultConfigId?.length) {
      this._activeConfigId = this.defaultConfigId;
    }

    this.toggleActive = this.toggleActive.bind(this);
    this.setActiveConfig = this.setActiveConfig.bind(this);
  }

  @computed
  get agent() {
    return this.root.agents.getById(this.id);
  }

  @computed
  get activeConfig() {
    if (!this.agent) return null;

    return (
      this.agent?.value.listeners.find((c) =>
        this._activeConfigId
          ? c.id === this._activeConfigId
          : c.type ===
            AgentViewUsecase.defaultConfigMap[
              this.agent?.value.type ?? AgentType.WebVisitIdentifier
            ],
      ) ??
      this.agent?.value.capabilities.find((c) =>
        this._activeConfigId
          ? c.id === this._activeConfigId
          : c.type ===
            AgentViewUsecase.defaultConfigMap[
              this.agent?.value.type ?? AgentType.WebVisitIdentifier
            ],
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

  toggleActive() {
    const span = Tracer.span('AgentViewUsecase.toggleActive');

    if (!this.agent) {
      console.error('AgentViewUsecase.toggleActive: Agent not found. Aborting');

      return;
    }

    this.agent.toggleStatus();
    this.service.saveAgent(this.agent);

    span.end();
  }

  private static defaultConfigMap: {
    [key in AgentType]?: AgentListenerEvent | CapabilityType;
  } = {
    [AgentType.WebVisitIdentifier]: AgentListenerEvent.NewWebSession,
    [AgentType.IcpQualifier]: CapabilityType.IcpQualify,
    [AgentType.SupportSpotter]: AgentListenerEvent.NewMeetingRecording,
    [AgentType.MeetingKeeper]: AgentListenerEvent.NewMeetingRecording,
    [AgentType.CashflowGuardian]: CapabilityType.GenerateInvoice,
    [AgentType.CampaignManager]: CapabilityType.ManageCampaignExecution,
  };
}
