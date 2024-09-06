import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateSequenceMutationVariables = Types.Exact<{
  input: Types.FlowSequenceUpdateInput;
}>;

export type UpdateSequenceMutation = {
  __typename?: 'Mutation';
  flow_sequence_Update: {
    __typename?: 'FlowSequence';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
