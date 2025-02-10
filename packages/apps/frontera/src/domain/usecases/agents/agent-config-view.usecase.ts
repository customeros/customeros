import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';

import { Capability } from '@graphql/types';

export class AgentConfigViewUsecase {
  private root = RootStore.getInstance();
  @observable private accessor _activeCapabilityId: string = '';

  constructor(private id: string, private defaultCapabilityId?: string | null) {
    if (this.defaultCapabilityId?.length) {
      this._activeCapabilityId = this.defaultCapabilityId;
    }
    this.setActiveCapability = this.setActiveCapability.bind(this);
  }

  @computed
  get agent() {
    return this.root.agents.getById(this.id);
  }

  @computed
  get activeCapability() {
    if (!this.agent) return null;

    if (!this._activeCapabilityId) {
      return this.agent?.value.capabilities?.find((e) => !!e.config.length);
    }

    return this.agent?.value.capabilities.find(
      (c) => c.id === this._activeCapabilityId,
    );
  }

  @action
  setActiveCapability(capability: Capability) {
    const span = Tracer.span('AgentViewUsecase.setActiveCapability', {
      previousActiveCapability: this.activeCapability?.type,
    });

    this._activeCapabilityId = capability.id;

    span.end({ currentActiveCapability: this.activeCapability?.type });
  }
}
