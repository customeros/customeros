import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

import { TenantSettings, TenantSettingsInput } from '@graphql/types';

export class TenantStore {
  value: TenantSettings | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;

  constructor(public root: RootStore, public transportLayer: Transport) {
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    this.load();
  }

  async load() {
    try {
      this.isLoading = true;
      const repsonse =
        await this.transportLayer.graphql.request<TENANT_SETTINGS_QUERY_RESULT>(
          TENANT_SETTINGS_QUERY,
        );
      runInAction(() => {
        this.value = repsonse.tenantSettings;
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  update(updated: (value: TenantSettings) => TenantSettings) {
    this.value = updated(this.value as TenantSettings);
    this.save();
  }

  async save() {
    try {
      this.isLoading = true;
      await this.transportLayer.graphql.request<
        TENANT_SETTINGS_UPDATE_RESULT,
        { input: TenantSettingsInput }
      >(TENANT_SETTINGS_UPDATE_MUTATION, {
        input: {
          ...(this.value as TenantSettingsInput),
          patch: true,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}

type TENANT_SETTINGS_QUERY_RESULT = {
  tenantSettings: TenantSettings;
};
const TENANT_SETTINGS_QUERY = gql`
  query TenantSettings {
    tenantSettings {
      logoUrl
      logoRepositoryFileId
      baseCurrency
      billingEnabled
      opportunityStages
    }
  }
`;

type TENANT_SETTINGS_UPDATE_RESULT = {
  tenant_UpdateSettings: TenantSettings;
};
const TENANT_SETTINGS_UPDATE_MUTATION = gql`
  mutation UpdateTenantSettings($input: TenantSettingsInput!) {
    tenant_UpdateSettings(input: $input) {
      logoUrl
      logoRepositoryFileId
      baseCurrency
      billingEnabled
    }
  }
`;
