import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type SearchContactsQueryVariables = Types.Exact<{
  limit?: Types.InputMaybe<Types.Scalars['Int']['input']>;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;


export type SearchContactsQuery = { __typename?: 'Query', ui_contacts_search: { __typename?: 'ContactSearchResult', ids: Array<string>, totalElements: any, totalAvailable: any } };
