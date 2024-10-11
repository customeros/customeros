import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type BulkUpdateOpportunityRenewalMutationVariables = Types.Exact<{
  input: Types.OpportunityRenewalUpdateAllForOrganizationInput;
}>;

export type BulkUpdateOpportunityRenewalMutation = {
  __typename?: 'Mutation';
  opportunityRenewal_UpdateAllForOrganization: {
    __typename?: 'Organization';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
