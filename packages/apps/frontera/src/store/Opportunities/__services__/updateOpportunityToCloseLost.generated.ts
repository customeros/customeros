import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOpportunityToCloseLostMutationVariables = Types.Exact<{
  opportunityId: Types.Scalars['ID']['input'];
}>;

export type UpdateOpportunityToCloseLostMutation = {
  __typename?: 'Mutation';
  opportunity_CloseLost: { __typename?: 'ActionResponse'; accepted: boolean };
};
