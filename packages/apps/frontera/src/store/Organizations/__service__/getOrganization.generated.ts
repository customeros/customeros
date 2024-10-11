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
    metadata: { __typename?: 'Metadata'; id: string; created: any };
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
      id: string;
      name: string;
      createdAt: any;
      updatedAt: any;
      appSource: string;
    }> | null;
    socialMedia: Array<{
      __typename?: 'Social';
      id: string;
      url: string;
      followersCount: any;
    }>;
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
      lastTouchPointTimelineEventId?: string | null;
      lastTouchPointAt?: any | null;
      lastTouchPointType?: Types.LastTouchpointType | null;
      lastTouchPointTimelineEvent?:
        | {
            __typename: 'Action';
            id: string;
            actionType: Types.ActionType;
            createdAt: any;
            source: Types.DataSource;
            createdBy?: {
              __typename?: 'User';
              id: string;
              firstName: string;
              lastName: string;
            } | null;
          }
        | {
            __typename: 'InteractionEvent';
            id: string;
            channel: string;
            eventType?: string | null;
            externalLinks: Array<{
              __typename?: 'ExternalSystem';
              type: Types.ExternalSystemType;
            }>;
            sentBy: Array<
              | {
                  __typename: 'ContactParticipant';
                  contactParticipant: {
                    __typename?: 'Contact';
                    id: string;
                    name?: string | null;
                    firstName?: string | null;
                    lastName?: string | null;
                  };
                }
              | {
                  __typename: 'EmailParticipant';
                  type?: string | null;
                  emailParticipant: {
                    __typename?: 'Email';
                    id: string;
                    email?: string | null;
                    rawEmail?: string | null;
                  };
                }
              | {
                  __typename: 'JobRoleParticipant';
                  jobRoleParticipant: {
                    __typename?: 'JobRole';
                    contact?: {
                      __typename?: 'Contact';
                      id: string;
                      name?: string | null;
                      firstName?: string | null;
                      lastName?: string | null;
                    } | null;
                  };
                }
              | { __typename: 'OrganizationParticipant' }
              | { __typename: 'PhoneNumberParticipant' }
              | {
                  __typename: 'UserParticipant';
                  userParticipant: {
                    __typename?: 'User';
                    id: string;
                    firstName: string;
                    lastName: string;
                  };
                }
            >;
          }
        | { __typename: 'InteractionSession' }
        | { __typename: 'Issue'; id: string; createdAt: any; updatedAt: any }
        | {
            __typename: 'LogEntry';
            id: string;
            createdBy?: {
              __typename?: 'User';
              lastName: string;
              firstName: string;
            } | null;
          }
        | {
            __typename: 'Meeting';
            id: string;
            name?: string | null;
            attendedBy: Array<
              | { __typename: 'ContactParticipant' }
              | { __typename: 'EmailParticipant' }
              | { __typename: 'OrganizationParticipant' }
              | { __typename: 'UserParticipant' }
            >;
          }
        | {
            __typename: 'Note';
            id: string;
            createdBy?: {
              __typename?: 'User';
              firstName: string;
              lastName: string;
            } | null;
          }
        | { __typename: 'PageView'; id: string }
        | null;
    } | null;
  } | null;
};
