import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type SaveOrganizationMutationVariables = Types.Exact<{
  input: Types.OrganizationSaveInput;
}>;

export type SaveOrganizationMutation = {
  __typename?: 'Mutation';
  organization_Save: {
    __typename?: 'Organization';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
