import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type OrganizationQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type OrganizationQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    name: string;
    stage?: Types.OrganizationStage | null;
    description?: string | null;
    industry?: string | null;
    website?: string | null;
    domains: Array<string>;
    isCustomer?: boolean | null;
    logo?: string | null;
    icon?: string | null;
    relationship?: Types.OrganizationRelationship | null;
    leadSource?: string | null;
    valueProposition?: string | null;
    employees?: any | null;
    yearFounded?: any | null;
    public?: boolean | null;
    metadata: {
      __typename?: 'Metadata';
      id: string;
      created: any;
      lastUpdated: any;
    };
    parentCompanies: Array<{
      __typename?: 'LinkedOrganization';
      organization: {
        __typename?: 'Organization';
        name: string;
        metadata: { __typename?: 'Metadata'; id: string };
      };
    }>;
    contracts?: Array<{
      __typename?: 'Contract';
      metadata: { __typename?: 'Metadata'; id: string };
    }> | null;
    owner?: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
      name?: string | null;
    } | null;
    tags?: Array<{
      __typename?: 'Tag';
      id?: string | null;
      name: string;
      createdAt?: any | null;
      updatedAt?: any | null;
      appSource?: string | null;
    }> | null;
    socialMedia: Array<{
      __typename?: 'Social';
      id: string;
      url: string;
      followersCount: any;
    }>;
    enrichDetails: {
      __typename?: 'EnrichDetails';
      enrichedAt?: any | null;
      failedAt?: any | null;
      requestedAt?: any | null;
    };
    accountDetails?: {
      __typename?: 'OrgAccountDetails';
      churned?: any | null;
      ltv?: number | null;
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
  } | null;
};
