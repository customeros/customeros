import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GlobalCacheQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GlobalCacheQuery = {
  __typename?: 'Query';
  global_Cache: {
    __typename?: 'GlobalCache';
    cdnLogoUrl: string;
    contactRegions: Array<string>;
    contactCities: Array<string>;
    mailboxes: Array<string>;
    isPlatformOwner: boolean;
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
      mailboxesV2: Array<{
        __typename?: 'Mailbox';
        userId?: string | null;
        mailbox: string;
      }>;
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
