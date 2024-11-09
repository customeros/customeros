import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type AddJobRoleMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  input: Types.JobRoleInput;
}>;

export type AddJobRoleMutation = {
  __typename?: 'Mutation';
  jobRole_Create: { __typename?: 'JobRole'; id: string };
};
