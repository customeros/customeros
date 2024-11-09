import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type RemoveTagFromContactMutationVariables = Types.Exact<{
  input: Types.ContactTagInput;
}>;

export type RemoveTagFromContactMutation = {
  __typename?: 'Mutation';
  contact_RemoveTag: { __typename?: 'ActionResponse'; accepted: boolean };
};
