import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type SaveOpportunityMutationVariables = Types.Exact<{
  input: Types.OpportunitySaveInput;
}>;

export type SaveOpportunityMutation = {
  __typename?: 'Mutation';
  opportunity_Save: {
    __typename?: 'Opportunity';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
