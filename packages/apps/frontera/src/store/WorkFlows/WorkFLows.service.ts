import { gql } from 'graphql-request';
import { Transport } from '@store/transport';

import { Workflow } from '@shared/types/__generated__/graphql.types';

class WorkFlowsService {
  private static instance: WorkFlowsService;
  private transport: Transport;
  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport) {
    if (!WorkFlowsService.instance) {
      WorkFlowsService.instance = new WorkFlowsService(transport);
    }

    return WorkFlowsService.instance;
  }

  async getWorkFlows() {
    return this.transport.graphql.request<WORKFLOWS_QUERY_RESPONSE>(
      WORKFLOWS_QUERY,
    );
  }
}

type WORKFLOWS_QUERY_RESPONSE = {
  workflows: Workflow[];
};

const WORKFLOWS_QUERY = gql`
  query workFlows {
    workflows {
      id
      name
      type
      live
      condition
    }
  }
`;

export { WorkFlowsService };
