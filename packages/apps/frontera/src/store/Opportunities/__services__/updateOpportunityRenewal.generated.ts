import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOpportunityRenewalMutationVariables = Types.Exact<{
  input: Types.OpportunityRenewalUpdateInput;
}>;

export type UpdateOpportunityRenewalMutation = {
  __typename?: 'Mutation';
  opportunityRenewalUpdate: {
    __typename?: 'Opportunity';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
