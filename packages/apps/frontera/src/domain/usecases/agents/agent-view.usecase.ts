import { Tracer } from '@infra/tracer';
import { action, observable } from 'mobx';

import { Capability, CapabilityType } from '@graphql/types';

export class AgentViewUsecase {
  @observable accessor activeCapability: Capability = {
    id: '',
    name: '',
    type: CapabilityType.IdentifyWebVisitor,
    action: '',
    active: false,
    config: '',
    errors: null,
  };

  constructor() {}

  @action
  setActiveCapability(capability: Capability) {
    Tracer.span('AgentViewUsecase.setActiveCapability', { capability });
    this.activeCapability = capability;
  }
}
