import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { CapabilityType } from '@graphql/types';

export class EditIcpQualificationCriteriaUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();
  private agentId: string = '';
  @observable accessor inputValue: string = '';
  @observable accessor qualificationCriteriaError: string = '';
  constructor() {
    this.execute = this.execute.bind(this);
    this.setInputValue = this.setInputValue.bind(this);
  }

  @action
  setInputValue(val: string) {
    this.inputValue = val;
  }

  @computed
  get capabilityErrors() {
    return this.root.agents
      .getById(this.agentId)
      ?.value.capabilities.find((c) => c.type === CapabilityType.IcpQualify)
      ?.errors;
  }

  @action
  init(agentId: string) {
    const span = Tracer.span('EditIcpQualificationCriteriaUsecase.init', {
      inputValue: this.inputValue,
    });

    this.agentId = agentId;

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'EditIcpQualificationCriteriaUsecase.init: Agent not found. aborting execution',
      );

      return;
    }

    const capability = agent.value.capabilities.find(
      (c) => c.type === CapabilityType.IcpQualify,
    );

    if (!capability) {
      console.error(
        'EditIcpQualificationCriteriaUsecase.init: Capability not found. aborting execution',
      );

      return;
    }

    const config = Agent.parseConfig(capability.config);

    if (!config) {
      console.error(
        'EditIcpQualificationCriteriaUsecase.init: Could not parse config. aborting',
      );

      return;
    }

    if (!config.qualificationCriteria) {
      console.error(
        'EditIcpQualificationCriteriaUsecase.init: QualificationCriteria not found in config. aborting execution',
      );

      return;
    }

    this.inputValue = config.qualificationCriteria.value as string;
    this.qualificationCriteriaError = config.qualificationCriteria
      .error as string;
    span.end({
      inputValue: config.qualificationCriteria.value,
      qualificationCriteriaError: this.qualificationCriteriaError,
    });
  }

  async execute() {
    const span = Tracer.span('EditIcpQualificationCriteriaUsecase.execute', {
      inputValue: this.inputValue,
    });
    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'EditIcpQualificationCriteriaUsecase: Agent not found. aborting execution',
      );

      return;
    }

    agent?.setCapabilityConfig(
      CapabilityType.IcpQualify,
      'qualificationCriteria',
      this.inputValue,
    );

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      console.error(
        'EditIcpQualificationCriteriaUsecase.execute: Error saving agent. aborting execution',
      );
    }

    if (res) {
      agent.put(res?.agent_Save);
      this.init(this.agentId);
    }

    span.end();
  }
}
