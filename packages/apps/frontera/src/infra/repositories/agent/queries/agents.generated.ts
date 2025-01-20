import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type AgentsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type AgentsQuery = {
  __typename?: 'Query';
  agents: Array<{
    __typename?: 'Agent';
    id: string;
    type: Types.AgentType;
    tenant: string;
    name: string;
    goal: string;
    isActive: boolean;
    flowId?: string | null;
    visible: boolean;
    createdAt: any;
    updatedAt: any;
    error?: string | null;
    color: string;
    icon: string;
    capabilities: Array<{
      __typename?: 'Capability';
      id: string;
      type: Types.CapabilityType;
      name: string;
      action: string;
      optional: boolean;
      values: string;
      errors?: string | null;
    }>;
  }>;
};
