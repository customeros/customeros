import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { CapabilityType } from '@graphql/types';

export class EditIcpDisqualificationCriteriaUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();
  private agentId: string = '';
  @observable accessor inputValue: string = '';
  @observable accessor validationError: string = '';
  @observable accessor disqualificationCriteriaError: string = '';

  constructor() {
    this.execute = this.execute.bind(this);
    this.validate = this.validate.bind(this);
    this.setInputValue = this.setInputValue.bind(this);
    this.init = this.init.bind(this);
  }

  @computed
  get isInvalid() {
    return this.validationError.length > 0;
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
    const span = Tracer.span('EditIcpDisqualificationCriteriaUsecase.validate');

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
  init(agentId: string) {
    const span = Tracer.span('EditIcpDisqualificationCriteriaUsecase.init', {
      inputValue: this.inputValue,
    });

    this.agentId = agentId;

    const agent = this.root.agents.getById(agentId);

    if (!agent) {
      console.error(
        'EditIcpDisqualificationCriteriaUsecase.init: Agent not found. aborting execution',
      );

      return;
    }

    const capability = agent.value.capabilities.find(
      (c) => c.type === CapabilityType.IcpQualify,
    );

    if (!capability) {
      console.error(
        'EditIcpDisqualificationCriteriaUsecase.init: Capability not found. aborting execution',
      );

      return;
    }

    const config = Agent.parseConfig(capability.config);

    if (!config) {
      console.error(
        'EditIcpDisqualificationCriteriaUsecase.init: Could not parse config. aborting',
      );

      return;
    }

    if (!config.disqualificationCriteria) {
      console.error(
        'EditIcpDisqualificationCriteriaUsecase.init: disqualificationCriteria not found in config. aborting execution',
      );

      return;
    }

    this.inputValue = config.disqualificationCriteria.value as string;
    this.disqualificationCriteriaError = config.disqualificationCriteria
      .error as string;
    span.end({
      inputValue: config.disqualificationCriteria.value,
      disqualificationCriteriaError: this.disqualificationCriteriaError,
    });
  }

  async execute() {
    const span = Tracer.span('EditIcpDisqualificationCriteriaUsecase.execute', {
      inputValue: this.inputValue,
    });
    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'EditIcpDisqualificationCriteriaUsecase: Agent not found. aborting execution',
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
        'disqualificationCriteria',
        this.inputValue,
      );

      const [res, err] = await this.service.saveAgent(agent);

      if (err) {
        console.error(
          'EditIcpDisqualificationCriteriaUsecase.execute: Error saving agent. aborting execution',
        );
      }

      if (res) {
        agent.put(res?.agent_Save);
        this.init(this.agentId);
      }
    }

    span.end();
  }
}
