import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { Capability } from '@graphql/types';
const agentData = {
  __typename: 'Query',
  value: {
    __typename: 'Agent',
    id: 'campaign-manager-1',
    type: 'campaign_manager',
    name: 'Campaign manager',
    goal: 'receive_reply',
    isActive: true,
    flowId: null,
    visible: true,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    error: null,
    color: '#000000',
    icon: '',
    capabilities: [
      {
        __typename: 'Capability',
        id: 'cap-1',
        type: 'enrich_email_address',
        name: 'Enrich Email Address',
        action: 'enrich_email_address',
        active: true,
        config: '',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-2',
        type: 'validate_email_deliverability',
        name: 'Validate Email Deliverability',
        action: 'validate_email_deliverability',
        active: true,
        config: '',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-3',
        type: 'manage_campaign_execution',
        name: 'Manage Campaign Execution',
        action: 'manage_campaign_execution',
        active: true,
        config: '{}',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-4',
        type: 'select_optimal_sending_mailbox',
        name: 'Select Optimal Sending Mailbox',
        action: 'select_optimal_sending_mailbox',
        active: true,
        config: '{}',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-5',
        type: 'manage_email_delivery_failures',
        name: 'Manage Email Delivery Failures',
        action: 'manage_email_delivery_failures',
        active: true,
        config: '',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-6',
        type: 'forward_email_replies',
        name: 'Forward Email Replies',
        action: 'forward_email_replies',
        active: true,
        config: '{}',
        errors: null,
      },
    ],
  },
};
export class AgentConfigViewUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();
  @observable private accessor _activeCapabilityId: string = '';

  constructor(private id: string, private defaultCapabilityId?: string | null) {
    console.log('🏷️ ----- defaultCapabilityId: ', defaultCapabilityId);

    if (this.defaultCapabilityId?.length) {
      this._activeCapabilityId = this.defaultCapabilityId;
    }
    this.toggleActive = this.toggleActive.bind(this);
    this.setActiveCapability = this.setActiveCapability.bind(this);
  }

  @computed
  get agent() {
    return agentData;
  }

  @computed
  get activeCapability() {
    if (!this.agent) return null;

    if (!this._activeCapabilityId) {
      console.log('🏷️ ----- : A', this._activeCapabilityId);
      console.log('🏷️ ----- : A', this._activeCapabilityId);

      return this.agent?.value.capabilities.find((c) => c.config.length > 0);
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

    console.log(
      '🏷️ ----- this.activeCapability: ',
      this.activeCapability?.type,
    );
    this._activeCapabilityId = capability.id;

    span.end({ currentActiveCapability: this.activeCapability?.type });
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
}
