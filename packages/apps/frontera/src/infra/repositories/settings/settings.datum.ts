import { TenantSettingsQuery } from './queries/tenantSettings.generated';
//
export type SettingsDatum = NonNullable<TenantSettingsQuery['tenantSettings']>;
