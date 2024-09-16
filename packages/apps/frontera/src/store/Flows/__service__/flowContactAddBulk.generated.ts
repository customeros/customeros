import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowContactAddBulkMutationVariables = Types.Exact<{
  flowId: Types.Scalars['ID']['input'];
  contactId: Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input'];
}>;

export type FlowContactAddBulkMutation = {
  __typename?: 'Mutation';
  flowContact_AddBulk: { __typename?: 'Result'; result: boolean };
};
