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

export type GetOrganizationsByIdsQueryVariables = Types.Exact<{
  ids?: Types.InputMaybe<
    Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input']
  >;
}>;

export type GetOrganizationsByIdsQuery = {
  __typename?: 'Query';
  ui_organizations: Array<{
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
  }>;
};
