import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetTagsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GetTagsQuery = {
  __typename?: 'Query';
  tags: Array<{
    __typename?: 'Tag';
    name: string;
    entityType: Types.EntityType;
    colorCode: string;
    metadata: { __typename?: 'Metadata'; id: string };
  }>;
};
