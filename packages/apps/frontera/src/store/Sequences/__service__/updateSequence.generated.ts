import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateSequenceMutationVariables = Types.Exact<{
  input: Types.FlowSequenceStoreInput;
}>;

export type UpdateSequenceMutation = {
  __typename?: 'Mutation';
  flow_sequence_store: {
    __typename?: 'FlowSequence';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
