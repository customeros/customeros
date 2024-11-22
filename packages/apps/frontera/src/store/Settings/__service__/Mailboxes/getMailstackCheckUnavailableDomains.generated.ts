import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type CheckUnavailableDomainsQueryVariables = Types.Exact<{
  domains:
    | Array<Types.Scalars['String']['input']>
    | Types.Scalars['String']['input'];
}>;

export type CheckUnavailableDomainsQuery = {
  __typename?: 'Query';
  mailstack_CheckUnavailableDomains: Array<string>;
};
