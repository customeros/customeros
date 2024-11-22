import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type MailstackDomainPurchaseSuggestionsQueryVariables = Types.Exact<{
  domain: Types.Scalars['String']['input'];
}>;

export type MailstackDomainPurchaseSuggestionsQuery = {
  __typename?: 'Query';
  mailstack_DomainPurchaseSuggestions: Array<string>;
};
