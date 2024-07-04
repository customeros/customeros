import { gql } from 'graphql-request';
import { Transport } from '@store/transport';

import {
  Workflow,
  WorkflowType,
  WorkflowUpdateInput,
} from '@shared/types/__generated__/graphql.types';

class WorkFlowService {
  private static instance: WorkFlowService;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport) {
    if (!WorkFlowService.instance) {
      WorkFlowService.instance = new WorkFlowService(transport);
    }

    return WorkFlowService.instance;
  }

  async getWorkFlowsByType() {
    return this.transport.graphql.request<
      WORKFLOW_QUERY_RESPONSE,
      WORKFLOW_QUERY_INPUT
    >(WORKFLOW_QUERY, { workflowType: WorkflowType.IdealCustomerProfile });
  }

  async updateWorkFlow(input: WorkflowUpdateInput) {
    return this.transport.graphql.request<
      UPDATE_WORKFLOW_MUTATION_RESPONSE,
      UPDATE_WORKFLOW_MUTATION_INPUT
    >(UPDATE_WORKFLOW_MUTATION, { input });
  }
}

type UPDATE_WORKFLOW_MUTATION_RESPONSE = {
  update_workflow: Workflow;
};
type UPDATE_WORKFLOW_MUTATION_INPUT = {
  input: WorkflowUpdateInput;
};

const UPDATE_WORKFLOW_MUTATION = gql`
  mutation update_workflow($input: WorkflowUpdateInput!) {
    workflow_Update(input: $input) {
      id
    }
  }
`;

type WORKFLOW_QUERY_RESPONSE = {
  workflow_ByType: Workflow[];
};

type WORKFLOW_QUERY_INPUT = {
  workflowType: WorkflowType;
};

const WORKFLOW_QUERY = gql`
  query workflow_ByType($workflowType: WorkflowType!) {
    workflow_ByType(workflowType: $workflowType) {
      id
      name
      type
      condition
      live
    }
  }
`;

export { WorkFlowService };
