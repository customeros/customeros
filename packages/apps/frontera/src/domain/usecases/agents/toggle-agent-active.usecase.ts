import { computed } from 'mobx';
import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { AgentService } from '@domain/services/agent/agent.service';

export class ToggleAgentActiveUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();

  constructor(private id: string) {
    this.toggleActive = this.toggleActive.bind(this);
  }

  @computed
  get agent() {
    return this.root.agents.getById(this.id);
  }

  toggleActive() {
    const span = Tracer.span('ToggleAgentActiveUsecase.toggleActive');

    if (!this.agent) {
      console.error(
        'ToggleAgentActiveUsecase.toggleActive: Agent not found. Aborting',
      );

      return;
    }

    this.agent.toggleStatus();
    this.service.saveAgent(this.agent);

    span.end();
  }
}
