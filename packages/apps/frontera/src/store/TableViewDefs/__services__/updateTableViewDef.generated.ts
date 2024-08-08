import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateTableViewDefMutationVariables = Types.Exact<{
  input: Types.TableViewDefUpdateInput;
}>;

export type UpdateTableViewDefMutation = {
  __typename?: 'Mutation';
  tableViewDef_Update: { __typename?: 'TableViewDef'; id: string };
};
