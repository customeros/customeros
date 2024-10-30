import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetArchivedOrganizationsAfterQueryVariables = Types.Exact<{
  date: Types.Scalars['Time']['input'];
}>;

export type GetArchivedOrganizationsAfterQuery = {
  __typename?: 'Query';
  organizations_HiddenAfter: Array<string>;
};
