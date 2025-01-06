import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type ValidateDomainsQueryVariables = Types.Exact<{
  domains: Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input'];
}>;


export type ValidateDomainsQuery = { __typename?: 'Query', mailstack_CheckUnavailableDomains: Array<string> };
