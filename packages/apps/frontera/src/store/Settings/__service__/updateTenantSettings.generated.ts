import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateTenantSettingsMutationVariables = Types.Exact<{
  input: Types.TenantSettingsInput;
}>;

export type UpdateTenantSettingsMutation = {
  __typename?: 'Mutation';
  tenant_UpdateSettings: {
    __typename?: 'TenantSettings';
    logoUrl: string;
    logoRepositoryFileId?: string | null;
    baseCurrency?: Types.Currency | null;
    billingEnabled: boolean;
    workspaceName?: string | null;
    workspaceLogo?: string | null;
  };
};
