import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GetFlowsQuery = {
  __typename?: 'Query';
  flows: Array<{
    __typename?: 'Flow';
    name: string;
    edges: string;
    nodes: string;
    status: Types.FlowStatus;
    metadata: { __typename?: 'Metadata'; id: string };
    statistics: {
      __typename?: 'FlowStatistics';
      total: any;
      pending: any;
      completed: any;
      goalAchieved: any;
    };
    contacts: Array<{
      __typename?: 'FlowContact';
      status: Types.FlowContactStatus;
      scheduledAction?: string | null;
      scheduledAt?: any | null;
      metadata: { __typename?: 'Metadata'; id: string };
      contact: {
        __typename?: 'Contact';
        metadata: { __typename?: 'Metadata'; id: string };
      };
    }>;
  }>;
};
