import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetTenantIndustriesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetTenantIndustriesQuery = {
  __typename?: 'Query';
  industries_InUse: Array<{
    __typename?: 'Industry';
    code: string;
    name: string;
  }>;
};
