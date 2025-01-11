import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetTagsByEntityTypeQueryVariables = Types.Exact<{
  entityType: Types.EntityType;
}>;

export type GetTagsByEntityTypeQuery = {
  __typename?: 'Query';
  tags_ByEntityType: Array<{
    __typename?: 'Tag';
    name: string;
    entityType: Types.EntityType;
    colorCode: string;
    metadata: { __typename?: 'Metadata'; id: string };
  }>;
};
