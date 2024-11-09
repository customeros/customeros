import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type AddTagsToContactMutationVariables = Types.Exact<{
  input: Types.ContactTagInput;
}>;

export type AddTagsToContactMutation = {
  __typename?: 'Mutation';
  contact_AddTag: { __typename?: 'ActionResponse'; accepted: boolean };
};
