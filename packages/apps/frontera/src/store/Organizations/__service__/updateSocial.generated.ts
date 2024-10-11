import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateSocialMutationVariables = Types.Exact<{
  input: Types.SocialUpdateInput;
}>;

export type UpdateSocialMutation = {
  __typename?: 'Mutation';
  social_Update: { __typename?: 'Social'; id: string; url: string };
};
