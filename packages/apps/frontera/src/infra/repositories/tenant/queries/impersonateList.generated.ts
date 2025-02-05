import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type TenantImpersonateListQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type TenantImpersonateListQuery = {
  __typename?: 'Query';
  tenant_impersonateList: Array<{
    __typename?: 'TenantImpersonateDetails';
    tenant: string;
    createdBy: string;
    personal: boolean;
  }>;
};
