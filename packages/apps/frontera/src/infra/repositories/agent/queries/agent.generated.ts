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
  } | null;
};
