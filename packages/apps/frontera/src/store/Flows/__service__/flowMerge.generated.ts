import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowMergeMutationVariables = Types.Exact<{
  input: Types.FlowMergeInput;
}>;

export type FlowMergeMutation = {
  __typename?: 'Mutation';
  flow_Merge: {
    __typename?: 'Flow';
    nodes: string;
    edges: string;
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
