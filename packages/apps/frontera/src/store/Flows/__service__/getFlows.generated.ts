import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GetFlowsQuery = {
  __typename?: 'Query';
  flows: Array<{
    __typename?: 'Flow';
    firstStartedAt?: any | null;
    name: string;
    edges: string;
    nodes: string;
    status: Types.FlowStatus;
    tableViewDefId: string;
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
    participants: Array<{
      __typename?: 'FlowParticipant';
      status: Types.FlowParticipantStatus;
      entityId: string;
      entityType: string;
      requirementsUnmeet: Array<Types.FlowParticipantRequirementsUnmeet>;
      executions: Array<{
        __typename?: 'FlowActionExecution';
        status: Types.FlowActionExecutionStatus;
        scheduledAt?: any | null;
        executedAt?: any | null;
        error?: string | null;
        metadata: { __typename?: 'Metadata'; id: string };
        action: {
          __typename?: 'FlowAction';
          action: Types.FlowActionType;
          metadata: { __typename?: 'Metadata'; id: string };
        };
      }>;
      metadata: { __typename?: 'Metadata'; id: string };
    }>;
  }>;
};
