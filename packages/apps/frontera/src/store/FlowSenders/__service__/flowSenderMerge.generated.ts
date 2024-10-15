import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowSenderMergeMutationVariables = Types.Exact<{
  flowId: Types.Scalars['ID']['input'];
  input: Types.FlowSenderMergeInput;
}>;

export type FlowSenderMergeMutation = {
  __typename?: 'Mutation';
  flowSender_Merge: {
    __typename?: 'FlowSender';
    metadata: { __typename?: 'Metadata'; id: string };
    user?: { __typename?: 'User'; id: string } | null;
  };
};
