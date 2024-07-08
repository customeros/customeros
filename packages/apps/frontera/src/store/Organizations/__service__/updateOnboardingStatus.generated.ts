import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateOnboardingStatusMutationVariables = Types.Exact<{
  input: Types.OnboardingStatusInput;
}>;

export type UpdateOnboardingStatusMutation = {
  __typename?: 'Mutation';
  organization_UpdateOnboardingStatus: {
    __typename?: 'Organization';
    id: string;
    accountDetails?: {
      __typename?: 'OrgAccountDetails';
      onboarding?: {
        __typename?: 'OnboardingDetails';
        status: Types.OnboardingStatus;
        comments?: string | null;
      } | null;
    } | null;
  };
};
