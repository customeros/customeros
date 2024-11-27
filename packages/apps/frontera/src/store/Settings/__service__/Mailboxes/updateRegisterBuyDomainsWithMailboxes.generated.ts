import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type RegisterBuyDomainsWithMailboxesMutationVariables = Types.Exact<{
  domains:
    | Array<Types.Scalars['String']['input']>
    | Types.Scalars['String']['input'];
  usernames:
    | Array<Types.Scalars['String']['input']>
    | Types.Scalars['String']['input'];
  amount: Types.Scalars['Float']['input'];
}>;

export type RegisterBuyDomainsWithMailboxesMutation = {
  __typename?: 'Mutation';
  mailstack_RegisterBuyDomainsWithMailboxes: {
    __typename?: 'RegisterBuyDomainWithMailboxes';
    id: string;
    clientSecret: string;
  };
};
