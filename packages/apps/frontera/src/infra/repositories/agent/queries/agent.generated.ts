import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type AgentQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type AgentQuery = {
  __typename?: 'Query';
  agent?: {
    __typename?: 'Agent';
    id: string;
    type: Types.AgentType;
    name: string;
    scope: Types.AgentScope;
    goal: string;
    isActive: boolean;
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
  } | null;
};
