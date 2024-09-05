import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateSequenceMutationVariables = Types.Exact<{
  input: Types.FlowSequenceStoreInput;
}>;

export type CreateSequenceMutation = {
  __typename?: 'Mutation';
  flow_sequence_store: {
    __typename?: 'FlowSequence';
    name: string;
    description: string;
    status: Types.FlowSequenceStatus;
    flow: {
      __typename?: 'Flow';
      name: string;
      description: string;
      status: Types.FlowStatus;
      metadata: { __typename?: 'Metadata'; id: string };
    };
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
