import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type SaveAgentMutationVariables = Types.Exact<{
  input: Types.AgentSaveInput;
}>;

export type SaveAgentMutation = {
  __typename?: 'Mutation';
  agent_Save: {
    __typename?: 'Agent';
    id: string;
    type: Types.AgentType;
    name: string;
    goal: string;
    isActive: boolean;
    visible: boolean;
    createdAt: any;
    updatedAt: any;
    error?: string | null;
    color: string;
    scope: Types.AgentScope;
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
  };
};
