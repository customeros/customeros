import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type HideOrganizationsMutationVariables = Types.Exact<{
  ids: Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input'];
}>;

export type HideOrganizationsMutation = {
  __typename?: 'Mutation';
  organization_HideAll?: { __typename?: 'Result'; result: boolean } | null;
};
