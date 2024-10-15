import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowSenderDeleteMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type FlowSenderDeleteMutation = {
  __typename?: 'Mutation';
  flowSender_Delete: { __typename?: 'Result'; result: boolean };
};
