import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateTableViewDefSharedMutationVariables = Types.Exact<{
  input: Types.TableViewDefUpdateInput;
}>;


export type UpdateTableViewDefSharedMutation = { __typename?: 'Mutation', tableViewDef_UpdateShared: { __typename?: 'TableViewDef', id: string } };
