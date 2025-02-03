import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetOrganizationsQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;

export type GetOrganizationsQuery = {
  __typename?: 'Query';
  dashboardView_Organizations?: {
    __typename?: 'OrganizationPage';
    totalElements: any;
    totalAvailable: any;
    content: Array<{
      __typename?: 'Organization';
      name: string;
      notes?: string | null;
      hide: boolean;
      wrongIndustry: boolean;
      icpFit?: Types.IcpFit | null;
      stage?: Types.OrganizationStage | null;
      description?: string | null;
      industry?: string | null;
      market?: Types.Market | null;
      website?: string | null;
      domains: Array<string>;
      logo?: string | null;
      icon?: string | null;
      relationship?: Types.OrganizationRelationship | null;
      lastFundingRound?: Types.FundingRound | null;
      leadSource?: string | null;
      slackChannelId?: string | null;
      public?: boolean | null;
      employees?: any | null;
      customerOsId: string;
      yearFounded?: any | null;
      metadata: {
        __typename?: 'Metadata';
        id: string;
        created: any;
        lastUpdated: any;
      };
      contracts?: Array<{
        __typename?: 'Contract';
        metadata: { __typename?: 'Metadata'; id: string };
      }> | null;
      parentCompanies: Array<{
        __typename?: 'LinkedOrganization';
        organization: {
          __typename?: 'Organization';
          name: string;
          metadata: { __typename?: 'Metadata'; id: string };
        };
      }>;
      owner?: {
        __typename?: 'User';
        id: string;
        firstName: string;
        lastName: string;
        name?: string | null;
      } | null;
      contacts: {
        __typename?: 'ContactsPage';
        content: Array<{
          __typename?: 'Contact';
          id: string;
          metadata: { __typename?: 'Metadata'; id: string };
        }>;
      };
      enrichDetails: {
        __typename?: 'EnrichDetails';
        enrichedAt?: any | null;
        failedAt?: any | null;
        requestedAt?: any | null;
      };
      socialMedia: Array<{
        __typename?: 'Social';
        id: string;
        url: string;
        followersCount: any;
        alias: string;
      }>;
      tags?: Array<{
        __typename?: 'Tag';
        name: string;
        entityType: Types.EntityType;
        metadata: { __typename?: 'Metadata'; id: string };
      }> | null;
      accountDetails?: {
        __typename?: 'OrgAccountDetails';
        ltv?: number | null;
        churned?: any | null;
        renewalSummary?: {
          __typename?: 'RenewalSummary';
          arrForecast?: number | null;
          maxArrForecast?: number | null;
          renewalLikelihood?: Types.OpportunityRenewalLikelihood | null;
          nextRenewalDate?: any | null;
        } | null;
        onboarding?: {
          __typename?: 'OnboardingDetails';
          status: Types.OnboardingStatus;
          comments?: string | null;
          updatedAt?: any | null;
        } | null;
      } | null;
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
      subsidiaries: Array<{
        __typename?: 'LinkedOrganization';
        organization: {
          __typename?: 'Organization';
          name: string;
          metadata: { __typename?: 'Metadata'; id: string };
          parentCompanies: Array<{
            __typename?: 'LinkedOrganization';
            organization: {
              __typename?: 'Organization';
              name: string;
              metadata: { __typename?: 'Metadata'; id: string };
            };
          }>;
        };
      }>;
      lastTouchpoint?: {
        __typename?: 'LastTouchpoint';
        lastTouchPointAt?: any | null;
        lastTouchPointType?: Types.LastTouchpointType | null;
      } | null;
    }>;
  } | null;
};
