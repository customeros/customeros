import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetPaymentIntentMutationVariables = Types.Exact<{
  domains: Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input'];
  usernames: Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input'];
  amount: Types.Scalars['Float']['input'];
}>;


export type GetPaymentIntentMutation = { __typename?: 'Mutation', mailstack_GetPaymentIntent: { __typename?: 'GetPaymentIntent', clientSecret: string } };
