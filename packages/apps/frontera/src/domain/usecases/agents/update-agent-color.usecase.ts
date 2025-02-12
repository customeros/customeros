import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { computed, observable } from 'mobx';
import { Agent } from '@store/Agents/Agent.dto.ts';
import { AgentService } from '@domain/services/agent/agent.service';

export class UpdateAgentColorUsecase {
  private root = RootStore.getInstance();
  private service = new AgentService();
  @observable accessor inputValue: string = '';
  @observable accessor isSaving: boolean = false;

  constructor(private id?: string) {
    this.execute = this.execute.bind(this);
  }

  @computed
  get agent(): undefined | null | Agent {
    return this.id ? this.root.agents.getById(this.id) : null;
  }

  public async execute(color: string) {
    this.isSaving = true;

    const agent = this.agent;

    if (!agent) {
      console.error('UpdateAgentColorUsecase.execute: Could not find agent');
      this.isSaving = false;

      return;
    }

    const span = Tracer.span('UpdateAgentColorUsecase.execute', {
      payload: {
        agent,
        color,
      },
    });

    agent.setColor(color);

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      this.isSaving = false;
      console.error(
        'UpdateAgentColorUsecase.execute: Could not change agent color',
        err,
      );
      span.end();

      return;
    }

    if (res?.agent_Save) {
      this.isSaving = false;
      agent.put(res?.agent_Save);
    }

    span.end();
  }
}
