import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type SkusQueryVariables = Types.Exact<{ [key: string]: never }>;

export type SkusQuery = {
  __typename?: 'Query';
  skus: Array<{
    __typename?: 'Sku';
    id: string;
    price: number;
    type: Types.SkuType;
    archived: boolean;
    name: string;
  }>;
};
