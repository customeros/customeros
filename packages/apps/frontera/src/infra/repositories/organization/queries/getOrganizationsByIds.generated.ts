import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

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
    industryName?: string | null;
    industryCode?: string | null;
    market?: Types.Market | null;
    website?: string | null;
    logoUrl?: string | null;
    iconUrl?: string | null;
    public?: boolean | null;
    stage?: Types.OrganizationStage | null;
    relationship?: Types.OrganizationRelationship | null;
    lastFundingRound?: Types.FundingRound | null;
    leadSource?: string | null;
    slackChannelId?: string | null;
    employees?: any | null;
    yearFounded?: any | null;
    enrichedAt?: any | null;
    enrichedFailedAt?: any | null;
    enrichedRequestedAt?: any | null;
    ltv?: number | null;
    hide: boolean;
    domains: Array<string>;
    icpFit?: Types.IcpFit | null;
    icpFitReasons: Array<string>;
    icpFitUpdatedAt?: any | null;
    wrongIndustry: boolean;
    createdAt: any;
    updatedAt: any;
    churnedAt?: any | null;
    customerOsId: string;
    referenceId: string;
    renewalSummaryArrForecast?: number | null;
    renewalSummaryMaxArrForecast?: number | null;
    renewalSummaryRenewalLikelihood?: Types.OpportunityRenewalLikelihood | null;
    renewalSummaryNextRenewalAt?: any | null;
    onboardingStatus: Types.OnboardingStatus;
    onboardingStatusUpdatedAt?: any | null;
    onboardingComments?: string | null;
    lastTouchPointAt?: any | null;
    lastTouchPointType?: Types.LastTouchpointType | null;
    contactCount?: number | null;
    parentId?: string | null;
    parentName?: string | null;
    contracts: Array<string>;
    contacts: Array<string>;
    subsidiaries: Array<string>;
    domainsDetails: Array<{
      __typename?: 'Domain';
      domain: string;
      primary?: boolean | null;
      primaryDomain?: string | null;
    }>;
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
      colorCode: string;
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
  }>;
};
