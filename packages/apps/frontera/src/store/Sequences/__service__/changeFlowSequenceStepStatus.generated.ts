import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ChangeFlowSequenceStepStatusMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  status: Types.FlowSequenceStepStatus;
}>;

export type ChangeFlowSequenceStepStatusMutation = {
  __typename?: 'Mutation';
  flow_sequence_step_ChangeStatus: {
    __typename?: 'FlowSequenceStep';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
