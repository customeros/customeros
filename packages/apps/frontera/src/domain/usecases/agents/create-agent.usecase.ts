import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { AgentService } from '@domain/services/agent/agent.service';

import { AgentType } from '@graphql/types';

export class CreateAgentUsecase {
  private root = RootStore.getInstance();
  private service = new AgentService();

  constructor(private navigateToAgent: (id: string) => void) {
    this.execute = this.execute.bind(this);
  }

  public async execute(type: AgentType) {
    const span = Tracer.span('CreateAgentUsecase.execute', {
      payload: type,
    });

    const [res, err] = await this.service.createAgent(type);

    if (err) {
      console.error('CreateAgentUsecase.execute: Could not create agent', err);
      span.end();

      return;
    }

    if (res?.agent_Save) {
      this.root.agents.addOne(res.agent_Save);
      this.root.ui.commandMenu.setType('AgentCommands');
      this.root.ui.commandMenu.setOpen(false);

      this.navigateToAgent(res.agent_Save.id);
    }

    span.end();
  }
}
