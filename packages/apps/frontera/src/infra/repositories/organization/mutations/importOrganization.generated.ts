import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type ImportOrganizationMutationVariables = Types.Exact<{
  globalOrganizationId: Types.Scalars['Int64']['input'];
  input?: Types.InputMaybe<Types.OrganizationSaveInputFromGlobalOrg>;
}>;

export type ImportOrganizationMutation = {
  __typename?: 'Mutation';
  organization_SaveByGlobalOrganization: {
    __typename?: 'OrganizationUiDetails';
    id: string;
  };
};
