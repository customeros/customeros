import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateTagMutationVariables = Types.Exact<{
  input: Types.TagInput;
}>;

export type CreateTagMutation = {
  __typename?: 'Mutation';
  tag_Create: {
    __typename?: 'Tag';
    name: string;
    entityType: Types.EntityType;
    colorCode: string;
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
