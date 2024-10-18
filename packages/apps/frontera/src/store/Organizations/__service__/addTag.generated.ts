import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type AddTagsToOrganizationMutationVariables = Types.Exact<{
  input: Types.OrganizationTagInput;
}>;

export type AddTagsToOrganizationMutation = {
  __typename?: 'Mutation';
  organization_AddTag: { __typename?: 'ActionResponse'; accepted: boolean };
};
