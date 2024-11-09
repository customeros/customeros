import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateContactPhoneNumberMutationVariables = Types.Exact<{
  input: Types.PhoneNumberUpdateInput;
}>;

export type UpdateContactPhoneNumberMutation = {
  __typename?: 'Mutation';
  phoneNumber_Update: { __typename?: 'PhoneNumber'; id: string };
};
