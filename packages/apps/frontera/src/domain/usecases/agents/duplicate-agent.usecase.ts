import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto.ts';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

export class DuplicateAgentUsecase {
  private root = RootStore.getInstance();
  private service = new AgentService();
  @observable accessor inputValue: string = '';
  @observable accessor isSaving: boolean = false;

  constructor(private id?: string) {
    this.execute = this.execute.bind(this);
    this.setInputValue = this.setInputValue.bind(this);
  }

  @action
  setInputValue(value: string) {
    this.inputValue = value;
  }

  @computed
  get agent(): undefined | null | Agent {
    return this.id ? this.root.agents.getById(this.id) : null;
  }

  public async execute() {
    this.isSaving = true;

    const agent = this.agent;

    if (!agent) {
      console.error('DuplicateAgentUsecase.execute: Could not find agent');
      this.isSaving = false;

      return;
    }

    const span = Tracer.span('DuplicateAgentUsecase.execute', {
      payload: {
        agent,
        name: this.inputValue,
      },
    });

    const [res, err] = await this.service.duplicateAgent(
      agent,
      this.inputValue,
    );

    if (err) {
      this.isSaving = false;
      console.error(
        'DuplicateAgentUsecase.execute: Could not duplicate agent',
        err,
      );
      span.end();

      return;
    }

    if (res?.agent_Save) {
      window.location.href = `/agents/${res.agent_Save.id}`;

      this.isSaving = false;
      this.root.agents.addOne(res.agent_Save);
      this.root.ui.commandMenu.setType('AgentCommands');
      this.root.ui.commandMenu.toggle();
    }

    span.end();
  }
}
