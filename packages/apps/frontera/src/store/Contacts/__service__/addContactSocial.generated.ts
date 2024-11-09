import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type AddContactSocialMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  input: Types.SocialInput;
}>;

export type AddContactSocialMutation = {
  __typename?: 'Mutation';
  contact_AddSocial: { __typename?: 'Social'; id: string };
};
