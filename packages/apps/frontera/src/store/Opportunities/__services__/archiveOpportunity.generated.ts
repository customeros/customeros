import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ArchiveOpportunityMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type ArchiveOpportunityMutation = {
  __typename?: 'Mutation';
  opportunity_Archive: { __typename?: 'ActionResponse'; accepted: boolean };
};
