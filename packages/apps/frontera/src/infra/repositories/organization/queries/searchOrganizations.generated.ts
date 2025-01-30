import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type SearchOrganizationsQueryVariables = Types.Exact<{
  limit?: Types.InputMaybe<Types.Scalars['Int']['input']>;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;

export type SearchOrganizationsQuery = {
  __typename?: 'Query';
  ui_organizations_search: {
    __typename?: 'OrganizationSearchResult';
    ids: Array<string>;
    totalElements: any;
    totalAvailable: any;
  };
};
