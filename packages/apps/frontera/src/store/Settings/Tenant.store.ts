import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';

import { TenantSettings } from '@graphql/types';

export class TenantStore {
  value: TenantSettings | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;

  constructor(private root: RootStore, private transportLayer: Transport) {
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
      this.value = repsonse.tenantSettings;
      this.isBootstrapped = true;
    } catch (err) {
      this.error = (err as Error).message;
    } finally {
      this.isLoading = false;
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
    }
  }
`;
