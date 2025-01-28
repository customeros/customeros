import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type SkuArchiveMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type SkuArchiveMutation = {
  __typename?: 'Mutation';
  sku_Archive: { __typename?: 'Result'; result: boolean };
};
