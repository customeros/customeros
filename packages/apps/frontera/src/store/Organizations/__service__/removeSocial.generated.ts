import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type RemoveSocialMutationVariables = Types.Exact<{
  socialId: Types.Scalars['ID']['input'];
}>;

export type RemoveSocialMutation = {
  __typename?: 'Mutation';
  social_Remove: { __typename?: 'Result'; result: boolean };
};
