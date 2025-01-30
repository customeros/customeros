import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type SearchGlobalOrganizationsQueryVariables = Types.Exact<{
  searchTerm: Types.Scalars['String']['input'];
  limit: Types.Scalars['Int']['input'];
}>;

export type SearchGlobalOrganizationsQuery = {
  __typename?: 'Query';
  globalOrganizations_Search: Array<{
    __typename?: 'GlobalOrganization';
    id: any;
    name: string;
    primaryDomain: string;
    website: string;
    logoUrl: string;
    iconUrl: string;
    organizationId?: string | null;
  }>;
};
