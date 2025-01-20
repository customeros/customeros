import { Transport } from '@infra/transport';

import AgentDocument from './queries/agent.graphql';
import AgentsDocument from './queries/agents.graphql';
import SaveAgentDocument from './mutations/saveAgent.graphql';
import { AgentQuery, AgentQueryVariables } from './queries/agent.generated';
import { AgentsQuery, AgentsQueryVariables } from './queries/agents.generated';
import {
  SaveAgentMutation,
  SaveAgentMutationVariables,
} from './mutations/saveAgent.generated';

export class AgentRepository {
  private instance: AgentRepository | null = null;
  private transport = Transport.getInstance();

  constructor() {
    if (!this.instance) {
      this.instance = new AgentRepository();
    }

    return this.instance;
  }

  public async getAgent(payload: AgentQueryVariables) {
    return this.transport.graphql.request<AgentQuery, AgentQueryVariables>(
      AgentDocument,
      payload,
    );
  }

  public async getAgents(payload: AgentsQueryVariables) {
    return this.transport.graphql.request<AgentsQuery, AgentsQueryVariables>(
      AgentsDocument,
      payload,
    );
  }

  public async saveAgent(payload: SaveAgentMutationVariables) {
    return this.transport.graphql.request<
      SaveAgentMutation,
      SaveAgentMutationVariables
    >(SaveAgentDocument, payload);
  }
}
