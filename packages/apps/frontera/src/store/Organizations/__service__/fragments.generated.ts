import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type OrganizationPlainFragment = {
  __typename?: 'OrganizationUiDetails';
  id: string;
  name: string;
  notes?: string | null;
  description?: string | null;
  industry?: string | null;
  market?: Types.Market | null;
  website?: string | null;
  logoUrl?: string | null;
  iconUrl?: string | null;
  public?: boolean | null;
  stage?: Types.OrganizationStage | null;
  relationship?: Types.OrganizationRelationship | null;
  lastFundingRound?: Types.FundingRound | null;
  leadSource?: string | null;
  valueProposition?: string | null;
  slackChannelId?: string | null;
  employees?: any | null;
  yearFounded?: any | null;
  enrichedAt?: any | null;
  enrichedFailedAt?: any | null;
  enrichedRequestedAt?: any | null;
  ltv?: number | null;
  churnedAt?: any | null;
  renewalSummaryArrForecast?: number | null;
  renewalSummaryMaxArrForecast?: number | null;
  renewalSummaryRenewalLikelihood?: Types.OpportunityRenewalLikelihood | null;
  renewalSummaryNextRenewalAt?: any | null;
  onboardingStatus: Types.OnboardingStatus;
  onboardingStatusUpdatedAt?: any | null;
  onboardingComments?: string | null;
  lastTouchPointAt?: any | null;
  contactCount?: number | null;
  parentId?: string | null;
  parentName?: string | null;
};

export type OrganizationNestedFragment = {
  __typename?: 'OrganizationUiDetails';
  contracts: Array<string>;
  contacts: Array<string>;
  subsidiaries: Array<string>;
  socialMedia: Array<{
    __typename?: 'Social';
    id: string;
    url: string;
    alias: string;
    followersCount: any;
  }>;
  tags: Array<{
    __typename?: 'Tag';
    name: string;
    entityType: Types.EntityType;
    metadata: { __typename?: 'Metadata'; id: string };
  }>;
  locations: Array<{
    __typename?: 'Location';
    id: string;
    name?: string | null;
    country?: string | null;
    region?: string | null;
    locality?: string | null;
    zip?: string | null;
    street?: string | null;
    postalCode?: string | null;
    houseNumber?: string | null;
    rawAddress?: string | null;
    countryCodeA2?: string | null;
    countryCodeA3?: string | null;
  }>;
  owner?: {
    __typename?: 'User';
    id: string;
    firstName: string;
    lastName: string;
    name?: string | null;
  } | null;
};
