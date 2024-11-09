import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateContactSocialMutationVariables = Types.Exact<{
  input: Types.SocialUpdateInput;
}>;

export type UpdateContactSocialMutation = {
  __typename?: 'Mutation';
  social_Update: { __typename?: 'Social'; id: string };
};
