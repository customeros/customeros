import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ArchiveTableViewDefMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type ArchiveTableViewDefMutation = {
  __typename?: 'Mutation';
  tableViewDef_Archive: { __typename?: 'ActionResponse'; accepted: boolean };
};
