import { observable } from 'mobx';
import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { AgentService } from '@domain/services/agent/agent.service';

export class ArchiveAgentUsecase {
  private root = RootStore.getInstance();
  private service = new AgentService();

  @observable accessor isSaving: boolean = false;

  constructor(private id?: string) {
    this.execute = this.execute.bind(this);
  }

  get agent() {
    return this.id ? this.root.agents.getById(this.id) : null;
  }

  async execute(onSuccess?: () => void) {
    const agent = this.agent;

    if (!agent) {
      console.error('ArchiveAgentUsecase.execute: Could not find agent');
      this.isSaving = false;

      return;
    }

    this.isSaving = true;

    const span = Tracer.span('ArchiveAgentUsecase.execute', {
      payload: {
        agent,
      },
    });

    const [res, err] = await this.service.archiveAgent(agent);

    if (err) {
      console.error(
        'ArchiveAgentUsecase.execute: Could not update agent icon',
        err,
      );
      this.isSaving = false;

      this.root.ui.toastError(
        'We could not archive this agent',
        `archive-agent-${err}`,
      );

      return;
    }

    if (res?.agent_Delete) {
      onSuccess && onSuccess();
      this.isSaving = false;
      this.root.ui.commandMenu.setOpen(false);
      this.root.ui.toastSuccess('Agent archived', `${agent.id}-archived`);

      setTimeout(() => {
        this.root.agents.value.delete(agent.id);
      }, 0);
    }

    span.end();
  }
}
