import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOpportunityToCloseWonMutationVariables = Types.Exact<{
  opportunityId: Types.Scalars['ID']['input'];
}>;

export type UpdateOpportunityToCloseWonMutation = {
  __typename?: 'Mutation';
  opportunity_CloseWon: { __typename?: 'ActionResponse'; accepted: boolean };
};
