import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateContactEmailMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  previousEmail: Types.Scalars['String']['input'];
  input: Types.EmailInput;
}>;

export type UpdateContactEmailMutation = {
  __typename?: 'Mutation';
  emailReplaceForContact: { __typename?: 'Email'; id: string };
};
