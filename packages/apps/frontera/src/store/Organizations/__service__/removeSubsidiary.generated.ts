import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type RemoveSubsidiaryToOrganizationMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  subsidiaryId: Types.Scalars['ID']['input'];
}>;

export type RemoveSubsidiaryToOrganizationMutation = {
  __typename?: 'Mutation';
  organization_RemoveSubsidiary: {
    __typename?: 'Organization';
    id: string;
    subsidiaries: Array<{
      __typename?: 'LinkedOrganization';
      organization: {
        __typename?: 'Organization';
        id: string;
        name: string;
        locations: Array<{
          __typename?: 'Location';
          id: string;
          address?: string | null;
        }>;
      };
    }>;
  };
};
