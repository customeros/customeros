import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type AddDomainMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  domain: Types.Scalars['String']['input'];
}>;

export type AddDomainMutation = {
  __typename?: 'Mutation';
  organization_AddDomain: { __typename?: 'ActionResponse'; accepted: boolean };
};
