import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type BuyDomainsMutationVariables = Types.Exact<{
  test: Types.Scalars['Boolean']['input'];
  paymentIntentId: Types.Scalars['String']['input'];
  domains: Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input'];
  username: Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input'];
  amount: Types.Scalars['Float']['input'];
  redirectWebsite?: Types.InputMaybe<Types.Scalars['String']['input']>;
}>;


export type BuyDomainsMutation = { __typename?: 'Mutation', mailstack_RegisterBuyDomainsWithMailboxes: { __typename?: 'Result', result: boolean } };
