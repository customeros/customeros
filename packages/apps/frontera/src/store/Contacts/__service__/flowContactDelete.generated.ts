import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowContactDeleteMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type FlowContactDeleteMutation = {
  __typename?: 'Mutation';
  flowContact_Delete: { __typename?: 'Result'; result: boolean };
};
