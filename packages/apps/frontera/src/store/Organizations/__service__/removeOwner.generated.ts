import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type RemoveOrganizationOwnerMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
}>;

export type RemoveOrganizationOwnerMutation = {
  __typename?: 'Mutation';
  organization_UnsetOwner: { __typename?: 'Organization'; id: string };
};
