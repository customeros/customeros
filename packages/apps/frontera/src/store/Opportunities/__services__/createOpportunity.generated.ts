import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateOpportunityMutationVariables = Types.Exact<{
  input: Types.OpportunityCreateInput;
}>;

export type CreateOpportunityMutation = {
  __typename?: 'Mutation';
  opportunity_Create: {
    __typename?: 'Opportunity';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
