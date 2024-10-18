import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type MergeOrganizationsMutationVariables = Types.Exact<{
  primaryOrganizationId: Types.Scalars['ID']['input'];
  mergedOrganizationIds:
    | Array<Types.Scalars['ID']['input']>
    | Types.Scalars['ID']['input'];
}>;

export type MergeOrganizationsMutation = {
  __typename?: 'Mutation';
  organization_Merge: { __typename?: 'Organization'; id: string };
};
