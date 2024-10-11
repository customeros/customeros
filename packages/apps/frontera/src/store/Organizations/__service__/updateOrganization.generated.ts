import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOrganizationMutationVariables = Types.Exact<{
  input: Types.OrganizationUpdateInput;
}>;

export type UpdateOrganizationMutation = {
  __typename?: 'Mutation';
  organization_Update: {
    __typename?: 'Organization';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
