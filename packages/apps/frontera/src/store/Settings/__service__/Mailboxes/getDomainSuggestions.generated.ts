import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetDomainSuggestionsQueryVariables = Types.Exact<{
  domain: Types.Scalars['String']['input'];
}>;


export type GetDomainSuggestionsQuery = { __typename?: 'Query', mailstack_DomainPurchaseSuggestions: Array<string> };
