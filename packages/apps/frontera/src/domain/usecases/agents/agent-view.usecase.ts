import { action, observable } from 'mobx';

import { Capability, CapabilityType } from '@graphql/types';

export class AgentViewUsecase {
  @observable accessor activeCapability: Capability = {
    id: '',
    name: '',
    type: CapabilityType.WebsiteTracker,
    action: '',
    optional: false,
    values: '',
    errors: null,
  };

  constructor() {}

  @action
  setActiveCapability(capability: Capability) {
    this.activeCapability = capability;
  }
}
