import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type SetOrganizationOwnerMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  userId: Types.Scalars['ID']['input'];
}>;

export type SetOrganizationOwnerMutation = {
  __typename?: 'Mutation';
  organization_SetOwner: { __typename?: 'Organization'; id: string };
};
