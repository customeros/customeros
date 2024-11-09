import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateContactRoleMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  input: Types.JobRoleUpdateInput;
}>;

export type UpdateContactRoleMutation = {
  __typename?: 'Mutation';
  jobRole_Update: { __typename?: 'JobRole'; id: string };
};
