import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type AddSocialMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  input: Types.SocialInput;
}>;

export type AddSocialMutation = {
  __typename?: 'Mutation';
  organization_AddSocial: { __typename?: 'Social'; id: string; url: string };
};
