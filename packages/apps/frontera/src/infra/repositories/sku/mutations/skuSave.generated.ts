import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type SkuSaveMutationVariables = Types.Exact<{
  input: Types.SkuInput;
}>;

export type SkuSaveMutation = {
  __typename?: 'Mutation';
  sku_Save: {
    __typename?: 'Sku';
    id: string;
    name: string;
    price: number;
    type: Types.SkuType;
    archived: boolean;
  };
};
