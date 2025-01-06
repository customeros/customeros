import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateTagMutationVariables = Types.Exact<{
  input: Types.TagUpdateInput;
}>;


export type UpdateTagMutation = { __typename?: 'Mutation', tag_Update?: { __typename?: 'Tag', metadata: { __typename?: 'Metadata', id: string } } | null };
