import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetContactsByIdsQueryVariables = Types.Exact<{
  ids?: Types.InputMaybe<
    Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input']
  >;
}>;

export type GetContactsByIdsQuery = {
  __typename?: 'Query';
  ui_contacts: Array<{
    __typename?: 'ContactUiDetails';
    id: string;
    createdAt: any;
    updatedAt: any;
    firstName: string;
    lastName: string;
    name: string;
    prefix: string;
    hide: boolean;
    description: string;
    timezone: string;
    profilePhotoUrl: string;
    enrichedAt?: any | null;
    enrichedFailedAt?: any | null;
    enrichedRequestedAt?: any | null;
    enrichedEmailRequestedAt?: any | null;
    enrichedEmailEnrichedAt?: any | null;
    enrichedEmailFound?: boolean | null;
    linkedInInternalId?: string | null;
    linkedInUrl?: string | null;
    linkedInAlias?: string | null;
    linkedInExternalId?: string | null;
    linkedInFollowerCount?: any | null;
    primaryOrganizationId?: string | null;
    primaryOrganizationName?: string | null;
    primaryOrganizationJobRoleId?: string | null;
    primaryOrganizationJobRoleTitle?: string | null;
    primaryOrganizationJobRoleDescription?: string | null;
    primaryOrganizationJobRoleStartDate?: any | null;
    primaryOrganizationJobRoleEndDate?: any | null;
    jobRoleIds: Array<string>;
    phones: Array<string>;
    connectedUsers: Array<string>;
    flows: Array<string>;
    emails: Array<{
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
    }>;
    tags: Array<{
      __typename?: 'Tag';
      name: string;
      colorCode: string;
      entityType: Types.EntityType;
      metadata: { __typename?: 'Metadata'; id: string };
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
  }>;
};
