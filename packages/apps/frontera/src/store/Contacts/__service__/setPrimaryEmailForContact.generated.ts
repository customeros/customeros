import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type SetPrimaryEmailForContactMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  email: Types.Scalars['String']['input'];
}>;

export type SetPrimaryEmailForContactMutation = {
  __typename?: 'Mutation';
  emailSetPrimaryForContact: { __typename?: 'Result'; result: boolean };
};
