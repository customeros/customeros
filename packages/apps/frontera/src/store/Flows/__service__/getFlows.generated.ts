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
    senders: Array<{
      __typename?: 'FlowSender';
      metadata: { __typename?: 'Metadata'; id: string };
      user?: { __typename?: 'User'; id: string } | null;
      flow?: {
        __typename?: 'Flow';
        metadata: { __typename?: 'Metadata'; id: string };
      } | null;
    }>;
    statistics: {
      __typename?: 'FlowStatistics';
      total: any;
      onHold: any;
      ready: any;
      scheduled: any;
      inProgress: any;
      completed: any;
      goalAchieved: any;
    };
    contacts: Array<{
      __typename?: 'FlowContact';
      status: Types.FlowParticipantStatus;
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
