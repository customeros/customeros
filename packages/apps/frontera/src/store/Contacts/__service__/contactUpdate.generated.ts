import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateContactMutationVariables = Types.Exact<{
  input: Types.ContactUpdateInput;
}>;

export type UpdateContactMutation = {
  __typename?: 'Mutation';
  contact_Update: { __typename?: 'Contact'; id: string };
};
