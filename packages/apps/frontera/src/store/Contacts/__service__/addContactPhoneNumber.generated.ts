import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type AddContactPhoneNumberMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  input: Types.PhoneNumberInput;
}>;

export type AddContactPhoneNumberMutation = {
  __typename?: 'Mutation';
  phoneNumberMergeToContact: {
    __typename?: 'PhoneNumber';
    id: string;
    rawPhoneNumber?: string | null;
  };
};
