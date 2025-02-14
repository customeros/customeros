import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

export class RenameAgentUsecase {
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
      console.error('RenameAgentUsecase.execute: Could not find agent');
      this.isSaving = false;

      return;
    }

    const span = Tracer.span('RenameAgentUsecase.execute', {
      payload: {
        agent,
        name: this.inputValue,
      },
    });

    // Default to 'Name me maybe' if the input is empty
    const newAgentName = !this.inputValue.trim().length
      ? 'Name me maybe'
      : this.inputValue;

    const prevName = agent.value.name;

    agent.setName(newAgentName);

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      this.isSaving = false;
      agent.setName(prevName);
      console.error('RenameAgentUsecase.execute: Could not rename agent', err);
      span.end();

      return;
    }

    if (res?.agent_Save) {
      this.isSaving = false;
      agent.put(res?.agent_Save);
      this.root.ui.commandMenu.setType('AgentCommands');
      this.root.ui.commandMenu.setOpen(false);
    }

    span.end();
  }
}
