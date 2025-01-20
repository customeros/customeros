import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GlobalCacheQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GlobalCacheQuery = {
  __typename?: 'Query';
  global_Cache: {
    __typename?: 'GlobalCache';
    cdnLogoUrl: string;
    mailboxes: Array<string>;
    isOwner: boolean;
    minARRForecastValue: number;
    maxARRForecastValue: number;
    contractsExist: boolean;
    isFirstLogin: boolean;
    user: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
      mailboxes: Array<string>;
      emails?: Array<{
        __typename?: 'Email';
        email?: string | null;
        rawEmail?: string | null;
        primary: boolean;
      }> | null;
      onboarding: {
        __typename?: 'UserOnboardingDetails';
        showOnboardingPage: boolean;
        onboardingInboundStepCompleted: boolean;
        onboardingOutboundStepCompleted: boolean;
        onboardingCrmStepCompleted: boolean;
        onboardingMailstackStepCompleted: boolean;
      };
    };
    inactiveEmailTokens: Array<{
      __typename?: 'GlobalCacheEmailToken';
      email: string;
      provider: string;
    }>;
    activeEmailTokens: Array<{
      __typename?: 'GlobalCacheEmailToken';
      email: string;
      provider: string;
    }>;
  };
};
