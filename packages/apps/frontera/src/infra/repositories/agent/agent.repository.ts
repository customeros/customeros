import { Transport } from '@infra/transport';

import AgentDocument from './queries/agent.graphql';
import AgentsDocument from './queries/agents.graphql';
import { AgentsQuery } from './queries/agents.generated';
import SaveAgentDocument from './mutations/saveAgent.graphql';
import { AgentQuery, AgentQueryVariables } from './queries/agent.generated';
import {
  SaveAgentMutation,
  SaveAgentMutationVariables,
} from './mutations/saveAgent.generated';

export class AgentRepository {
  private transport = Transport.getInstance();

  constructor() {}

  public async getAgent(payload: AgentQueryVariables) {
    return this.transport.graphql.request<AgentQuery, AgentQueryVariables>(
      AgentDocument,
      payload,
    );
  }

  public async getAgents() {
    return this.transport.graphql.request<AgentsQuery>(AgentsDocument);
  }

  public async saveAgent(payload: SaveAgentMutationVariables) {
    return this.transport.graphql.request<
      SaveAgentMutation,
      SaveAgentMutationVariables
    >(SaveAgentDocument, payload);
  }
}
