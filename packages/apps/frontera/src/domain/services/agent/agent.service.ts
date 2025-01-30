import { Tracer } from '@infra/tracer';
import { type Agent } from '@store/Agents/Agent.dto';
import { AgentRepository } from '@infra/repositories/agent';

import { unwrap } from '@utils/unwrap';

export class AgentService {
  private repo = new AgentRepository();

  public async saveAgent(agent: Agent) {
    const span = Tracer.span('AgentService.saveAgent', {
      payload: agent.toPayload(),
    });

    const req = await unwrap(this.repo.saveAgent({ input: agent.toPayload() }));

    span.end();

    return req;
  }
}
