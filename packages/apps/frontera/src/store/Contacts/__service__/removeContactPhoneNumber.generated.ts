import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type RemoveContactPhoneNumberMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  id: Types.Scalars['ID']['input'];
}>;

export type RemoveContactPhoneNumberMutation = {
  __typename?: 'Mutation';
  phoneNumberRemoveFromContactById: { __typename?: 'Result'; result: boolean };
};
