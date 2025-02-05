import { Transport } from '@infra/transport';

import TenantSettingsDocument from './queries/tenantSettings.graphql';
import UpdateTenantSettingsDocument from './mutations/updateTenantSettings.graphql';
import {
  TenantSettingsQuery,
  TenantSettingsQueryVariables,
} from './queries/tenantSettings.generated.ts';
import {
  UpdateTenantSettingsMutation,
  UpdateTenantSettingsMutationVariables,
} from './mutations/updateTenantSettings.generated.ts';

export class SettingsRepository {
  private transport = Transport.getInstance();

  constructor() {}

  async getTenantSettings() {
    return this.transport.graphql.request<
      TenantSettingsQuery,
      TenantSettingsQueryVariables
    >(TenantSettingsDocument);
  }

  async updateTenantSettings(payload: UpdateTenantSettingsMutationVariables) {
    return this.transport.graphql.request<
      UpdateTenantSettingsMutation,
      UpdateTenantSettingsMutationVariables
    >(UpdateTenantSettingsDocument, payload);
  }
}
