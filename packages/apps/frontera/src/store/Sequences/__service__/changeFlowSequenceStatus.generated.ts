import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ChangeFlowSequenceStatusMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  stage: Types.FlowSequenceStatus;
}>;

export type ChangeFlowSequenceStatusMutation = {
  __typename?: 'Mutation';
  flow_sequence_ChangeStatus: {
    __typename?: 'FlowSequence';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
