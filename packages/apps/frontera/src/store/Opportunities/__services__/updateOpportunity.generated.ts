import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOpportunityMutationVariables = Types.Exact<{
  input: Types.OpportunityUpdateInput;
}>;

export type UpdateOpportunityMutation = {
  __typename?: 'Mutation';
  opportunity_Update: { __typename?: 'Opportunity'; id: string };
};
