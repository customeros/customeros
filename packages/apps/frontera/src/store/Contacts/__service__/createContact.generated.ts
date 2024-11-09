import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateContactMutationVariables = Types.Exact<{
  contactInput: Types.ContactInput;
}>;

export type CreateContactMutation = {
  __typename?: 'Mutation';
  contact_Create: string;
};
