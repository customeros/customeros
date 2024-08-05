import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOpportunityOwnerMutationVariables = Types.Exact<{
  opportunityId: Types.Scalars['ID']['input'];
  userID: Types.Scalars['ID']['input'];
}>;

export type UpdateOpportunityOwnerMutation = {
  __typename?: 'Mutation';
  opportunity_SetOwner: { __typename?: 'ActionResponse'; accepted: boolean };
};
