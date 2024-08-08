import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateTableViewDefMutationVariables = Types.Exact<{
  input: Types.TableViewDefCreateInput;
}>;

export type CreateTableViewDefMutation = {
  __typename?: 'Mutation';
  tableViewDef_Create: { __typename?: 'TableViewDef'; id: string };
};
