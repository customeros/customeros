import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type TenantSettingsQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type TenantSettingsQuery = {
  __typename?: 'Query';
  tenantSettings: {
    __typename?: 'TenantSettings';
    logoUrl: string;
    logoRepositoryFileId?: string | null;
    baseCurrency?: Types.Currency | null;
    workspaceName?: string | null;
    workspaceLogo?: string | null;
    billingEnabled: boolean;
    opportunityStages: Array<{
      __typename?: 'TenantSettingsOpportunityStageConfiguration';
      id: string;
      value: string;
      order: number;
      label: string;
      visible: boolean;
      likelihoodRate: any;
    }>;
  };
};
