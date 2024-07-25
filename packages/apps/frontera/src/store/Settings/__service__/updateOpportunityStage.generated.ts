import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOpportunityStageMutationVariables = Types.Exact<{
  input: Types.TenantSettingsOpportunityStageConfigurationInput;
}>;

export type UpdateOpportunityStageMutation = {
  __typename?: 'Mutation';
  tenant_UpdateSettingsOpportunityStage: {
    __typename?: 'ActionResponse';
    accepted: boolean;
  };
};
