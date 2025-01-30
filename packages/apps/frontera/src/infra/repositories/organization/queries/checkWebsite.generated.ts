import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type CheckWebsiteQueryVariables = Types.Exact<{
  website: Types.Scalars['String']['input'];
}>;

export type CheckWebsiteQuery = {
  __typename?: 'Query';
  organization_CheckWebsite: {
    __typename?: 'WebsiteCheckDetails';
    accepted: boolean;
    primary: boolean;
    domain: string;
    primaryDomain: string;
    globalOrganization?: {
      __typename?: 'GlobalOrganization';
      id: any;
      name: string;
      logoUrl: string;
      iconUrl: string;
      domains: Array<string>;
      website: string;
      primaryDomain: string;
      organizationId?: string | null;
    } | null;
  };
};
