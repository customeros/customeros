import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateContactMutationVariables = Types.Exact<{
  input: Types.ContactInput;
  organizationId: Types.Scalars['ID']['input'];
}>;

export type CreateContactMutation = {
  __typename?: 'Mutation';
  contact_CreateForOrganization: { __typename?: 'Contact'; id: string };
};
