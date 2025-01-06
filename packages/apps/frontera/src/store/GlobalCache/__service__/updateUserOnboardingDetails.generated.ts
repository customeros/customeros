import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UserUpdateOnboardingDetailsMutationVariables = Types.Exact<{
  input: Types.UserOnboardingDetailsInput;
}>;


export type UserUpdateOnboardingDetailsMutation = { __typename?: 'Mutation', user_UpdateOnboardingDetails: { __typename?: 'User', id: string, firstName: string, lastName: string, emails?: Array<{ __typename?: 'Email', email?: string | null, rawEmail?: string | null, primary: boolean }> | null, onboarding: { __typename?: 'UserOnboardingDetails', showOnboardingPage: boolean, onboardingInboundStepCompleted: boolean, onboardingOutboundStepCompleted: boolean, onboardingCrmStepCompleted: boolean, onboardingMailstackStepCompleted: boolean } } };
