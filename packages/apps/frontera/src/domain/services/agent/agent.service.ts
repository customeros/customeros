import { type Agent } from '@store/Agents/Agent.dto';
import { AgentRepository } from '@infra/repositories/agent';

import { unwrap } from '@utils/unwrap';

export class AgentService {
  private repo = new AgentRepository();

  public async saveAgent(agent: Agent) {
    return await unwrap(this.repo.saveAgent({ input: agent.value }));
  }
}
