import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type SaveJobRolesMutationVariables = Types.Exact<{
  input: Types.JobRoleSaveInput;
}>;

export type SaveJobRolesMutation = {
  __typename?: 'Mutation';
  jobRole_Save: string;
};
