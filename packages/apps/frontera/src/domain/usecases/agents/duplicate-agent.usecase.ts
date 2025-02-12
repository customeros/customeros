import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { AgentService } from '@domain/services/agent/agent.service';

import { AgentType } from '@graphql/types';

export class DuplicateAgentUsecase {
  private root = RootStore.getInstance();
  private service = new AgentService();

  constructor() {
    this.execute = this.execute.bind(this);
  }

  public async execute(type: AgentType) {
    const span = Tracer.span('DuplicateAgentUsecase.execute', {
      payload: type,
    });

    const [res, err] = await this.service.createAgent(type);

    if (err) {
      console.error(
        'DuplicateAgentUsecase.execute: Could not duplicate agent',
        err,
      );
      span.end();

      return;
    }

    if (res?.agent_Save) {
      this.root.agents.addOne(res.agent_Save);
      this.root.ui.commandMenu.setType('AgentCommands');
      this.root.ui.commandMenu.toggle();
    }

    span.end();
  }
}
