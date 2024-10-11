import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type RemoveTagFromOrganizationMutationVariables = Types.Exact<{
  input: Types.OrganizationTagInput;
}>;

export type RemoveTagFromOrganizationMutation = {
  __typename?: 'Mutation';
  organization_RemoveTag: { __typename?: 'ActionResponse'; accepted: boolean };
};
