import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowParticipantQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type GetFlowParticipantQuery = { __typename?: 'Query', flowParticipant: { __typename?: 'FlowParticipant', status: Types.FlowParticipantStatus, entityId: string, entityType: string, executions: Array<{ __typename?: 'FlowActionExecution', status: Types.FlowActionExecutionStatus, scheduledAt?: any | null, executedAt?: any | null, error?: string | null, metadata: { __typename?: 'Metadata', id: string }, action: { __typename?: 'FlowAction', action: Types.FlowActionType } }>, metadata: { __typename?: 'Metadata', id: string } } };
