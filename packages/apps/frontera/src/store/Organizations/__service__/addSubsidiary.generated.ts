import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type AddSubsidiaryToOrganizationMutationVariables = Types.Exact<{
  input: Types.LinkOrganizationsInput;
}>;

export type AddSubsidiaryToOrganizationMutation = {
  __typename?: 'Mutation';
  organization_AddSubsidiary: {
    __typename?: 'Organization';
    metadata: { __typename?: 'Metadata'; id: string };
    subsidiaries: Array<{
      __typename?: 'LinkedOrganization';
      organization: {
        __typename?: 'Organization';
        name: string;
        metadata: { __typename?: 'Metadata'; id: string };
        locations: Array<{
          __typename?: 'Location';
          id: string;
          address?: string | null;
        }>;
      };
    }>;
  };
};
