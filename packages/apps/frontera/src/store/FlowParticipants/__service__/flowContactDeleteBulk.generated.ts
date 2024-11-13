import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowContactDeleteBulkMutationVariables = Types.Exact<{
  id: Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input'];
}>;

export type FlowContactDeleteBulkMutation = {
  __typename?: 'Mutation';
  flowContact_DeleteBulk: { __typename?: 'Result'; result: boolean };
};
