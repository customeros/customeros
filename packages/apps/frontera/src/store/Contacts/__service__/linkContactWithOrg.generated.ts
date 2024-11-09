import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type LinkOrganizationMutationVariables = Types.Exact<{
  input: Types.ContactOrganizationInput;
}>;

export type LinkOrganizationMutation = {
  __typename?: 'Mutation';
  contact_AddOrganizationById: { __typename?: 'Contact'; id: string };
};
