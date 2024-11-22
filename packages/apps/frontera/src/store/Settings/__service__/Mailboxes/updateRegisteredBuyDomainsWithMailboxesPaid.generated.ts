import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type UpdateRegisteredBuyDomainsWithMailboxesPaidMutationVariables =
  Types.Exact<{
    id: Types.Scalars['String']['input'];
  }>;

export type UpdateRegisteredBuyDomainsWithMailboxesPaidMutation = {
  __typename?: 'Mutation';
  mailstack_RegisteredBuyDomainsWithMailboxesPaid: {
    __typename?: 'Result';
    result: boolean;
  };
};
