import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type RemoveDomainMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  domain: Types.Scalars['String']['input'];
}>;

export type RemoveDomainMutation = {
  __typename?: 'Mutation';
  organization_RemoveDomain: {
    __typename?: 'ActionResponse';
    accepted: boolean;
  };
};
