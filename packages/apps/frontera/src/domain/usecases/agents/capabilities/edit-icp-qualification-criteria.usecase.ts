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
  @observable accessor validationError: string = '';

  constructor() {
    this.execute = this.execute.bind(this);
    this.validate = this.validate.bind(this);
    this.setInputValue = this.setInputValue.bind(this);
  }

  @computed
  get isInvalid() {
    return this.validationError.length > 0;
  }

  @action
  setAgentId(agentId: string) {
    this.agentId = agentId;
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
  validate() {
    const span = Tracer.span('EditIcpQualificationCriteriaUsecase.validate');

    if (this.inputValue.length === 0) {
      this.validationError = 'Houston we have a blank...';
      span.end();

      return false;
    } else {
      this.validationError = '';
    }

    span.end();

    return true;
  }

  @action
  init() {
    const span = Tracer.span('EditIcpQualificationCriteriaUsecase.init', {
      inputValue: this.inputValue,
    });
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

    const config = Agent.parseCapabilityConfig(capability.config);

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
    span.end({
      inputValue: config.qualificationCriteria.value,
    });
  }

  execute() {
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

    const isValid = this.validate();

    if (!isValid) {
      return;
    }

    if (isValid) {
      this.validationError = '';
      agent?.setCapabilityConfig(
        CapabilityType.IcpQualify,
        'qualificationCriteria',
        this.inputValue,
      );

      this.service.saveAgent(agent);
    }

    span.end();
  }
}
