import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type AgentsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type AgentsQuery = {
  __typename?: 'Query';
  agents: Array<{
    __typename?: 'Agent';
    id: string;
    type: Types.AgentType;
    name: string;
    goal: string;
    isActive: boolean;
    scope: Types.AgentScope;
    visible: boolean;
    createdAt: any;
    updatedAt: any;
    error?: string | null;
    color: string;
    icon: string;
    isConfigured: boolean;
    listeners: Array<{
      __typename?: 'AgentListener';
      id: string;
      type: Types.AgentListenerEvent;
      name: string;
      active: boolean;
      config: string;
      errors?: string | null;
    }>;
    capabilities: Array<{
      __typename?: 'Capability';
      id: string;
      type: Types.CapabilityType;
      name: string;
      active: boolean;
      config: string;
      errors?: string | null;
    }>;
  }>;
};
