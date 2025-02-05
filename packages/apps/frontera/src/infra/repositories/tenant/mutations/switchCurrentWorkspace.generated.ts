import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type AdminSwitchCurrentWorkspaceMutationVariables = Types.Exact<{
  tenant: Types.Scalars['String']['input'];
}>;

export type AdminSwitchCurrentWorkspaceMutation = {
  __typename?: 'Mutation';
  admin_switchCurrentWorkspace: boolean;
};
