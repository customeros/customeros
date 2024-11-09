import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type DeleteContactMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
}>;

export type DeleteContactMutation = {
  __typename?: 'Mutation';
  contact_HardDelete: { __typename?: 'Result'; result: boolean };
};
