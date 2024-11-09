import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ContactsQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
}>;

export type ContactsQuery = {
  __typename?: 'Query';
  contacts: {
    __typename?: 'ContactsPage';
    totalPages: number;
    totalElements: any;
    content: Array<{
      __typename?: 'Contact';
      id: string;
      name?: string | null;
      firstName?: string | null;
      lastName?: string | null;
      prefix?: string | null;
      description?: string | null;
      timezone?: string | null;
      updatedAt: any;
      createdAt: any;
      profilePhotoUrl?: string | null;
      metadata: { __typename?: 'Metadata'; id: string };
      flows: Array<{
        __typename?: 'Flow';
        metadata: { __typename?: 'Metadata'; id: string };
      }>;
      tags?: Array<{
        __typename?: 'Tag';
        id: string;
        name: string;
        source: Types.DataSource;
        updatedAt: any;
        createdAt: any;
        appSource: string;
        metadata: {
          __typename?: 'Metadata';
          id: string;
          source: Types.DataSource;
          sourceOfTruth: Types.DataSource;
          appSource: string;
          created: any;
          lastUpdated: any;
        };
      }> | null;
      organizations: {
        __typename?: 'OrganizationPage';
        totalElements: any;
        totalAvailable: any;
        content: Array<{
          __typename?: 'Organization';
          id: string;
          name: string;
          metadata: { __typename?: 'Metadata'; id: string };
        }>;
      };
      primaryEmail?: {
        __typename?: 'Email';
        id: string;
        primary: boolean;
        email?: string | null;
        emailValidationDetails: {
          __typename?: 'EmailValidationDetails';
          verified: boolean;
          verifyingCheckAll: boolean;
          isValidSyntax?: boolean | null;
          isRisky?: boolean | null;
          isFirewalled?: boolean | null;
          provider?: string | null;
          firewall?: string | null;
          isCatchAll?: boolean | null;
          canConnectSmtp?: boolean | null;
          deliverable?: Types.EmailDeliverable | null;
          isMailboxFull?: boolean | null;
          isRoleAccount?: boolean | null;
          isFreeAccount?: boolean | null;
          smtpSuccess?: boolean | null;
        };
      } | null;
      latestOrganizationWithJobRole?: {
        __typename?: 'OrganizationWithJobRole';
        organization: {
          __typename?: 'Organization';
          name: string;
          metadata: { __typename?: 'Metadata'; id: string };
        };
        jobRole: {
          __typename?: 'JobRole';
          id: string;
          primary: boolean;
          jobTitle?: string | null;
          description?: string | null;
          company?: string | null;
          startedAt?: any | null;
          endedAt?: any | null;
        };
      } | null;
      jobRoles: Array<{
        __typename?: 'JobRole';
        id: string;
        primary: boolean;
        jobTitle?: string | null;
        description?: string | null;
        company?: string | null;
        startedAt?: any | null;
        endedAt?: any | null;
      }>;
      locations: Array<{
        __typename?: 'Location';
        id: string;
        address?: string | null;
        locality?: string | null;
        postalCode?: string | null;
        country?: string | null;
        countryCodeA2?: string | null;
        countryCodeA3?: string | null;
        region?: string | null;
      }>;
      phoneNumbers: Array<{
        __typename?: 'PhoneNumber';
        id: string;
        e164?: string | null;
        rawPhoneNumber?: string | null;
        label?: Types.PhoneNumberLabel | null;
        primary: boolean;
      }>;
      emails: Array<{
        __typename?: 'Email';
        id: string;
        primary: boolean;
        email?: string | null;
        work?: boolean | null;
        emailValidationDetails: {
          __typename?: 'EmailValidationDetails';
          verified: boolean;
          verifyingCheckAll: boolean;
          isValidSyntax?: boolean | null;
          isRisky?: boolean | null;
          isFirewalled?: boolean | null;
          provider?: string | null;
          firewall?: string | null;
          isCatchAll?: boolean | null;
          canConnectSmtp?: boolean | null;
          deliverable?: Types.EmailDeliverable | null;
          isMailboxFull?: boolean | null;
          isRoleAccount?: boolean | null;
          isFreeAccount?: boolean | null;
          smtpSuccess?: boolean | null;
        };
      }>;
      socials: Array<{
        __typename?: 'Social';
        id: string;
        url: string;
        alias: string;
        followersCount: any;
      }>;
      enrichDetails: {
        __typename?: 'EnrichDetails';
        enrichedAt?: any | null;
        failedAt?: any | null;
        requestedAt?: any | null;
        emailEnrichedAt?: any | null;
        emailFound?: boolean | null;
        emailRequestedAt?: any | null;
      };
      connectedUsers: Array<{ __typename?: 'User'; id: string }>;
    }>;
  };
};
